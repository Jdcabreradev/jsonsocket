package interfaces

type IServerSocket interface {
	Bind() string
	Listen(clientId string) ([]interface{}, error)
	Response(ProcessID string, data []interface{}) error
	Close(ProcessID string)
}
