package jsonsocket

import (
	"sync"

	"github.com/Jdcabreradev/logify/v3"
	logify_enums "github.com/Jdcabreradev/logify/v3/enums"
)

// socketBroker manages communication between clients and groups, handling subscriptions,
// broadcasting messages, and managing client connections and disconnections.
type socketBroker struct {
	clients    map[string]*socketProcess // Map of connected clients by their unique ID.
	groups     map[string]*SocketGroup   // Map of groups with their respective clients.
	clientSync sync.RWMutex              // Mutex to synchronize access to the clients map.
	groupSync  sync.RWMutex              // Mutex to synchronize access to the groups map.
	Logger     *logify.Logger            // Logger for logging events and errors.
}

// newBroker creates and returns a new socketBroker with initialized maps for clients and groups.
func newBroker(logger *logify.Logger) *socketBroker {
	return &socketBroker{
		clients: make(map[string]*socketProcess),
		groups:  make(map[string]*SocketGroup),
		Logger:  logger,
	}
}

// AddClient adds a new client to the broker, registering it with its unique ID.
func (b *socketBroker) AddClient(clientId string, client *socketProcess) {
	b.clientSync.Lock()          // Lock the clients map for writing.
	defer b.clientSync.Unlock()  // Ensure unlocking after operation.
	b.clients[clientId] = client // Add the client to the map.
	b.Logger.Log(logify_enums.INFO, "New client added with process id: "+clientId)
}

// RemoveClient removes a client from the broker and unsubscribes it from all groups.
func (b *socketBroker) RemoveClient(clientId string) {
	b.clientSync.Lock() // Lock the clients map for writing.
	client, exists := b.clients[clientId]
	if !exists {
		b.clientSync.Unlock() // Unlock and return if the client doesn't exist.
		return
	}
	delete(b.clients, clientId) // Remove the client from the map.
	b.clientSync.Unlock()       // Unlock after modification.

	// Unsubscribe client from all groups it was subscribed to.
	b.groupSync.Lock() // Lock the groups map for writing.
	for groupId := range client.SubscribedGroups {
		group, exists := b.groups[groupId]
		if exists {
			delete(group.Clients, clientId) // Remove the client from the group.
			b.Logger.Log(logify_enums.INFO, "Client "+clientId+" unsubscribed from group "+groupId)
			if len(group.Clients) == 0 {
				delete(b.groups, groupId) // Delete the group if it's empty.
				b.Logger.Log(logify_enums.INFO, "Group "+groupId+" was deleted due to zero clients listening")
			}
		}
	}
	b.groupSync.Unlock() // Unlock after modification.

	b.Logger.Log(logify_enums.INFO, "Client "+clientId+" was disconnected and unsubscribed from all groups")
}

// SubscribeToGroup subscribes a client to a specific group, creating the group if it doesn't exist.
func (b *socketBroker) SubscribeToGroup(groupId, clientId string) {
	b.clientSync.RLock() // Lock the clients map for reading.
	client, exists := b.clients[clientId]
	b.clientSync.RUnlock() // Unlock after reading.

	if !exists {
		return
	}

	b.groupSync.Lock() // Lock the groups map for writing.
	group, exists := b.groups[groupId]
	if !exists { // Create the group if it doesn't exist.
		group = &SocketGroup{
			Id:      groupId,
			Clients: make(map[string]*socketProcess),
		}
		b.groups[groupId] = group
		b.Logger.Log(logify_enums.INFO, "New group created: "+groupId)
	}
	group.Clients[clientId] = client         // Add the client to the group.
	client.SubscribedGroups[groupId] = group // Update the client's subscribed groups.
	b.groupSync.Unlock()                     // Unlock after modification.
	b.Logger.Log(logify_enums.INFO, "Client "+clientId+" subscribed to group "+groupId)
}

// UnsubscribeFromGroup unsubscribes a client from a specific group.
func (b *socketBroker) UnsubscribeFromGroup(groupId, clientId string) {
	b.clientSync.RLock() // Lock the clients map for reading.
	client, exists := b.clients[clientId]
	b.clientSync.RUnlock() // Unlock after reading.

	if !exists {
		return
	}

	b.groupSync.Lock() // Lock the groups map for writing.
	group, exists := b.groups[groupId]
	if exists {
		delete(group.Clients, clientId) // Remove the client from the group.
		b.Logger.Log(logify_enums.INFO, "Client "+clientId+" unsubscribed from group "+groupId)
		if len(group.Clients) == 0 {
			delete(b.groups, groupId) // Delete the group if it's empty.
			b.Logger.Log(logify_enums.INFO, "Group "+groupId+" was deleted due to zero clients listening")
		}
	}
	delete(client.SubscribedGroups, groupId) // Remove the group from the client's subscribed groups.
	b.groupSync.Unlock()                     // Unlock after modification.
}

// BroadcastToGroup sends a message to all clients subscribed to a specific group.
func (b *socketBroker) BroadcastToGroup(groupId string, data []interface{}) {
	b.groupSync.RLock() // Lock the groups map for reading.
	group, exists := b.groups[groupId]
	if exists {
		b.Logger.Log(logify_enums.DEBUG, "Group "+groupId+" is receiving a broadcast message")
		for _, client := range group.Clients { // Send the data to each client in the group.
			client.Response(data)
		}
	}
	b.groupSync.RUnlock() // Unlock after reading.
}
