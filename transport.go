package roc

import (
	"io"
	"sync"

	"fmt"
	"net"
	"time"

	"github.com/abferm/serial"
)

type Transport interface {
	Transceive(request Message) (response Message, err error)
	Broadcast(request Message) (err error)
}

type baseTransport struct {
}

func (client baseTransport) send(request Message, transport io.Writer) (err error) {
	err = request.updateCalculatedFields()
	if err != nil {
		return
	}

	logger.Debugf("sending % x\n", request.bytes())
	_, err = transport.Write(request.bytes())
	if err != nil {
		return
	}
	return
}

func (client baseTransport) receive(transport io.Reader) (response Message, err error) {
	err = response.read(transport)
	if err != nil {
		return
	}
	if response.Opcode == ErrorResponse {
		var errM ErrorResponseMessage
		errM, err = errorResponseMessage(response)
		if err == nil {
			err = errM
		}
	}
	return
}

func (client baseTransport) transceive(request Message, transport io.ReadWriter) (response Message, err error) {
	err = client.send(request, transport)
	if err != nil {
		return
	}

	response, err = client.receive(transport)
	return
}

type SerialTransport struct {
	baseTransport
	Port   *serial.LockingSerialPort
	Config serial.Config
}

func NewSerialTransport(port *serial.LockingSerialPort, config serial.Config) *SerialTransport {
	trans := new(SerialTransport)
	trans.Port = port
	trans.Config = config
	return trans
}

func (trans SerialTransport) Transceive(request Message) (response Message, err error) {
	trans.Port.Lock()
	defer trans.Port.Unlock()
	err = trans.Port.Connect(trans.Config)
	if err != nil {
		return
	}
	defer trans.Port.Close()

	logger.Debugf("flushing serial port\n")
	if err = trans.Port.Flush(); err != nil {
		return
	}

	response, err = trans.transceive(request, trans.Port)
	return
}

func (trans SerialTransport) Broadcast(request Message) (err error) {
	trans.Port.Lock()
	defer trans.Port.Unlock()
	err = trans.Port.Connect(trans.Config)
	if err != nil {
		return
	}
	defer trans.Port.Close()

	logger.Debugf("flushing serial port\n")
	if err = trans.Port.Flush(); err != nil {
		return
	}

	err = trans.send(request, trans.Port)
	if err != nil {
		return
	}

	// Sleep 1/2 timeout in order to ensure the message makes it out before the
	// next flush
	time.Sleep(trans.Config.Timeout / 2)
	return
}

type TCPTransport struct {
	baseTransport
	Address string
	Port    int
	Timeout time.Duration
}

func NewTCPTransport(address string, port int, timeout time.Duration) *TCPTransport {
	trans := new(TCPTransport)
	trans.Address = address
	trans.Port = port
	trans.Timeout = timeout
	return trans
}

func (trans TCPTransport) Transceive(request Message) (response Message, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", trans.Address, trans.Port))
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.SetDeadline(time.Now().Add(trans.Timeout))
	if err != nil {
		return
	}
	response, err = trans.transceive(request, conn)
	return
}

func (trans TCPTransport) Broadcast(request Message) (err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", trans.Address, trans.Port))
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.SetDeadline(time.Now().Add(trans.Timeout))
	if err != nil {
		return
	}
	err = trans.send(request, conn)
	return
}

// SharedTCPTransport: This is a TCPTransport wrapped in a mutex lock. It should
// be used when multiple devices are available at the same address and port, for
// example several serial devices connected to a serial to tcp converter via RS-485
// or a radio network.
type SharedTCPTransport struct {
	TCPTransport
	lock *sync.Mutex
}

func NewSharedTCPTransport(address string, port int, timeout time.Duration) *SharedTCPTransport {
	trans := new(SharedTCPTransport)
	trans.Address = address
	trans.Port = port
	trans.Timeout = timeout
	trans.lock = new(sync.Mutex)
	return trans
}

func (trans SharedTCPTransport) Transceive(request Message) (response Message, err error) {
	trans.lock.Lock()
	response, err = trans.TCPTransport.Transceive(request)
	trans.lock.Unlock()
	return
}

func (trans SharedTCPTransport) Broadcast(request Message) (err error) {
	trans.lock.Lock()
	err = trans.TCPTransport.Broadcast(request)
	trans.lock.Unlock()
	return
}
