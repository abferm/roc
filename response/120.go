package response

type Opcode120 struct {
	AlarmLogPointer   uint16 // offset 6 length 2
	EventLogPointer   uint16 // offset 8 length 2
	HourlyPointerRAM0 uint16 // offset 10 length 2
	HourlyPointerRAM1 uint16 // offset 12 length 2
	HourlyPointerRAM2 uint16 // offset 14 length 2
	_                 uint16 // offset 16 length 2
	DailyPointerRAM0  uint8  // offset 18 length 1
	DailyPointerRAM1  uint8  // offset 19 length 1
	DailyPointerRAM2  uint8  // offset 20 length 1
	_                 uint8  // offset 21 length 1
	MaxAlarms         uint16 // offset 22 length 2
	MaxEvents         uint16 // offset 24 length 2
	MaxDaysRAM0       uint8  // offset 26 length 1
	MaxDaysRAM1       uint8  // offset 27 length 1
	MaxDaysRAM2       uint8  // offset 28 length 1
	_                 uint8  // offset 29 length 1
	AuditPointer      uint16 // offset 30 length 2 (Canada Only)
}
