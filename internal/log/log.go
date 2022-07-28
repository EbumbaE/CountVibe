package log

type Logger interface{
	Process(...any)
	Error(...any)
}