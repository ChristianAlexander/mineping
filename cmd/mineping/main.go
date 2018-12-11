package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/christianalexander/mineping"
)

var config struct {
	serverName string
	serverPort uint
}

func init() {
	flag.StringVar(&config.serverName, "n", "beardcraft.us", "The server to connect to, defaults to beardcraft.us")
	flag.UintVar(&config.serverPort, "p", 25565, "The port to connect to, defaults to 25565")
	flag.Parse()
}

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.serverName, config.serverPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to '%s:%d': %v", config.serverName, config.serverPort, err)
		os.Exit(1)
	}
	defer conn.Close()

	err = mineping.WriteServerListPingRequest(conn, config.serverName, uint16(config.serverPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write handshake packet: %v", err)
		os.Exit(1)
	}

	response, err := mineping.ReadServerListPing(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read handshake response: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%d/%d players online\n\n", response.Players.Online, response.Players.Max)

	for _, p := range response.Players.Sample {
		fmt.Println(p.Name)
	}
}
