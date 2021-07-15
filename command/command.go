package command

import (
	"errors"
)

// type Command struct{}

// func (c *Command) ExecuteEmpty(code CMD_CODE, subCode CMD_SUBCODE) error {
// 	cmd, err := getCmd(code, subCode)
// 	if err != nil {
// 		return err
// 	}

// 	if cmd.Tipe != reflect.Invalid {
// 		return errors.New("command has payload")
// 	}

// 	return nil
// }

// func (c *Command) GenInfo() (string, error) {

// }

func getCmd(code CMD_CODE, subCode CMD_SUBCODE) (Cmd, error) {
	for _, cmd := range CMD_LIST {
		if cmd.Code == code && cmd.SubCode == subCode {
			return cmd, nil
		}
	}
	return Cmd{}, errors.New("command code not found")
}
