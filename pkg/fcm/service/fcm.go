package sender

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type PushSender interface {
	SendToToken(ctx context.Context, token, title, body string, data map[string]string) error
}

type firebaseService struct {
	client *messaging.Client
}

func NewFirebaseService(ctx context.Context, app *firebase.App) (PushSender, error) {
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseService{
		client: client,
	}, nil
}

func (s *firebaseService) SendToToken(ctx context.Context, token string, title string, body string, data map[string]string) error {

	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	_, err := s.client.Send(ctx, msg)
	return err
}
