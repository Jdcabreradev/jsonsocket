package models

import "sync"

// SocketChannel represents a communication channel with its clients and synchronization mechanism.
type SocketChannel struct {
	Id      string                    // Id to map specific channel.
	Clients map[string]*SocketMessage // Stores all the clients listening to the channel.
	Sync    sync.RWMutex              // Sync read and write to support multithreading.
}
