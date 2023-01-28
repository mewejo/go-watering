package arduino

type Command string

const (
	WATER_OFF   Command = "a"
	WATER_ON    Command = "b"
	WATER_1_ON  Command = "c"
	WATER_1_OFF Command = "d"
	WATER_2_ON  Command = "e"
	WATER_2_OFF Command = "f"
	WATER_3_ON  Command = "g"
	WATER_3_OFF Command = "h"
	WATER_4_ON  Command = "i"
	WATER_4_OFF Command = "j"
)
