package roc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/abferm/roc/response"
	"github.com/abferm/serial"
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

func NewClientSerial(host, controller Address, port *serial.LockingSerialPort, config serial.Config) *Client {
	return NewClient(host, controller, NewSerialTransport(port, config))
}

func NewClientTCP(host, controller Address, networkAddress string, port int, timeout time.Duration) *Client {
	return NewClient(host, controller, NewTCPTransport(networkAddress, port, timeout))
}

func (client Client) SendTimeAndDate(loc *time.Location) (now time.Time, err error) {
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
	// Assume it's at in the 2000s
	year := int(response.Data[5]) + 2000
	//leapYear := int(response.Data[6])
	//weekDay := int(response.Data[7])
	now = time.Date(year, month, day, hour, minute, second, 0, loc)
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

func (client Client) SendSpecifiedParameters(parameters []TLP) (data []byte, err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SendSpecifiedParameters

	if len(parameters) > 255 {
		err = fmt.Errorf("Too many parameters requested.")
		return
	}

	request.Data = []byte{byte(len(parameters))}

	for _, tlp := range parameters {
		request.Data = append(request.Data, tlp.PointType, tlp.LogicNumber, tlp.Parameter)
	}

	response, err := client.Transport.Transceive(request)
	if err != nil {
		return
	}

	data = response.Data
	return
}

func (client Client) SetTimeAndDate(now time.Time) (err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SetTimeAndDate
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

func (client Client) SendMultipleHistoryPoints(segment uint8, index uint16, historyType, startingHistoryPoint, pointCount, periodCount uint8) (data []byte, err error) {
	request := Message{}
	request.Source, request.Destination = client.Host, client.Controller
	request.Opcode = SendMultipleHistoryPoints
	indexBytes := []byte{0, 0}
	binary.LittleEndian.PutUint16(indexBytes, index)
	request.Data = []byte{segment, indexBytes[0], indexBytes[1], historyType, startingHistoryPoint, pointCount, periodCount}

	response, err := client.Transport.Transceive(request)
	data = response.Data
	return
}

func (client Client) SendArchivedHistoryFromPointer(index uint16, historyType, historyPoint, periodCount uint8) (resp *response.Opcode130, err error) {
	resp = new(response.Opcode130)
	requestMSG := Message{}
	requestMSG.Source, requestMSG.Destination = client.Host, client.Controller
	requestMSG.Opcode = SendArchivedHistoryFromPointer
	indexBytes := []byte{0, 0}
	binary.LittleEndian.PutUint16(indexBytes, index)
	requestMSG.Data = []byte{historyType, historyPoint, periodCount, indexBytes[0], indexBytes[1]}

	responseMSG, err := client.Transport.Transceive(requestMSG)
	if err != nil {
		return
	}
	err = resp.FromData(responseMSG.Data)
	return
}

func (client Client) SendEventPointers() (resp *response.Opcode120, err error) {
	resp = new(response.Opcode120)
	requestMSG := Message{}
	requestMSG.Source, requestMSG.Destination = client.Host, client.Controller
	requestMSG.Opcode = SendEventPointers

	responseMSG, err := client.Transport.Transceive(requestMSG)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(responseMSG.Data), binary.LittleEndian, resp)
	return
}

func (client Client) SendArchivedHistoryFromDate(historyPoint, day, month uint8) (resp *response.Opcode128, err error) {
	resp = new(response.Opcode128)
	requestMSG := Message{}
	requestMSG.Source, requestMSG.Destination = client.Host, client.Controller
	requestMSG.Opcode = SendArchivedHistoryFromDate
	requestMSG.Data = []byte{historyPoint, day, month}

	responseMSG, err := client.Transport.Transceive(requestMSG)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(responseMSG.Data), binary.LittleEndian, resp)
	return
}
