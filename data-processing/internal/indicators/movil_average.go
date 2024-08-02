package indicators

func WMA(prices []float64, period int) float64 {
    var weightedSum, weightSum float64
    for i := 0; i < period; i++ {
        weight := float64(period - i)
        weightedSum += prices[len(prices)-1-i] * weight
        weightSum += weight
    }
    return weightedSum / weightSum
}