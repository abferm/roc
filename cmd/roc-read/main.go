package main

import (
	"flag"
	"time"

	"fmt"

	"strconv"
	"strings"

	"github.com/abferm/roc"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("")

func main() {
	netAddr := flag.String("ip", "127.0.0.1", "Network address of controller")
	port := flag.Int("p", 4000, "TCP port")
	timeout := flag.Duration("o", time.Second, "Response timeout")

	hostGroup := flag.Uint("hg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	hostUnit := flag.Uint("hu", uint(roc.DefaultUnit), "Unit ID for the host to use (0-255)")
	controllerGroup := flag.Uint("cg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	controllerUnit := flag.Uint("cu", uint(roc.BroadcastUnit), "Unit ID for the host to use (0-255)")

	tlpString := flag.String("tlps", "0.0.0", "Comma separated list of TLPs to read")
	ascii := flag.Bool("ascii", false, "Encode response bytes as ascii rather than hex")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	tlps, err := parseTLPList(*tlpString)
	if err != nil {
		panic(err)
	}

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	data, err := client.SendSpecifiedParameters(tlps)
	if err != nil {
		panic(err)
	}

	if *ascii {
		fmt.Println(string(data))
	} else {
		fmt.Printf("%x\n", data)
	}
}

func parseTLPList(tlpsString string) (tlps []roc.TLP, err error) {
	for _, tlpString := range strings.Split(tlpsString, ",") {
		var tlp roc.TLP
		tlp, err = parseTLP(tlpString)
		if err != nil {
			return
		}
		tlps = append(tlps, tlp)
	}
	return
}

func parseTLP(tlpString string) (tlp roc.TLP, err error) {
	separated := strings.Split(tlpString, ".")
	if len(separated) != 3 {
		err = fmt.Errorf("Invalid TLP %q", tlpString)
		return
	}

	// Using a loop to avoid code duplication
	for i, s := range separated {
		var val uint64
		val, err = strconv.ParseUint(s, 10, 8)
		if err != nil {
			err = fmt.Errorf("Invalid TLP %q: %s", tlpString, err.Error())
			return
		}
		switch i {
		case 0:
			tlp.PointType = uint8(val)
		case 1:
			tlp.LogicNumber = uint8(val)
		case 2:
			tlp.Parameter = uint8(val)
		}
	}

	return
}
