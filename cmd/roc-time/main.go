package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/abferm/roc"
)

func main() {
	netAddr := flag.String("ip", "127.0.0.1", "Network address of controller")
	port := flag.Int("p", 4000, "TCP port")
	timeout := flag.Duration("o", time.Second, "Response timeout")

	hostGroup := flag.Uint("hg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	hostUnit := flag.Uint("hu", uint(roc.DefaultUnit), "Unit ID for the host to use (0-255)")
	controllerGroup := flag.Uint("cg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	controllerUnit := flag.Uint("cu", uint(roc.BroadcastUnit), "Unit ID for the host to use (0-255)")
	flag.Parse()

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	date, err := client.SendTimeAndDate()
	if err != nil {
		panic(err)
	}

	fmt.Println(date.String())
}
