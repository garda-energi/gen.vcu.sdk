package command_test

import (
	"testing"

	"github.com/pudjamansyurin/gen_vcu_sdk/command"
)

func TestCommand(t *testing.T) {
	cmd := command.New(nil)

	cmd.GenInfo()
}
