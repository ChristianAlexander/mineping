package encoding

import "io"

func WriteUint8(writer io.Writer, val uint8) error {
	_, err := writer.Write([]byte{val})
	return err
}

func ReadUint8(reader io.Reader) (uint8, error) {
	var dat [1]byte
	_, err := reader.Read(dat[:1])
	return dat[0], err
}
