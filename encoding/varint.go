package encoding

import "io"

func WriteVarInt(writer io.Writer, val int) error {
	for val >= 0x80 {
		err := WriteUint8(writer, byte(val)|0x80)
		if err != nil {
			return nil
		}
		val >>= 7
	}
	err := WriteUint8(writer, byte(val))
	return err
}

func ReadVarInt(reader io.Reader) (int, error) {
	var result int
	var bytes byte

	for {
		b, err := ReadUint8(reader)
		if err != nil {
			return 0, nil
		}

		result |= int(uint(b&0x7F) << uint(bytes*7))
		bytes++
		if (b & 0x80) == 0x80 {
			continue
		}
		break
	}

	return result, nil
}
