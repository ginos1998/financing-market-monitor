package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	mdb "github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	integrationsRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/integrations"
	dto "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos/integrations"
)

type Client struct {
	Client   *discordgo.Session
	Webhooks []dto.Webhook
}

func NewClient(token string, mdb mdb.MongoRepository) (*Client, error) {
	client, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	cli, err := integrationsRepository.FindIntegrationByAppName(mdb, "discord")
	if err != nil {
		return nil, errors.New("Error finding discord integration config: " + err.Error())
	}

	return &Client{
		Client:   client,
		Webhooks: cli.Webhooks,
	}, nil
}
