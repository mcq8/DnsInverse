package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	dnstap "github.com/dnstap/golang-dnstap"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		tcpInputs   = flag.String("l", "127.0.0.1:1058", "read dnstap payloads from tcp/ip")
		apiAddr     = flag.String("p", ":80", "api address for example :80 or 127.0.0.1:80 ")
		flagTimeout = flag.Duration("t", 0, "I/O timeout for tcp/ip and unix domain sockets")
	)
	flag.Parse()

	inverseDnsMap := newInverseDnsMap()
	output := newInverseDnsFiller(inverseDnsMap)
	go output.RunOutputLoop()

	l, err := net.Listen("tcp", *tcpInputs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen on %s: %v\n", tcpInputs, err)
		os.Exit(1)
	}
	i := dnstap.NewFrameStreamSockInput(l)
	i.SetTimeout(*flagTimeout)
	go i.ReadInto(output.GetOutputChannel())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%q", inverseDnsMap.Get(r.URL.Query().Get("ip")))
	})
	log.Fatal(http.ListenAndServe(*apiAddr, nil))

}
