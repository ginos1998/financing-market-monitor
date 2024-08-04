package utils

import (
	"sort"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

type ByDateDesc []dtos.TimeSeries

func (a ByDateDesc) Len() int           { return len(a) }
func (a ByDateDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDateDesc) Less(i, j int) bool {
    dateI, _ := time.Parse("2006-01-02", a[i].Date)
    dateJ, _ := time.Parse("2006-01-02", a[j].Date)
    return dateI.After(dateJ) // Comparación descendente
}


func OrderTimeSeriesByDateDesc(timeSeries []dtos.TimeSeries) []dtos.TimeSeries {
	sort.Sort(ByDateDesc(timeSeries))
	return timeSeries
}

func OrderTimeSeriesByDateAsc(timeSeries *[]dtos.TimeSeries) []dtos.TimeSeries {
	sort.Slice(*timeSeries, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", (*timeSeries)[i].Date)
		dateJ, _ := time.Parse("2006-01-02", (*timeSeries)[j].Date)
		return dateI.Before(dateJ)
	})
	return *timeSeries
}