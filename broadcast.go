package roc

import (
	"time"

	"github.com/abferm/serial"
)

// BroadcastClient broadcasts a command to the specified group

type BroadcastClient struct {
	Host       Address
	Controller Address
	Transport  Transport
}

func NewBroadcastClient(group byte, transport Transport) *BroadcastClient {
	client := new(BroadcastClient)
	client.Host.Group, client.Host.Unit = BroadcastUnit, BroadcastUnit
	client.Controller.Group, client.Controller.Unit = group, BroadcastUnit
	client.Transport = transport
	return client
}

func NewBroadcastClientSerial(group byte, port *serial.LockingSerialPort, config serial.Config) *BroadcastClient {
	return NewBroadcastClient(group, NewSerialTransport(port, config))
}

func NewBroadcastClientTCP(group byte, networkAddress string, port int, timeout time.Duration) *BroadcastClient {
	return NewBroadcastClient(group, NewTCPTransport(networkAddress, port, timeout))
}

func (client BroadcastClient) SetTimeAndDate(now time.Time) (err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SetTimeAndData
	request.Data = []byte{byte(now.Second()), byte(now.Minute()), byte(now.Hour()), byte(now.Day()), byte(now.Month()), byte(now.Year() % 100)}

	err = client.Transport.Broadcast(request)
	return
}
