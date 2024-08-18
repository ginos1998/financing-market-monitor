package stock_data

import (
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka/producers"

	"github.com/robfig/cron/v3"
)

func InitHistStockDataProducer(producer *producers.KafkaProducer, server server.Server) {
	server.Logger.Info("Cron <updateHistoricalStockData> created. Schedule: every day at 9:10 AM")
	useYahooAPI := true
	c := cron.New()
	_, err := c.AddFunc("10 9 * * *", // every day at 9:10 AM
		func() {
			server.Logger.Info("Cron <updateHistoricalStockData> started at ", time.Now().Format(time.RFC3339))
			producer.UpdateHistoricalStockData(useYahooAPI, server)
			server.Logger.Info("Cron <updateHistoricalStockData> finished at ", time.Now().Format(time.RFC3339))
		})
	if err != nil {
		server.Logger.Error("Error creating cron <updateHistoricalStockData>: ", err)
		return
	}
	c.Start()
}
