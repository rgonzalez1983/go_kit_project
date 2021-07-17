package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go_kit_project/internal/config"
	"go_kit_project/internal/handler"
	"go_kit_project/internal/middleware"
	"go_kit_project/internal/static"
	"os"
	"os/signal"
	"syscall"
)

var httpHandler = handler.App{}

func main() {
	if err, isConfigurable := config.ConfigEnv(); !isConfigurable {
		fmt.Printf(""+static.MsgResponseStartError+", %s", err)
		values := []interface{}{static.KeyType, static.ERROR, static.KeyMessage, static.MsgResponseStartError + ", " + err.Error()}
		middleware.LoggingOperation(httpHandler.Logg, values...)
	} else {
		fmt.Println(static.MsgResponseStartProcess)
		addr := ":" + viper.GetString(static.APP_PORT)
		_ = httpHandler.Initialize(static.ValueEmpty, static.ValueEmpty)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		_ = httpHandler.Run(addr)
	}
}
