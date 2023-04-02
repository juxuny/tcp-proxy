package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func randPort() uint32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func (r *run) getClientAddress(conn net.Conn) (string, error) {
	if r.FromDeXun {
		header := make([]uint8, 8)
		_, err := conn.Read(header)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d.%d.%d.%d:%d", header[4], header[5], header[6], header[7], uint16(randPort())), nil
	}
	return conn.RemoteAddr().String(), nil
}
