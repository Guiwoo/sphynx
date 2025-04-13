package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/sphynx/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ContextHook struct{}

func (c *ContextHook) Run(e *zerolog.Event, lv zerolog.Level, msg string) {
	ctx := e.GetCtx()
	value, _ := ctx.Value("trid").(string)
	e.Str("trid", value)
}

func New() *zerolog.Logger {

	fileWriter := &lumberjack.Logger{
		Filename:   config.Get[string]("log_folder"),
		MaxSize:    config.Get[int]("log_max_size"),
		MaxBackups: config.Get[int]("log_max_backup"),
		MaxAge:     config.Get[int]("log_max_age"),
		Compress:   true,
	}

	multiLevelWriter := zerolog.MultiLevelWriter(
		consloeWriter(),
		fileWriter,
	)

	log := zerolog.
		New(multiLevelWriter).
		Hook(&ContextHook{}).
		With().
		Caller().
		Timestamp().
		Logger()

	return &log
}

func consloeWriter() zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatTimestamp = func(i interface{}) string {
		if t, ok := i.(time.Time); ok {
			return fmt.Sprintf("[ %s ]", t.Format("2006-01-02 15:04:05"))
		}
		return fmt.Sprintf("[ %s ]", i)
	}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%-6s]", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("[ msg : %-6s ]", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("[ %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s ]", i)
	}
	output.FormatCaller = func(i interface{}) string {
		var caller string
		if s, ok := i.(string); ok {
			caller = s
		}

		dirPath, _ := os.Getwd()
		relPath, _ := filepath.Rel(dirPath, caller)

		parts := strings.Split(relPath, ":")
		file := parts[0]
		line := parts[1]

		return fmt.Sprintf("[ %-s:%-s ]", file, line)
	}
	output.FormatPartValueByName = func(i interface{}, name string) string {
		if name == "trid" {
			return fmt.Sprintf("[ %s ]", i)
		}
		return i.(string)
	}

	output.PartsOrder = []string{"time", "level", "trid", "message", "caller"}
	output.FieldsExclude = []string{"trid"}

	return output
}
