package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirebaseApp(ctx context.Context, credentialsPath []byte) (*firebase.App, error) {

	opt := option.WithCredentialsJSON(credentialsPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase newApp process: %w", err)
	}

	return app, nil
}
