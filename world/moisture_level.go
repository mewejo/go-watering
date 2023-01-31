package world

// Percentage of hysteresis to use (decimal, 0.05 = 5%)
const hysteresis = 0.05

type MoistureLevel struct {
	Percentage uint
}

func percentageToMoistureLevel(level float64) MoistureLevel {
	if level > 100 {
		level = 100
	}

	if level < 0 {
		level = 0
	}

	return MoistureLevel{
		Percentage: uint(level),
	}
}

func (ml MoistureLevel) HysteresisOnLevel() MoistureLevel {
	return percentageToMoistureLevel(
		float64(ml.Percentage) * (1.0 - hysteresis),
	)
}

func (ml MoistureLevel) HysteresisOffLevel() MoistureLevel {
	return percentageToMoistureLevel(
		float64(ml.Percentage) * (1.0 + hysteresis),
	)
}
