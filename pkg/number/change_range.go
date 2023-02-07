package number

// Source: https://www.arduino.cc/reference/en/language/functions/math/map/
func ChangeRange(input float64, inputMin float64, inputMax float64, outputMin float64, outputMax float64) float64 {
	return (input-inputMin)*(outputMax-outputMin)/(inputMax-inputMin) + outputMin
}
