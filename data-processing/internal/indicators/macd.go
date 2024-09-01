package indicators

// MACD calculates the MACD line, Signal line, and MACD histogram
func MACD(prices []float64) ([]float64, []float64, []float64) {
	fastPeriod := 12
	slowPeriod := 26
	signalPeriod := 9

	fastEMA := EMA(prices, fastPeriod)
	slowEMA := EMA(prices, slowPeriod)

	macdLine := make([]float64, len(fastEMA)-len(slowEMA))
	for i := range macdLine {
		macdLine[i] = fastEMA[i+len(fastEMA)-len(slowEMA)] - slowEMA[i]
	}

	signalLine := EMA(macdLine, signalPeriod)

	macdHistogram := make([]float64, len(macdLine)-len(signalLine))
	for i := range macdHistogram {
		macdHistogram[i] = macdLine[i+len(macdLine)-len(signalLine)] - signalLine[i]
	}

	return macdLine, signalLine, macdHistogram
}
