package jsonsocket

type IClientSocket interface {
	Connect() error
	Close() error
}
