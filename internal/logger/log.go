package logger

type Logger interface {
	Process(...any)
	Error(...any)
}
