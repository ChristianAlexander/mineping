package encoding

import "io"

func WriteString(writer io.Writer, val string) error {
	valBytes := []byte(val)
	err := WriteVarInt(writer, len(valBytes))
	if err != nil {
		return err
	}

	_, err = writer.Write(valBytes)
	return err
}

func ReadString(reader io.Reader) (string, error) {
	length, err := ReadVarInt(reader)
	if err != nil {
		return "", err
	}

	bytes := make([]byte, length)
	_, err = reader.Read(bytes)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
