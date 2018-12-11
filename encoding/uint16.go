package encoding

import (
	"encoding/binary"
	"io"
)

func WriteUint16(writer io.Writer, val uint16) error {
	var valArr [2]byte
	binary.BigEndian.PutUint16(valArr[:2], val)
	_, err := writer.Write(valArr[:2])
	return err
}

func ReadUint16(reader io.Reader) (uint16, error) {
	var dat [2]byte
	_, err := reader.Read(dat[:2])
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(dat[:2]), nil
}
