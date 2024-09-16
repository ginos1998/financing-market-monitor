package alerts

import (
	"errors"
	"fmt"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/discord"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	alertService "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/alerts"
	alertsRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/alerts"
	ds "github.com/ginos1998/financing-market-monitor/data-processing/internal/integrations/discord"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
	alertsDTO "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos/alerts"
	redisService "github.com/ginos1998/financing-market-monitor/data-processing/internal/services/redis"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger logrus.Logger

var alertTypes = map[string]func(float64, dtos.IntraDayPrices) bool{
	"crossing":      isCrossing,
	"over":          isOver,
	"under":         isUnder,
	"reaching":      isReaching,
	"crossing-up":   isCrossingUp,
	"crossing-down": isCrossingDown,
}

func FindActiveAlerts(mdb mongod.MongoRepository) ([]alertsDTO.Alert, error) {
	return alertService.FindActiveAlerts(mdb)
}

func ProcessAlerts(server *server.Server, alerts []alertsDTO.Alert) ([]alertsDTO.Alert, error) {
	logger = *server.Logger
	triggeredAlerts := make([]alertsDTO.Alert, 0)
	for _, alert := range alerts {
		symbolKey := fmt.Sprintf("INTRA-DAY_%s", alert.Symbol)
		symbolPrices, err := redisService.GetIntraDaySymbolPrices(&server.RedisClient, symbolKey)
		if err != nil || symbolPrices == "" {
			logger.Warnf("Couldn't process alert %s, symbol prices not found for symbol %s", alert.Id, alert.Symbol)
			continue
		}
		symbolValues := dtos.IntraDayPrices{}
		err = symbolValues.FromJSON(symbolPrices)
		if err != nil {
			errorMessage := fmt.Sprintf("Couldn't process alert %s: Error parsing symbol prices for symbol %s: %v", alert.Id, alert.Symbol, err)
			logger.Error(errorMessage)
			return nil, errors.New(errorMessage + ". Cause: " + err.Error())
		}

		if _, exists := alertTypes[alert.Type]; exists && alertTypes[alert.Type](alert.Price, symbolValues) {
			if alert.Trigger == "once" {
				err = DisableAlert(server.MongoRepository, alert.Id)
			}
			if err != nil {
				errorMessage := fmt.Sprintf("Error disabling alert %s: %v", alert.Id, err)
				logger.Error(errorMessage)
				continue
			}
			triggeredAlerts = append(triggeredAlerts, alert)
		}
	}
	err := triggerAlerts(triggeredAlerts, *server.DiscordClient)
	if err != nil {
		errorMessage := fmt.Sprintf("Error triggering alerts: %v. Cause: %v", triggeredAlerts, err.Error())
		logger.Error(errorMessage)
		return nil, errors.New(errorMessage)
	}
	return triggeredAlerts, nil
}

func DisableAlert(mdb mongod.MongoRepository, alertId primitive.ObjectID) error {
	return alertsRepository.DisableAlertById(mdb, alertId)
}

func triggerAlerts(alert []alertsDTO.Alert, dsCli discord.Client) error {
	for _, alert := range alert {
		err := ds.NotifyAlert(dsCli, alert.AlertMessageFull())
		if err != nil {
			return errors.New("Error triggering alert: " + err.Error())
		}
	}
	return nil
}

func isReaching(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return alertPrice == intraDayPrices.Current
}

func isOver(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return intraDayPrices.Current > alertPrice
}

func isUnder(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return intraDayPrices.Current < alertPrice
}

func isCrossingUp(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return intraDayPrices.Previous < alertPrice && intraDayPrices.Current >= alertPrice
}

func isCrossingDown(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return intraDayPrices.Previous > alertPrice && intraDayPrices.Current <= alertPrice
}

func isCrossing(alertPrice float64, intraDayPrices dtos.IntraDayPrices) bool {
	return isCrossingUp(alertPrice, intraDayPrices) || isCrossingDown(alertPrice, intraDayPrices)
}
