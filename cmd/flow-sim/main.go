package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"math/rand"
	"time"

	"github.com/abferm/roc"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("")

var start = roc.TLP{
	PointType:   17,
	LogicNumber: 1,
	Parameter:   2,
}

const count uint8 = 3

func main() {
	netAddr := flag.String("ip", "127.0.0.1", "Network address of controller")
	port := flag.Int("p", 4000, "TCP port")
	timeout := flag.Duration("o", time.Second*2, "Response timeout")

	hostGroup := flag.Uint("hg", uint(1), "Group ID for the host to use (0-255)")
	hostUnit := flag.Uint("hu", uint(1), "Unit ID for the host to use (0-255)")
	controllerGroup := flag.Uint("cg", uint(roc.DefaultGroup), "Group ID for the host to use (0-255)")
	controllerUnit := flag.Uint("cu", uint(roc.DefaultUnit), "Unit ID for the host to use (0-255)")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	rand.Seed(time.Now().UnixNano())

	diffP := randBetween(0, 10)
	staticP := randBetween(10, 20)
	temp := randBetween(15, 120)

	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.LittleEndian, diffP)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buff, binary.LittleEndian, staticP)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buff, binary.LittleEndian, temp)
	if err != nil {
		panic(err)
	}

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	err = client.SetContiguousParameters(start, count, buff.Bytes())
	if err != nil {
		panic(err)
	}
}

func randBetween(min, max float32) (randF float32) {
	randF = min + (rand.Float32() * (max - min))
	return
}
