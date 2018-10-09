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

	hostGroup := flag.Uint("hg", uint(1), "Group ID for the host to use (0-255)")
	hostUnit := flag.Uint("hu", uint(1), "Unit ID for the host to use (0-255)")
	controllerGroup := flag.Uint("cg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	controllerUnit := flag.Uint("cu", uint(roc.DefaultGroup), "Unit ID for the host to use (0-255)")
	set := flag.Bool("set", false, "Set the time rather than read")
	zone := flag.String("tz", "Local", "Time zone of the controller")
	broadcast := flag.Bool("broadcast", false, "Broadcast time to controller group")
	flag.Parse()

	loc, err := time.LoadLocation(*zone)
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)

	if *set {
		if *broadcast {
			client := roc.NewBroadcastClientTCP(byte(*controllerGroup), *netAddr, *port, *timeout)
			err = client.SetTimeAndDate(now)
		} else {
			err = client.SetTimeAndDate(now)
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Set time to %s\n", now.String())
	} else {
		date, err := client.SendTimeAndDate(loc)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Host:       %s\n", now.String())
		fmt.Printf("Controller: %s\n", date.String())
	}

}
