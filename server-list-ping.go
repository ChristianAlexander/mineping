package mineping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/christianalexander/mineping/encoding"
)

type ServerListPingResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon string `json:"favicon"`
}

// WriteServerListPingRequest writes a the request for a SLP.
func WriteServerListPingRequest(client net.Conn, serverName string, serverPort uint16) error {
	// -- Create packet --
	buf := new(bytes.Buffer)

	// Handshake packet ID
	encoding.WriteVarInt(buf, 0x00)

	// Version 404 (1.13.2)
	encoding.WriteVarInt(buf, 404)

	encoding.WriteString(buf, serverName)

	encoding.WriteUint16(buf, serverPort)

	// Status request
	encoding.WriteVarInt(buf, 1)

	// -- Write packet --
	err := encoding.WriteVarInt(client, buf.Len())
	if err != nil {
		return fmt.Errorf("failed to write packet length: %v", err)
	}

	_, err = buf.WriteTo(client)
	if err != nil {
		return fmt.Errorf("failed to write packet: %v", err)
	}

	// -- Create request packet --
	buf = new(bytes.Buffer)

	encoding.WriteVarInt(buf, 0x00)

	// -- Write packet --
	err = encoding.WriteVarInt(client, buf.Len())
	if err != nil {
		return fmt.Errorf("failed to write packet length: %v", err)
	}

	_, err = buf.WriteTo(client)
	if err != nil {
		return fmt.Errorf("failed to write packet: %v", err)
	}

	return nil
}

// ReadServerListPing reads a the response to a SLP request.
func ReadServerListPing(client net.Conn) (*ServerListPingResponse, error) {
	length, err := encoding.ReadVarInt(client)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %v", err)
	}

	packet := make([]byte, length)

	_, err = io.ReadFull(client, packet)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet: %v", err)
	}

	buf := bytes.NewBuffer(packet)
	id, err := encoding.ReadVarInt(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %v", err)
	}

	if id != 0x00 {
		return nil, fmt.Errorf("expected packet with ID 0x00, but received '%#x' instead", id)
	}

	responseStr, err := encoding.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response ServerListPingResponse
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &response, nil
}
