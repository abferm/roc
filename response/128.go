package response

type Opcode128 struct {
	HistoryPointNumber                           uint8
	Month                                        uint8
	Day                                          uint8
	Hour                                         uint8
	Minute                                       uint8
	DatabasePointer                              uint16
	HourlyValues                                 [24]float32
	DailyArchivedValue                           float32
	MinValue                                     float32
	MaxValue                                     float32
	MinSec, MinMinute, MinHour, MinDay, MinMonth uint8
	MaxSec, MaxMinute, MaxHour, MaxDay, MaxMonth uint8
	DatabasePointType                            uint8
	_                                            float32
	_                                            float32
}
