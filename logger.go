package sl

import (
	"fmt"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"os"
	"runtime"
	"time"
)

type LogLevel int

const (
	Trace   LogLevel = 0
	Debug   LogLevel = 1
	Info    LogLevel = 2
	Warning LogLevel = 3
	Error   LogLevel = 4
	Fatal   LogLevel = 5
)

func getLevelName(level LogLevel) string {
	switch level {
	case Debug:
		return "Debug"
	case Trace:
		return "Trace"
	case Info:
		return "Info"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Fatal:
		return "Fatal"
	}
	return ""
}

type Logger interface {
	Log(level LogLevel, message string, params ...any)
	LogTrace(message string, a ...any)
	LogDebug(message string, a ...any)
	LogInformation(message string, a ...any)
	LogWarning(message string, a ...any)
	LogError(err error, message string)
	LogFatal(err error, message string)
}
type LoggerConfig struct {
	ElasticConfig     *ElasticConfig
	ServiceName       string
	PrintServiceName  bool
	PrintErrorLogLine bool
}

type LogMessage struct {
	TimeStamp time.Time `json:"@timeStamp"`
	Level     string    `json:"level"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
}

type ElasticConfig struct {
	Host          string `json:"host"`
	User          string `json:"user"`
	Password      string `json:"password"`
	IndexTemplate string `json:"indexTemplate"`
}

type logger struct {
	printErrorLogLine bool
	//elasticConfig     *ElasticConfig
	elasticClient *elastic.Client
}

func (l logger) Log(level LogLevel, message string, params ...any) {
	_message := fmt.Sprintf(message, params...)
	if l.printErrorLogLine {
		_, _, name := getCallerName()
		fmt.Printf("%s [%s] [%s] %s\n", time.Now().Format("2006-01-02T15:04:05"), getLevelName(level), name, _message)
	} else {
		fmt.Printf("%s [%s] %s\n", time.Now().Format("2006-01-02T15:04:05"), getLevelName(level), _message)
	}
}

func (l logger) LogTrace(message string, a ...any) {
	l.Log(Trace, message, a...)
}

func (l logger) LogDebug(message string, a ...any) {
	l.Log(Debug, message, a...)
}

func (l logger) LogInformation(message string, a ...any) {
	l.Log(Debug, message, a...)
}

func (l logger) LogWarning(message string, a ...any) {
	l.Log(Warning, message, a...)
}

func (l logger) LogError(err error, message string) {
	file, line, _ := getErrorCallerName()
	errorMessage := fmt.Sprintf("%s at %s:%d\n%+v", message, file, line, err)
	l.Log(Error, errorMessage)
}

func (l logger) LogFatal(err error, message string) {
	file, line, _ := getErrorCallerName()
	errorMessage := fmt.Sprintf("%s at %s:%d\n%+v", message, file, line, err)
	l.Log(Fatal, errorMessage)
}

func NewLogger(config LoggerConfig) Logger {
	if config.ElasticConfig != nil {
		elcli, _ := elastic.NewClient(elastic.Config{Logger: &elastictransport.JSONLogger{Output: os.Stdout}})
		return &logger{
			printErrorLogLine: config.PrintErrorLogLine,
			//elasticConfig:     config.ElasticConfig,
			elasticClient: elcli,
		}
	} else {
		return &logger{
			printErrorLogLine: config.PrintErrorLogLine,
			//elasticConfig:     config.ElasticConfig,
		}
	}

}

func getCallerName() (string, int, string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(4, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return file, line, f.Name()
}

func getErrorCallerName() (string, int, string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return file, line, f.Name()
}
