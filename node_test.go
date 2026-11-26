package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

const addr = "localhost:6368"

// startServer starts the cache store, and the server
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

// receive is used to receive data from server
func receive(t *testing.T, conn net.Conn) ([]byte, error) {
	t.Helper()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// send is used to send data to server
func send(t *testing.T, conn net.Conn, data []byte) error {
	t.Helper()

	_, err := conn.Write(data)
	if err != nil {
		return err
	}

	return nil
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

	// TODO: compare results
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
	cache := startServer(t)

	KEY := "location"
	VAL := "Vice city"

	cache.Set(KEY, []byte(VAL))

	conn := connect(t)

	_, err := fmt.Fprintf(conn, "DEL %s", KEY)
	if err != nil {
		t.Fatal(err)
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		t.Fatal("Error while receiving data. ", err.Error())
	}

	if cache.Exists(KEY) {
		t.Error("Deletion failed, key still exists.")
	}
}

func TestExistsCmd(t *testing.T) {
	cache := startServer(t)
	cache.Set("user:100", []byte("John"))

	conn := connect(t)

	// check for user:100
	// should return 1
	if err := send(t, conn, []byte("EXISTS user:100")); err != nil {
		t.Fatal(err)
	}

	buffer, err := receive(t, conn)
	if err != nil {
		t.Fatal("Error while receiving data from server.", err)
	}

	if string(buffer) != "TRUE" {
		t.Errorf("Unexpected value for 'EXISTS user:100'. Expected: %s, Got: %s", "1", string(buffer))
	}

	// check for user:200
	// should return 0

	if err := send(t, conn, []byte("EXISTS user:200")); err != nil {
		t.Fatal(err)
	}

	buffer, err = receive(t, conn)
	if err != nil {
		t.Fatal("Error while receiving data from server. ", err)
	}

	if string(buffer) != "FALSE" {
		t.Errorf("Unexpected value for 'EXISTS user:200'. Expected: %s, Got: %s", "0", string(buffer))
	}
}
