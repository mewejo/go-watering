package arduino

import "errors"

type WaterOutlet int

const (
	WATER_OUTLET_1 WaterOutlet = iota
	WATER_OUTLET_2
	WATER_OUTLET_3
	WATER_OUTLET_4
)

func (o WaterOutlet) OnCommand() (Command, error) {
	if o == WATER_OUTLET_1 {
		return WATER_1_ON, nil
	} else if o == WATER_OUTLET_2 {
		return WATER_2_ON, nil
	} else if o == WATER_OUTLET_3 {
		return WATER_3_ON, nil
	} else if o == WATER_OUTLET_4 {
		return WATER_4_ON, nil
	} else {
		return WATER_OFF, errors.New("Invalid outlet")
	}
}

func (o WaterOutlet) OffCommand() (Command, error) {
	if o == WATER_OUTLET_1 {
		return WATER_1_OFF, nil
	} else if o == WATER_OUTLET_2 {
		return WATER_2_OFF, nil
	} else if o == WATER_OUTLET_3 {
		return WATER_3_OFF, nil
	} else if o == WATER_OUTLET_4 {
		return WATER_4_OFF, nil
	} else {
		return WATER_OFF, errors.New("Invalid outlet")
	}
}
