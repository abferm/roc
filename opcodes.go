package roc

const (
	GeneralUpdate          byte = 0
	SendTestData           byte = 2 // ROC300-Series and FloBoss 407 ONLY
	SendConfiguration      byte = 6
	SendTimeAndDate        byte = 7
	SetTimeAndData         byte = 8
	SendConfigurableOpcode byte = 10
	SetConfigurableOpcode  byte = 11
	SetOperatorID          byte = 17
	LogsEvent              byte = 18
	RESERVED_19            byte = 19
	StoreAndForward        byte = 24
	RESERVED_80            byte = 80
	RESERVED_100           byte = 100
	SetSystemVariables     byte = 102
	SendSystemInfo         byte = 103
	SendHistoryPointDef    byte = 105
	SendHistory            byte = 107
	SendEventPointers      byte = 120
	SendAlarms             byte = 121
	SendEvents             byte = 122
	ReadUserTemplate       byte = 123
)
