package command

type RSP_CODE uint8

const (
	RSP_CODE_ERROR RSP_CODE = iota
	RSP_CODE_OK
	RSP_CODE_INVALID
)
