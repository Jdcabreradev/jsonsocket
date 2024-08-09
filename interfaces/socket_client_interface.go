package interfaces

type IClientSocket interface {
	Connect() error
	Close() error
}
