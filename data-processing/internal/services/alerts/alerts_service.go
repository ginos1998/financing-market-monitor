package alerts

import (
	"errors"
	"fmt"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	alertService "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/alerts"
	alertsRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/alerts"
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

func ProcessAlerts(mdb mongod.MongoRepository, redisClient *redis.RedisClient, alerts []alertsDTO.Alert) ([]alertsDTO.Alert, error) {
	triggeredAlerts := make([]alertsDTO.Alert, 0)
	for _, alert := range alerts {
		symbolKey := fmt.Sprintf("INTRA-DAY_%s", alert.Symbol)
		symbolPrices, err := redisService.GetIntraDaySymbolPrices(redisClient, symbolKey)
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
			err = triggerAlert(alert)
			if err != nil {
				errorMessage := fmt.Sprintf("Error triggering alert %s: %v", alert.Id, err)
				logger.Error(errorMessage)
				continue
			}
			if alert.Trigger == "once" {
				err = DisableAlert(mdb, alert.Id)
			}
			if err != nil {
				errorMessage := fmt.Sprintf("Error disabling alert %s: %v", alert.Id, err)
				logger.Error(errorMessage)
				continue
			}
			triggeredAlerts = append(triggeredAlerts, alert)
		}
	}
	return triggeredAlerts, nil
}

func DisableAlert(mdb mongod.MongoRepository, alertId primitive.ObjectID) error {
	return alertsRepository.DisableAlertById(mdb, alertId)
}

func triggerAlert(alert alertsDTO.Alert) error {
	logger.Info("Alert triggered: ", alert.ToString())
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
