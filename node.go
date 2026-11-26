package main

import (
	"bytes"
	"errors"
	"io"
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

	maxSize := 1024
	buffer := make([]byte, maxSize)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			// close the connection
			if errors.Is(err, io.EOF) {
				return
			}

			// TODO: fix logging
			log.Print(err)
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
			conn.Write([]byte("OK"))
		case "GET":
			val, err := handleGET(c, input)
			if err != nil {
				conn.Write([]byte("Error while retrieving data. " + err.Error() + "\n"))
				return
			}
			conn.Write(val)
		case "DEL":
			if err := handleDELETE(c, input); err != nil {
				conn.Write([]byte("Error while deleting data. " + err.Error() + ""))
				return
			}
			conn.Write([]byte("OK"))
		case "EXISTS":
			val, err := handleEXISTS(c, input)
			if err != nil {
				conn.Write([]byte("Error while checking key exists. " + err.Error() + ""))
				return
			}
			conn.Write(val)
		default:
			// TODO: log
			conn.Write([]byte("Error: Unsupported Command."))
		}
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

func handleGET(c *Cache, input []byte) ([]byte, error) {
	// TODO: sanitize key
	key := string(input)

	return c.Get(key)
}

func handleDELETE(c *Cache, input []byte) error {
	// TODO: sanitize key
	key := string(input)

	return c.Delete(key)
}

func handleEXISTS(c *Cache, input []byte) ([]byte, error) {
	// TODO: sanitize input
	key := string(input)

	if c.Exists(key) {
		return []byte("TRUE"), nil
	}

	return []byte("FALSE"), nil
}
