package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	ds "github.com/ginos1998/financing-market-monitor/data-processing/config/discord"

	"github.com/sirupsen/logrus"
)

var logger logrus.Logger

func NotifyAlert(dsCli ds.Client, message string) error {
	token := dsCli.Webhooks[0].Token
	webhookId := dsCli.Webhooks[0].Id

	_, err := dsCli.Client.WebhookExecute(
		webhookId,
		token,
		true,
		&discordgo.WebhookParams{
			Content: message,
		},
	)
	if err != nil {
		return errors.New("Error sending alert to Discord: " + err.Error())
	}
	return nil
}
