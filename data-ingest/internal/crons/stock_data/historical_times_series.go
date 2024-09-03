package stock_data

import (
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka/producers"

	"github.com/robfig/cron/v3"
)

func InitHistStockDataProducer(producer *producers.KafkaProducer, server server.Server) {
	server.Logger.Info("Cron <UpdateTickersTimeSeries> created. Schedule: every day at 9:10 AM")
	c := cron.New()
	_, err := c.AddFunc("21 21 * * *", // every day at 9:10 AM
		func() {
			server.Logger.Info("Cron <UpdateTickersTimeSeries> started at ", time.Now().Format(time.RFC3339))
			producer.UpdateTickersTimeSeries(server)
			server.Logger.Info("Cron <UpdateTickersTimeSeries> finished at ", time.Now().Format(time.RFC3339))
		})
	if err != nil {
		server.Logger.Error("Error creating cron <UpdateTickersTimeSeries>: ", err)
		return
	}
	c.Start()
}
