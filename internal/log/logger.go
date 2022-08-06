package log

import(
	"os"
	"log"
)

type FileLogger struct{
	file *os.File
	processLogger *log.Logger
	errorLogger *log.Logger
}

func NewLogger(path string)(*FileLogger, error){

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	lg := &FileLogger{}

	lg.file = file
	lg.processLogger = log.New(file, "Process: ", log.Ldate|log.Ltime|log.Lshortfile)
	lg.errorLogger = log.New(file, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

	return lg, nil
}

func (lg *FileLogger) Error(a ...any) {
	lg.errorLogger.Println(a...)
}

func (lg *FileLogger) Process(a ...any) {
	lg.processLogger.Println(a...)
}

func (lg *FileLogger) CloseFile() error {
	return lg.file.Close()
}