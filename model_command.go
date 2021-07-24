package sdk

type HeaderCommand struct {
	Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type CommandPacket struct {
	Header  *HeaderCommand
	Payload []byte
}

type HeaderResponse struct {
	HeaderCommand
	ResCode resCode `type:"uint8"`
}

type ResponsePacket struct {
	Header  *HeaderResponse
	Message []byte `type:"char"`
}
