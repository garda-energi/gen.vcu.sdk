package command

import (
	"errors"
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type Cmd struct{}

func (c *Cmd) ExecuteEmpty(code CMD_CODE, subCode CMD_SUBCODE) error {
	// cmd, err := getCmdPacket(code, subCode)
	// if err != nil {
	// 	return err
	// }

	cmd, ok := CmdList[code][subCode]
	if !ok {
		return errors.New("command not found")
	}

	if cmd.Tipe != reflect.Invalid {
		return errors.New("command need payload")
	}

	util.Debug(cmd)

	return nil
}

// func (c *Command) GenInfo() (string, error) {

// }

func getCmdPacket(code CMD_CODE, subCode CMD_SUBCODE) (CmdPacket, error) {
	for _, cmd := range CMD_LIST {
		if cmd.Code == code && cmd.SubCode == subCode {
			return cmd, nil
		}
	}
	return CmdPacket{}, errors.New("command code not found")
}
