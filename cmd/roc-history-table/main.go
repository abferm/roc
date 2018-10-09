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
	controllerUnit := flag.Uint("cu", uint(roc.DefaultGroup), "Unit ID for the host to use (0-255)")

	daily := flag.Bool("daily", false, "Pull daily archives rather than hourly")
	// historyPoint := flag.Uint("point", 0, "History point number")
	// index := flag.Uint("index", 0, "Starting history pointer")
	count := flag.Uint("count", 1, "Number of periods to request")
	pointCount := flag.Uint("pointCount", 8, "Number of history points to read")

	// ts := flag.Bool("timestamp", false, "Decode as timestamp rather than floating point value")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	client := roc.NewClientTCP(roc.Address{Group: byte(*hostGroup), Unit: byte(*hostUnit)}, roc.Address{Group: byte(*controllerGroup), Unit: byte(*controllerUnit)}, *netAddr, *port, *timeout)
	eventPointers, err := client.SendEventPointers()
	if err != nil {
		panic(err)
	}

	fmt.Println("Event Pointers")
	fmt.Printf("Alarm Log Pointer:               %d\n", eventPointers.AlarmLogPointer)
	fmt.Printf("Event Log Pointer:               %d\n", eventPointers.EventLogPointer)
	fmt.Printf("Hourly Pointer (RAM0):           %d\n", eventPointers.HourlyPointerRAM0)
	fmt.Printf("Hourly Pointer (RAM1):           %d\n", eventPointers.HourlyPointerRAM1)
	fmt.Printf("Hourly Pointer (RAM2):           %d\n", eventPointers.HourlyPointerRAM2)
	fmt.Printf("Daily Pointer (RAM0):            %d\n", eventPointers.DailyPointerRAM0)
	fmt.Printf("Daily Pointer (RAM1):            %d\n", eventPointers.DailyPointerRAM1)
	fmt.Printf("Daily Pointer (RAM2):            %d\n", eventPointers.DailyPointerRAM2)
	fmt.Printf("Max Alarm Count:                 %d\n", eventPointers.MaxAlarms)
	fmt.Printf("Max Event Count:                 %d\n", eventPointers.MaxEvents)
	fmt.Printf("Max Days (RAM0):                 %d\n", eventPointers.MaxDaysRAM0)
	fmt.Printf("Max Days (RAM1):                 %d\n", eventPointers.MaxDaysRAM1)
	fmt.Printf("Max Days (RAM2):                 %d\n", eventPointers.MaxDaysRAM2)
	fmt.Printf("(Not Canada) Minutes Per Period: %d\n", eventPointers.AuditPointer)
	fmt.Printf("(Canada) Audit Log Pointer:      %d\n\n\n", eventPointers.AuditPointer)

	index := eventPointers.HourlyPointerRAM0 - uint16(*count)
	if *daily {
		index = (uint16(eventPointers.DailyPointerRAM0) + 840) - uint16(*count)
		if index < 840 {
			index += uint16(eventPointers.MaxDaysRAM0)
		}
	} else if index >= 840 {
		index -= 840
	}

	tsResp, err := client.SendArchivedHistoryFromPointer(index, 0, 254, uint8(*count))
	if err != nil {
		panic(err)
	}

	timestamps, err := tsResp.AsTimestamps()
	if err != nil {
		panic(err)
	}

	dataValues := [][]float32{}
	header := "timestamp\t"
	for i := 0; i < int(*pointCount); i++ {
		resp, err := client.SendArchivedHistoryFromPointer(index, 0, uint8(i), uint8(*count))
		if err != nil {
			panic(err)
		}
		respValues, err := resp.AsFloats()
		if err != nil {
			panic(err)
		}
		dataValues = append(dataValues, respValues)
		header += (fmt.Sprint(i) + "\t")
	}

	fmt.Println(header)
	for row := 0; row < int(tsResp.ValueCount); row++ {
		rowString := timestamps[row].String() + "\t"
		for _, column := range dataValues {
			rowString += (fmt.Sprint(column[row]) + "\t")
		}
		fmt.Println(rowString)
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
