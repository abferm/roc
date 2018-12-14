package main

import (
	"encoding/json"
	"flag"
	"time"

	"fmt"

	"github.com/abferm/roc"
	"github.com/abferm/roc/response"
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
	startMonth := flag.Uint("start-month", 1, "First month to request")
	startDay := flag.Uint("start-day", 1, "First day to request")
	endMonth := flag.Uint("end-month", 1, "First month to request")
	endDay := flag.Uint("end-day", 1, "First day to request")

	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}
	startDate := time.Date(2018, time.Month(*startMonth), int(*startDay), 0, 0, 0, 0, time.Local)
	endDate := time.Date(2018, time.Month(*endMonth), int(*endDay), 0, 0, 0, 0, time.Local)

	responses := []*response.Opcode128{}
	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	for date := startDate; date.Unix() <= endDate.Unix(); date = date.Add(time.Hour * 24) {
		resp, err := client.SendArchivedHistoryFromDate(byte(*historyPoint), uint8(date.Day()), uint8(date.Month()))
		if err != nil {
			panic(err)
		}

		logger.Debugf("%+v", resp)
		responses = append(responses, resp)
	}
	b, err := json.MarshalIndent(responses, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}
