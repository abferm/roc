package roc

import "fmt"

type TLP struct {
	PointType   uint8
	LogicNumber uint8
	Parameter   uint8
}

func (tlp TLP) String() string {
	return fmt.Sprintf("%d.%d.%d", tlp.PointType, tlp.LogicNumber, tlp.Parameter)
}
