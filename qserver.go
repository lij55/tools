package main

import (
	"net"
	"log"
	"fmt"
	"net/http"
	"os"
	"io"
	"flag"
)

func ShowURL(port string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Get interfaces error: %x", err)
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		// handle err
		if err != nil {
			log.Fatal("Get address error: %x", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.IsLoopback() {
				continue
			}
			// process IP address
			if ip.To4() != nil {
				fmt.Printf("http://%s:%s\n", ip.String(), port)
			}
		}
	}
}

func ShowExtIP() {
	resp, err := http.Get("http://ipv4.myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func runHTTP(port string, directory string) {
	http.Handle("/", http.FileServer(http.Dir(directory)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	port := flag.String("p", "8000", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	showip := flag.Bool("v", true, "show full url")
	flag.Parse()
	if *showip == true {
		ShowURL(*port)
	}
	runHTTP(*port, *directory)
}