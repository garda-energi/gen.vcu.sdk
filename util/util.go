package util

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForCtrlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func Debug(data interface{}) {
	fmt.Printf("%+v\n", data)
}
