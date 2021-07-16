package util

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func WaitForCtrlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func Debug(data interface{}) {
	fmt.Printf("%+v\n", data)
	// spew.Dump(data)
}

func HexString(payload []byte) string {
	return strings.ToUpper(hex.EncodeToString(payload))
}
