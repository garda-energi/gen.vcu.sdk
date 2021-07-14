package util

import (
	"encoding/hex"
	"os"
	"os/signal"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func WaitForCtrlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func Debug(data interface{}) {
	// fmt.Printf("%+v\n", data)
	spew.Dump(data)
}

func HexString(payload []byte) string {
	return strings.ToUpper(hex.EncodeToString(payload))
}
