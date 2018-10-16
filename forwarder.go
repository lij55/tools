package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	_ "os"
	_ "sync"
	"time"
	"flag"
	"os"
)

type ErrorCallback func(addr string, err error)

func logerr(addr string, err error) {
	fmt.Println(addr)
	fmt.Println(err)
}

type Forwarder struct {
	addr      string
	remotes   []string
	reporterr ErrorCallback
}

func NewForwarder(listenaddr string, toaddrs []string, errorcb ErrorCallback) *Forwarder {
	var ret = Forwarder{}
	if errorcb == nil {
		errorcb = logerr
	}
	ret.reporterr = errorcb
	ret.addr = listenaddr
	ret.remotes = toaddrs
	return &ret
}

func (fwd *Forwarder) forward(l net.Conn, raddr string) {
	r, err := net.Dial("tcp", raddr)
	if r == nil {
		fwd.reporterr(raddr, err)
		return
	}
	go io.Copy(l, r)
	go io.Copy(r, l)
}

func (fwd *Forwarder) Start() error {
	local, err := net.Listen("tcp", fwd.addr)
	if local == nil {
		return err
	}
	for {
		remote_no := len(fwd.remotes)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		conn, err := local.Accept()
		if conn == nil {
			fmt.Println(err)
			return err
		}

		go fwd.forward(conn, fwd.remotes[r.Intn(remote_no)])
	}
	return err
}



func main() {
	listen := flag.String("f", "", "address to listen on")
	dest := flag.String("d", "", "the destnation")
	flag.Parse()

	if *listen == "" {
		fmt.Fprintf(os.Stderr,"listening address is required\n")
		os.Exit(2)
	}
	if *dest == "" {
		fmt.Fprintf(os.Stderr,"dest address is required\n")
		os.Exit(2)
	}
	fmt.Fprintf(os.Stdout,"from %s to %s\n", *listen, *dest)
	remoteaddrs := []string{*dest}
	f := NewForwarder(*listen, remoteaddrs, nil)
	f.Start()
}