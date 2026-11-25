package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

const addr = "localhost:6368"

func startServer(t *testing.T) *Cache {
	t.Helper()

	cache := NewCache()
	l := Start(cache)

	t.Cleanup(func() {
		l.Close()
	})

	return cache
}

func connect(t *testing.T) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func TestSETCmd(t *testing.T) {
	cache := startServer(t)

	// Connect to server
	conn := connect(t)

	// Send SET key value
	_, err := conn.Write([]byte("SET name \"John Doe\""))
	if err != nil {
		t.Fatal(err)
	}

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatal("Error while reading from server. ", err)
	}

	EXPECTED_REPLY := "OK\n"
	if reply != EXPECTED_REPLY {
		t.Errorf("Unexpected reply from server. Expected: %q, Got: %q", EXPECTED_REPLY, reply)
	}

	if !cache.Exists("name") {
		t.Error("SET operation failed to store key.")
	}
}

func TestGetCmd(t *testing.T) {
	cache := startServer(t)

	KEY := "location"
	VAL := "Vice city"

	cache.Set(KEY, []byte(VAL))

	conn := connect(t)

	// send GET cmd
	_, err := fmt.Fprintf(conn, "GET %s", KEY)
	if err != nil {
		t.Error("Failed to send GET cmd to server. ", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Error("Failed to get reply from server. ", err)
	}

	if string(buffer[:n]) != VAL {
		t.Errorf("Unexpected reply from server. Expected: %q, Got: %q", VAL, string(buffer[:n]))
	}
}

func TestDeleteCmd(t *testing.T) {
	t.Error("Not implemented.")
}
