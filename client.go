package roc

import (
	"fmt"
	"time"

	"github.com/abferm/roc/serial"
)

type Client struct {
	Host       Address
	Controller Address
	Transport  Transport
}

func NewClient(host, controller Address, transport Transport) *Client {
	client := new(Client)
	client.Host, client.Controller = host, controller
	client.Transport = transport
	return client
}

func NewClientSerial(host, controller Address, port *LockingSerialPort, config serial.Config) *Client {
	return NewClient(host, controller, NewSerialTransport(port, config))
}

func NewClientTCP(host, controller Address, networkAddress string, port int, timeout time.Duration) *Client {
	return NewClient(host, controller, NewTCPTransport(networkAddress, port, timeout))
}

func (client Client) SendTimeAndDate() (now time.Time, err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SendTimeAndDate

	response, err := client.Transport.Transceive(request)
	if err != nil {
		return
	}

	if response.DataLength != 8 {
		err = fmt.Errorf("ROC date is 8 bytes, received %d", request.DataLength)
	}
	second := int(response.Data[0])
	minute := int(response.Data[1])
	hour := int(response.Data[2])
	day := int(response.Data[3])
	month := time.Month(response.Data[4])
	year := int(response.Data[5])
	//leapYear := int(response.Data[6])
	//weekDay := int(response.Data[7])
	now = time.Date(year, month, day, hour, minute, second, 0, time.Local)
	return
}

func (client Client) SendContiguousParameters(start TLP, count uint8) (data []byte, err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SendContiguousParameters
	request.Data = []byte{start.PointType, start.LogicNumber, count, start.Parameter}

	response, err := client.Transport.Transceive(request)
	if err != nil {
		return
	}

	data = response.Data[4:]
	return
}

func (client Client) SetTimeAndDate(now time.Time) (err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SetTimeAndData
	request.Data = []byte{byte(now.Second()), byte(now.Minute()), byte(now.Hour()), byte(now.Day()), byte(now.Month()), byte(now.Year() % 100)}

	_, err = client.Transport.Transceive(request)
	return
}
func (client Client) SetContiguousParameters(start TLP, count uint8, data []byte) (err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SetContiguousParameters
	request.Data = append([]byte{start.PointType, start.LogicNumber, count, start.Parameter}, data...)

	_, err = client.Transport.Transceive(request)
	return
}
