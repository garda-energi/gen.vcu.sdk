package command

import "testing"

func TestCommand(t *testing.T) {
	cmd := Cmd{}

	cmd.ExecuteEmpty(CMDC_GEN, CMD_SUBCODE(CMD_GEN_INFO))

}
