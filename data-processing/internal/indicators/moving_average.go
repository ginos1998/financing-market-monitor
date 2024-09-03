package indicators

func WMA(prices []float64, period int) float64 {
	if prices == nil || len(prices) == 0 || period <= 0 {
		return -1
	}

	var weightedSum, weightSum float64
	for i := 0; i < period; i++ {
		weight := float64(period - i)
		weightedSum += prices[i] * weight
		weightSum += weight
	}
	return weightedSum / weightSum
}

func SMA(prices []float64, period int) float64 {
	if prices == nil || len(prices) == 0 || period <= 0 {
		return -1
	}

	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	return sum / float64(period)
}

// EMA calculates the Exponential Moving Average
func EMA(prices []float64, period int) []float64 {
	multiplier := 2.0 / float64(period+1)
	ema := make([]float64, len(prices))

	// Start with an initial SMA for the first period
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate the EMA for each subsequent period
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
	}

	return ema[period-1:]
}
