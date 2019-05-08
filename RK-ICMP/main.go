package main

import (
	"flag"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"golang.org/x/net/trace"
	"net"
	"net/http"
	"os"
	"time"
)

func ping(host string) {
	p := fastping.NewPinger()

	ra, err := net.ResolveIPAddr("ip4:icmp", host)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}

	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func traceHandler(w http.ResponseWriter, req *http.Request) {
	tr := trace.New("mypkg.Foo", req.URL.Path)
	defer tr.Finish()
}

func main() {
	var (
		host  = flag.String("h", "localhost", "Host to ping")
		count = flag.Int("c", 5, "Number of ICMP packets")
	)

	flag.Parse()

	for i := 0; i < *count; i++ {
		go ping(*host)
	}

	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
}
