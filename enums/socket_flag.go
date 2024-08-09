package enums

// SocketFlag represents different types of socket flags.
type SocketFlag uint8

const (
	START_TX          SocketFlag = iota + 1 // Start Data transmission flag.
	TX                                      // Continuous data transmission flag.
	END_TX                                  // End data transmission.
	CLOSE_CONN                              // Close connection flag.
	SUB_CHANNEL                             // Subscribe to channel flag.
	BROADCAST_CHANNEL                       // Broadcast to channel flag.
	UNSUB_CHANNEL                           // Unsubscibe to channel flag.
)
