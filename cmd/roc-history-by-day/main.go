package main

import (
	"flag"
	"time"

	"fmt"

	"github.com/abferm/roc"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("")

func main() {
	netAddr := flag.String("ip", "127.0.0.1", "Network address of controller")
	port := flag.Int("p", 4000, "TCP port")
	timeout := flag.Duration("o", time.Second, "Response timeout")

	hostGroup := flag.Uint("hg", uint(1), "Group ID for the host to use (0-255)")
	hostUnit := flag.Uint("hu", uint(1), "Unit ID for the host to use (0-255)")
	controllerGroup := flag.Uint("cg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	controllerUnit := flag.Uint("cu", uint(roc.DefaultUnit), "Unit ID for the host to use (0-255)")

	historyPoint := flag.Uint("point", 0, "History point number (254 for timestamp)")
	month := flag.Uint("month", 1, "Month to request")
	day := flag.Uint("day", 1, "Day to request")

	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	resp, err := client.SendArchivedHistoryFromDate(byte(*historyPoint), byte(*day), byte(*month))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

}
