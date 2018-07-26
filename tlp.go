package roc

import "fmt"

type TLP struct {
	PointType   uint8 `yaml:"point_type"`
	LogicNumber uint8 `yaml:"logic_number"`
	Parameter   uint8 `yaml:"parameter"`
}

func (tlp TLP) String() string {
	return fmt.Sprintf("%d.%d.%d", tlp.PointType, tlp.LogicNumber, tlp.Parameter)
}
