package interfaces

type IClientSocket interface {
	Connect() error
	Send(data []interface{}) error
	Await() ([]interface{}, error)
	JoinGroup(groupID string) error
	LeaveGroup(groupID string) error
	BroadcastToGroup(groupID string, data interface{}) error
	GetGroups() []string
	Close() error
}
