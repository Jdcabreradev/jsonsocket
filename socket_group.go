package jsonsocket

import (
	"sync"
)

// SocketGroup represents a communication channel with its clients and synchronization mechanism.
type SocketGroup struct {
	Id      string                    // Id to map specific group.
	Clients map[string]*socketProcess // Stores all the clients listening to the group.
	Sync    sync.RWMutex              // Sync read and write to support multithreading.
}
