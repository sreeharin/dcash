package main

import (
	"bytes"
	"errors"
	"net"
	"log"
)

// Start is the entrypoint of node module
// It runs the server 
func Start(c *Cache) {
	l, err := net.Listen("tcp", "localhost:6368")	
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO: fix logging
			log.Print(err)
			continue
		}

		go handleConn(c, conn)
	}
}

// handleConn is used to handle the commands send by the client
func handleConn(c *Cache, conn net.Conn) {
	defer conn.Close()	
	
	MAX_SIZE := 1024
	buffer := make([]byte, MAX_SIZE)

	_, err := conn.Read(buffer)
	if err != nil {
		// TODO: fix logging
		log.Print(err)
		return
	}

	command, input, found := bytes.Cut(buffer, []byte(" "))
	if !found {
		return
	}

	switch string(command) {
	case "SET":
		if err := handleSET(c, input); err != nil {
			conn.Write([]byte("Unexpected error while storing. Error: "+err.Error()+"\n"))
			return
		}
		conn.Write([]byte("OK\n"))
	default:
		// TODO: log
		conn.Write([]byte("Error: Unsupported Command."))
	}

}

// handleSET handles the 'SET' command
// it adds an entry to the cache
func handleSET(c *Cache, input []byte) error {
	key, value, exists := bytes.Cut(input, []byte(" "))

	if !exists {
		return errors.New("Corrupted data sent by client.")
	}

	c.Set(string(key), value)

	return nil
}
