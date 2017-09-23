package roc

import (
	"errors"
	"fmt"
)

/* Note:
 *  This error code list is valid for FlashPACs, FloBoss 500-Series,
 *  FloBoss 100-Series, and RegFlo.
 *  ROCPACs and FloBoss 407 have a more complicated error set, which
 *  depends on a combination of the Opcode and Error code
 */

const (
	_ byte = iota
	InvalidOpcode
	InvalidParameterNumber
	InvalidLogicalNumber
	InvalidPointType
	TooManyBytes
	TooFewBytes
	DidNotReceive1Byte
	DidNotReceive2Byte
	DidNotReceive3Byte
	DidNotReceive4Byte
	DidNotReceive5Byte
	DidNotReceive16Byte
	OutsideValidAddressRange
	InvalidHistoryRequest
	InvalidFSTRequest
	InvalidEventEntry
	RequestedTooManyAlarms
	RequestedTooManyEvents
	WriteReadOnlyParameter
	SecurityError
	InvalidSecurityLogon
	InvalidStoreAndForwardPath
	FlashProgrammingError
	HistoryConfigurationInProgress

	RequestedSecurityLevelTooHigh byte = 63
)

type ErrorResponseMessage struct {
	Message
}

func (err ErrorResponseMessage) GetErrorCode() byte {
	return err.Data[0]
}

func (err ErrorResponseMessage) GetCauseOpcode() byte {
	return err.Data[1]
}

func (err ErrorResponseMessage) GetCauseByte() byte {
	return err.Data[2]
}

func (err ErrorResponseMessage) Error() string {
	return fmt.Sprintf("Error code %d caused by byte %d of request with opcode %d", err.GetErrorCode(), err.GetCauseByte(), err.GetCauseOpcode())
}

func errorResponseMessage(message Message) (errResponse ErrorResponseMessage, err error) {
	if message.Opcode != ErrorResponse {
		err = errors.New("Incorrect Opcode")
	} else if len(message.Data) != 3 {
		err = errors.New("Incorrect Data Length")
	} else {
		errResponse.Message = message
	}
	return
}
