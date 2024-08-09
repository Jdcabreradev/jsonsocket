package interfaces

type IServerSocket interface {
	Bind() string
	Close() error
}
