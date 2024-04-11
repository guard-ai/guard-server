package notifications

import (
	"fmt"

	"github.com/guard-ai/guard-server/pkg"
	"github.com/guard-ai/guard-server/pkg/models"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

type Notifier struct {
	client *expo.PushClient
}

func NewNotifier() (*Notifier, error) {
	config := &expo.ClientConfig{AccessToken: pkg.Env().ExpoAccessToken}
	client := expo.NewPushClient(config)
	return &Notifier{client}, nil
}

func (n *Notifier) Broadcast(event models.Event, users []string) error {
	to := []expo.ExponentPushToken{}
	for _, user := range users {
		pushToken, err := expo.NewExponentPushToken(user)
		if err != nil {
			continue
		}

		to = append(to, pushToken)
	}

	response, err := n.client.Publish(&expo.PushMessage{
		To:         to,
		Title:      fmt.Sprintf("%s: %s", event.Level, event.Category),
		Body:       "Tap for more information!",
		Data:       map[string]string{},
		Sound:      "default",
		TTLSeconds: 15,
		Priority:   expo.HighPriority,
		Badge:      0,
		ChannelID:  pkg.Env().ExpoChannelId,
	})

	if err != nil {
		return err
	}

	return response.ValidateResponse()
}
