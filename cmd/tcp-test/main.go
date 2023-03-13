package main

import (
	"context"
	"fmt"
	"github.com/juxuny/yc/log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":20001")
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		defer ln.Close()
		client, err := ln.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		_, err = client.Write([]byte("1234567890"))
		if err != nil {
			log.Error(err)
		}
		<-context.Background().Done()
	}()
	conn, err := net.Dial("tcp", ":20001")
	if err != nil {
		log.Error(err)
		return
	}
	buf := make([]byte, 2)
	_, err = conn.Read(buf)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(buf))
	//conn.Close()
	<-context.Background().Done()
}
