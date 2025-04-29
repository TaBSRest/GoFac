package interfaces

type TaBSError interface {
	error
	GetMessage() string
}
