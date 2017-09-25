package roc

import (
	"sync"

	"github.com/abferm/roc/serial"
)

type LockingSerialPort struct {
	sync.Mutex
	serial.Port
	Path   string
	isOpen bool
	// TODO: make config a private member and access with getters/setters
}

func NewLockingSerialPort(path string) *LockingSerialPort {
	lp := &LockingSerialPort{}
	lp.Path = path
	return lp
}

func (lp *LockingSerialPort) Connect(config serial.Config) (err error) {
	if lp.isOpen {
		return
	}

	// If the config address doesn't match lp.Path log a warning and override config.Address with lp.Path
	if config.Address != lp.Path {
		logger.Warningf("Config address does not match port path. ( %q != %q )", config.Address, lp.Path)
		config.Address = lp.Path
	}

	if lp.Port == nil {
		lp.Port, err = serial.Open(&config)
	} else {
		err = lp.Port.Open(&config)
	}
	if err == nil {
		lp.isOpen = true
	}
	return
}

func (lp *LockingSerialPort) Close() (err error) {
	if !lp.isOpen {
		return
	}
	err = lp.Port.Close()
	lp.isOpen = false
	return
}
