package main

import (
	"context"
	"fmt"
	"github.com/juxuny/yc/cmd"
	"github.com/juxuny/yc/log"
	"github.com/spf13/cobra"
	"net"
	"strings"
	"time"
)

type run struct {
	ToProxyProtocol bool
	FromDeXun       bool
	Listen          string
	Remote          string
	Timeout         int
}

func (r *run) Prepare(cmd *cobra.Command) {
}

func (r *run) InitFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&r.ToProxyProtocol, "to-proxy-protocol", false, "enable proxy_protocol to nginx")
	cmd.PersistentFlags().BoolVar(&r.FromDeXun, "from-de-xun", false, "enable dexunyun proxy, ref: https://www.dexunyun.com/")
	cmd.PersistentFlags().StringVarP(&r.Listen, "listen", "l", ":20000", "listen address")
	cmd.PersistentFlags().StringVarP(&r.Remote, "remote", "r", "127.0.0.1:10000", "backend nginx port")
	cmd.PersistentFlags().IntVarP(&r.Timeout, "timeout", "t", 60, "timeout")
}

const bufferLen = 10240

func (r run) transfer(ctx context.Context, cancel context.CancelFunc, from net.Conn, to net.Conn) {
	defer func() {
		_ = from.Close()
		_ = to.Close()
	}()
	buf := make([]byte, bufferLen)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		_ = from.SetDeadline(time.Now().Add(time.Second * time.Duration(r.Timeout)))
		n, err := from.Read(buf)
		if err != nil {
			log.Debug(err)
			cancel()
			return
		}
		_ = to.SetDeadline(time.Now().Add(time.Second * time.Duration(r.Timeout)))
		_, err = to.Write(buf[:n])
		if err != nil {
			log.Error(err)
			cancel()
			return
		}
	}
}

func (r run) sendClientAddress(clientAddress string, conn net.Conn) error {
	s := strings.Split(clientAddress, ":")
	d := strings.Split(conn.RemoteAddr().String(), ":")
	buf := []byte(fmt.Sprintf("PROXY TCP4 %s %s %s %s\r\n", s[0], d[0], s[1], d[1]))
	log.Debug(string(buf))
	_, err := conn.Write(buf)
	return err
}

func (r run) handleClient(client net.Conn) {
	clientAddress, err := r.getClientAddress(client)
	if err != nil {
		log.Error(err)
		return
	}
	backendConn, err := net.Dial("tcp", r.Remote)
	if err != nil {
		log.Error(err)
		_ = client.Close()
		return
	}
	if r.ToProxyProtocol {
		if err := r.sendClientAddress(clientAddress, backendConn); err != nil {
			log.Error(err)
			_ = client.Close()
			_ = backendConn.Close()
			return
		}
	}
	log.Info("accepted from:", clientAddress)
	ctx, cancel := context.WithCancel(context.Background())
	go r.transfer(ctx, cancel, client, backendConn)
	go r.transfer(ctx, cancel, backendConn, client)
}

func (r run) start() {
	ln, err := net.Listen("tcp", r.Listen)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		client, err := ln.Accept()
		if err != nil {
			log.Error(err)
			break
		}
		go r.handleClient(client)
	}
}

func (r run) Run() {
	for {
		r.start()
		time.Sleep(time.Second * 5)
	}
}

func main() {
	runCommand := cmd.NewCommandBuilder("", &run{})
	if err := runCommand.Build().Execute(); err != nil {
		log.Error(err)
	}
}
