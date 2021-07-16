package command

import (
	"errors"
	"reflect"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type Command struct {
	client *mqtt.Client
}

func New(client *mqtt.Client) *Command {
	return &Command{
		client: client,
	}
}

func (c *Command) GenInfo() (string, error) {
	cmd, err := getCmdPacket(CMDC_GEN, CMD_SUBCODE(CMD_GEN_INFO))
	if err != nil {
		return "", err
	}

	util.Debug(cmd)
	return "", nil
}

func (c *Command) ExecuteEmpty(code CMD_CODE, subCode CMD_SUBCODE) error {

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

func getCmdPacket(code CMD_CODE, subCode CMD_SUBCODE) (CmdPacket, error) {
	// for _, cmd := range CMD_LIST {
	// 	if cmd.Code == code && cmd.SubCode == subCode {
	// 		return cmd, nil
	// 	}
	// }

	cmd, ok := CmdList[code][subCode]
	if ok {
		return cmd, nil
	}
	return CmdPacket{}, errors.New("command code not found")
}
