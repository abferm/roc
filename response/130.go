package response

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Opcode130 struct {
	HistoryType        uint8  // offset 6 length 1
	HistoryPointNumber uint8  // offset 7 length 1
	ValueCount         uint8  // offset 8 length 1
	ValueData          []byte // offset 9 length 4 * ValueCount
}

func (resp *Opcode130) FromData(messageData []byte) (err error) {
	if len(messageData) < 3 {
		err = fmt.Errorf("At least 3 bytes required, received %d", len(messageData))
		return
	}
	resp.HistoryType, resp.HistoryPointNumber, resp.ValueCount = messageData[0], messageData[1], messageData[2]
	if resp.ValueCount > 0 {
		resp.ValueData = messageData[3:]
	}
	return
}

func (resp *Opcode130) AsFloats() (values []float32, err error) {
	values = make([]float32, int(resp.ValueCount))
	err = binary.Read(bytes.NewReader(resp.ValueData), binary.LittleEndian, &values)
	return
}

func (resp *Opcode130) AsTimestamps() (values []Opcode130Timestamp, err error) {
	values = make([]Opcode130Timestamp, int(resp.ValueCount))
	err = binary.Read(bytes.NewReader(resp.ValueData), binary.LittleEndian, &values)
	return
}

type Opcode130Timestamp struct {
	Minute uint8
	Hour   uint8
	Day    uint8
	Month  uint8
}

func (ts Opcode130Timestamp) String() string {
	return fmt.Sprintf("%d-%d %d:%d", ts.Month, ts.Day, ts.Hour, ts.Minute)
}

func (ts Opcode130Timestamp) ToTime(loc *time.Location) (timestamp time.Time) {
	now := time.Now()
	year := now.Year()
	if now.Month() < time.Month(ts.Month) {
		// Assume the timestamp is from last year
		// May be wrong if the ROC has been powered off for longer than one year.
		year -= 1
	}

	// Avoid panic
	if loc == nil {
		loc = time.Local
	}

	timestamp = time.Date(year, time.Month(ts.Month), int(ts.Day), int(ts.Hour), int(ts.Minute), 0, 0, loc)
	return
}
