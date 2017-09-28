package roc

import (
	"bytes"
	"fmt"
)

// splitParameters: split data between the supplied TLPs
// Mainly for use with SendSpecifiedParameters
// Note: this implementation has the following known flaws
//		Case: Next TLP's bytes occur naturally within
// 				the data for the current TLP.
//		Result: Short byte slice followed by a long byte slice
//
// The surest way to do this is build a lookup table for data lengths.
// The two biggest problems with that are that it would be a pain to
// set up and there are differences between ROC devices, so you would
// need to include the device model in all look-ups.
func splitParameters(parameters []TLP, data []byte)(splitData map[TLP][]byte, err error){
	if len(data) < 4{
		err = fmt.Errorf("Too little data")
		return
	}

	if len(parameters)<1{
		err = fmt.Errorf("At least one parameter is required")
		return
	}

	// Assert the byte slice starts with the first TLP
	if !bytes.HasPrefix(data,[]byte{parameters[0].PointType, parameters[0].LogicNumber, parameters[0].Parameter}){
		err = fmt.Errorf("Missing TLP")
		return
	}

	// Make a copy of the data, dropping the TLP bytes for the first parameter
	dataCopy := append([]byte{}, data[3:]...)

	splitData = make(map[TLP][]byte)

	for i, current := range parameters{

		// We're done! all TLPs have been handled
		if i == (len(parameters)-1){
			splitData[current]=dataCopy
			return
		}

		// Split on next TLP to get current bytes, and remainder
		next := []byte{parameters[i+1].PointType, parameters[i+1].LogicNumber, parameters[i+1].Parameter}
		chunks := bytes.SplitN(dataCopy, next, 2)
		if len(chunks) != 2{
			err = fmt.Errorf("Too little data")
			return
		}else if len(chunks[1]) < 1{
			err = fmt.Errorf("Too little data")
			return
		}

		splitData[current] = chunks[0]
		dataCopy = chunks[1]
	}

	return
}
