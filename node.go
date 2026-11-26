package main

import (
	"bytes"
	"errors"
	"log"
	"net"
)

// Start is the entrypoint of node module
// It runs the server
func Start(c *Cache) net.Listener {
	l, err := net.Listen("tcp", "localhost:6368")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				// TODO: fix logging
				log.Print(err)
				continue
			}

			go handleConn(c, conn)
		}
	}()

	return l
}

// handleConn is used to handle the commands send by the client
// TODO: Get rid of '\n'
func handleConn(c *Cache, conn net.Conn) {
	defer conn.Close()

	MAX_SIZE := 1024
	buffer := make([]byte, MAX_SIZE)
	n, err := conn.Read(buffer)

	if err != nil {
		// TODO: fix logging
		log.Print(err)
		return
	}

	command, input, found := bytes.Cut(buffer[:n], []byte(" "))
	if !found {
		return
	}

	switch string(command) {
	case "SET":
		if err := handleSET(c, input); err != nil {
			conn.Write([]byte("Error while storing data. " + err.Error() + "\n"))
			return
		}
		conn.Write([]byte("OK\n"))
	case "GET":
		//TODO: handle input sanitization
		key := string(input)
		val, err := c.Get(key)
		if err != nil {
			conn.Write([]byte("Error while retrieving data. " + err.Error() + "\n"))
			return
		}
		conn.Write(val)
	case "DEL":
		//TODO: handle input sanitization
		if err := handleDELETE(c, input); err != nil {
			conn.Write([]byte("Error while deleting data. " + err.Error() + ""))
			return
		}
		conn.Write([]byte("OK"))
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

func handleDELETE(c *Cache, input []byte) error {
	// TODO: sanitize key
	key := string(input)

	return c.Delete(key)
}
