package command

func contains(b []byte, value ...uint8) bool {
	for _, v := range value {
		if toUint8(b) == v {
			return true
		}
	}
	return false
}

func max(b []byte, max uint8) bool {
	return toUint8(b) <= max
}

func between(b []byte, min, max uint8) bool {
	return toUint8(b) >= min && toUint8(b) <= max
}
