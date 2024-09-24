package integrations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Webhook struct {
	ApplicationId string `bson:"application_id"`
	Avatar        string `bson:"avatar"`
	ChannelId     string `bson:"channel_id"`
	GuildId       string `bson:"guild_id"`
	Id            string `bson:"id"`
	Name          string `bson:"name"`
	Token         string `bson:"token"`
	Type          int    `bson:"type"`
	URL           string `bson:"url"`
}

type Integration struct {
	Id       primitive.ObjectID `bson:"_id"`
	App      string             `bson:"app"`
	Webhooks []Webhook          `bson:"webhooks"`
}
