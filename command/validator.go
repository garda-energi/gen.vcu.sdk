package command

func max(b []byte, max uint8) bool {
	return toUint8(b) <= max
}

func between(b []byte, min, max uint8) bool {
	return toUint8(b) >= min && toUint8(b) <= max
}
