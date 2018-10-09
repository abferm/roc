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

	historyType := flag.Uint("type", 0, "History type (0-periodic, 1-extended)")
	historyPoint := flag.Uint("point", 0, "History point number (254 for timestamp)")
	index := flag.Uint("index", 0, "Starting history pointer")
	count := flag.Uint("count", 1, "Number of values to request")

	ts := flag.Bool("timestamp", false, "Decode as timestamp rather than floating point value")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	resp, err := client.SendArchivedHistoryFromPointer(uint16(*index), byte(*historyType), byte(*historyPoint), byte(*count))
	if err != nil {
		panic(err)
	}

	fmt.Printf("History Type:          %d\n", resp.HistoryType)
	fmt.Printf("History Point:         %d\n", resp.HistoryPointNumber)
	fmt.Printf("Number of Values Sent: %d\n\n", resp.ValueCount)

	if *ts {
		values, err := resp.AsTimestamps()
		if err != nil {
			panic(err)
		}
		for i, ts := range values {
			fmt.Printf("%d : %s\n", i, ts.String())
		}
	} else {
		values, err := resp.AsFloats()
		if err != nil {
			panic(err)
		}
		for i, value := range values {
			fmt.Printf("%d : %f\n", i, value)
		}
	}

}

func splitBytes(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
