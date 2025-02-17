package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"sphinx/config"
)

var logger zerolog.Logger

func InitLog() error {
	//console := zerolog.ConsoleWriter{Out: os.Stderr}
	//file := lumberjack.Logger{}
	//ouput := zerolog.MultiLevelWriter(console)
	return nil
}

func main() {
	cfg := config.NewConfig(nil)

	value := config.GetConfig[string](cfg, "log_level")

	fmt.Println(value)
}
