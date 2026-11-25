package main

import (
	"bufio"
	"net"
	"testing"
)

func TestSETCmd(t *testing.T) {
	cache := NewCache()
	go Start(cache)

	// Connect to server
	conn, err := net.Dial("tcp", "localhost:6368")
	if err != nil {
		t.Fatal(err)
	}

	// Send SET key value
	_, err = conn.Write([]byte("SET name \"John Doe\""))
	if err != nil {
		t.Fatal(err)
	}

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatal(err)
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
	t.Error("Not implemented.")
}

func TestDeleteCmd(t *testing.T) {
	t.Error("Not implemented.")
}

func TestExistsCmd(t *testing.T) {
	t.Error("Not implemented")
}
