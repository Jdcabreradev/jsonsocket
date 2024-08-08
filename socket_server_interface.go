package jsonsocket

type IServerSocket interface {
	Bind() string
	Close() error
}
