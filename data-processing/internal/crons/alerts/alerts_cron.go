package alerts

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	alertService "github.com/ginos1998/financing-market-monitor/data-processing/internal/services/alerts"
)

var logger logrus.Logger

func StartAlertsCron(s *server.Server) error {
	logger = *s.Logger
	c := cron.New()
	_, err := c.AddFunc("@every 30s", // every day at 9:10 AM
		func() {
			manageAlerts(&s.RedisClient, s.MongoRepository)
		})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func manageAlerts(redisClient *redis.RedisClient, mongoRepository mongod.MongoRepository) {
	logger.Info("ALERTS CRON | Checking alerts...")
	alerts, err := alertService.FindActiveAlerts(mongoRepository)
	if err != nil {
		logger.Error("Error getting alerts: ", err)
		return
	}

	if len(alerts) == 0 {
		logger.Info("ALERTS CRON | No active alerts")
		return
	}
	logger.Infof("ALERTS CRON | Found %d active alerts", len(alerts))

	alertsTriggered, err := alertService.ProcessAlerts(mongoRepository, redisClient, alerts)
	if err != nil {
		logger.Errorf("Error processing alerts:\n%v\nCause:%v", alerts, err)
		return
	}
	if len(alertsTriggered) > 0 {
		logger.Infof("ALERTS CRON | Alerts triggered:%d.\nIds:", len(alertsTriggered))
		for _, alert := range alertsTriggered {
			logger.Infof("Alert: %s\n", alert.ToString())
		}
	} else {
		logger.Info("ALERTS CRON | No alerts triggered")
	}

}
