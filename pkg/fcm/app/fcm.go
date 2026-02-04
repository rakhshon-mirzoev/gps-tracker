package fcm

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirebaseApp(ctx context.Context, credentialsPath string) (*firebase.App, error) {
	credentialsJSON, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("NewFirebaseApp read file: %w", err)
	}

	opt := option.WithCredentialsJSON(credentialsJSON)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase newApp process: %w", err)
	}

	return app, nil
}
