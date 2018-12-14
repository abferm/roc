package response

type Opcode128 struct {
	HistoryPointNumber uint8       `json:"history_point_number"`
	Month              uint8       `json:"month"`
	Day                uint8       `json:"day"`
	Hour               uint8       `json:"hour"`
	Minute             uint8       `json:"minute"`
	DatabasePointer    uint16      `json:"database_pointer"`
	HourlyValues       [24]float32 `json:"hourly_values"`
	DailyArchivedValue float32     `json:"daily_archived_value"`
	MinValue           float32     `json:"min_value"`
	MaxValue           float32     `json:"max_value"`
	MinSecond          uint8       `json:"min_second"`
	MinMinute          uint8       `json:"min_minute"`
	MinHour            uint8       `json:"min_hour"`
	MinDay             uint8       `json:"min_day"`
	MinMonth           uint8       `json:"min_month"`
	MaxSecond          uint8       `json:"max_second"`
	MaxMinute          uint8       `json:"max_minute"`
	MaxHour            uint8       `json:"max_hour"`
	MaxDay             uint8       `json:"max_day"`
	MaxMonth           uint8       `json:"max_month"`
	DatabasePointType  uint8       `json:"database_point_type"`
	_                  float32
	_                  float32
}
