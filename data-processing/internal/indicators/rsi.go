package indicators

const defaultPeriod = 14

// RSI calculates the RSI of a given set of prices
func RSI(prices []float64, period int) float64 {
	if prices == nil || len(prices) == 0 {
		return -1
	}

	if period == 0 {
		period = defaultPeriod
	}

	// initial gains and losses
	var gainSum, lossSum float64
	for i := 1; i <= period; i++ {
		change := prices[i-1] - prices[i]
		if change > 0 {
			gainSum += change
		} else {
			lossSum -= change
		}
	}

	// SMA of gains and losses
	avgGain := gainSum / float64(period)
	avgLoss := lossSum / float64(period)

	// calculate RSI for the rest of the prices
	for i := period + 1; i < len(prices); i++ {
		change := prices[i-1] - prices[i]
		gain := 0.0
		loss := 0.0
		if change > 0 {
			gain = change
		} else {
			loss = -change
		}

		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)
	}

	// calculate RSI
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}
