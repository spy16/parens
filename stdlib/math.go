package stdlib

// Add returns sum of all the arguments.
func Add(vals ...float64) float64 {
	sum := 0.0
	for _, val := range vals {
		sum += val
	}

	return sum
}

// Sub returns result of subtracting from left-to-right.
func Sub(vals ...float64) float64 {
	if len(vals) == 1 {
		if vals[0] == 0 {
			return 0
		}
		return -1 * vals[0]
	}

	for i := 1; i < len(vals); i++ {
		vals[i] = -1 * vals[i]
	}

	return Add(vals...)
}

// Mul multiplies all numbers.
func Mul(vals ...float64) float64 {
	result := 1.0
	for _, val := range vals {
		result = result * val
	}

	return result
}

// Div divides from left to right.
func Div(vals ...float64) float64 {
	for i := range vals {
		vals[i] = 1.0 / vals[i]
	}

	return Mul(vals...)
}
