package fcm

import (
	"context"
	"errors"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/school-management/backend/internal/config"
)

var (
	ErrFCMNotInitialized = errors.New("klien FCM belum diinisialisasi")
	ErrInvalidToken      = errors.New("token FCM tidak valid")
	ErrSendFailed        = errors.New("gagal mengirim notifikasi")
)

// Client wraps the Firebase Cloud Messaging client
// Requirements: 13.1, 13.2 - Firebase Cloud Messaging integration for push notifications
type Client struct {
	app       *firebase.App
	messaging *messaging.Client
	projectID string
}

// NewClient creates a new FCM client
// Requirements: 13.1, 13.2 - THE System SHALL send push notification to parent's device via FCM
func NewClient(cfg config.FCMConfig) (*Client, error) {
	if cfg.CredentialsFile == "" {
		log.Println("FCM credentials file not configured, FCM client will be disabled")
		return &Client{}, nil
	}

	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.CredentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize FCM client: %w", err)
	}

	return &Client{
		app:       app,
		messaging: messagingClient,
		projectID: cfg.ProjectID,
	}, nil
}

// IsInitialized returns true if the FCM client is properly initialized
func (c *Client) IsInitialized() bool {
	return c.messaging != nil
}

// SendPushNotification sends a push notification to a single device
// Requirements: 13.1, 13.2 - THE System SHALL send push notification to parent's device via FCM
func (c *Client) SendPushNotification(ctx context.Context, deviceToken, title, body string, data map[string]string) error {
	if !c.IsInitialized() {
		log.Println("FCM client not initialized, skipping push notification")
		return nil
	}

	if deviceToken == "" {
		return ErrInvalidToken
	}

	message := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				Sound:       "default",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}

	_, err := c.messaging.Send(ctx, message)
	if err != nil {
		// Check if it's an invalid token error
		if messaging.IsUnregistered(err) || messaging.IsInvalidArgument(err) {
			return ErrInvalidToken
		}
		return fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	return nil
}


// SendMulticast sends a push notification to multiple devices
// Requirements: 13.1, 13.2 - THE System SHALL send push notification to parent's device via FCM
func (c *Client) SendMulticast(ctx context.Context, deviceTokens []string, title, body string, data map[string]string) (*MulticastResult, error) {
	if !c.IsInitialized() {
		log.Println("FCM client not initialized, skipping multicast notification")
		return &MulticastResult{}, nil
	}

	if len(deviceTokens) == 0 {
		return &MulticastResult{}, nil
	}

	// FCM allows max 500 tokens per multicast
	const maxTokens = 500
	result := &MulticastResult{
		SuccessCount: 0,
		FailureCount: 0,
		FailedTokens: make([]string, 0),
	}

	for i := 0; i < len(deviceTokens); i += maxTokens {
		end := i + maxTokens
		if end > len(deviceTokens) {
			end = len(deviceTokens)
		}

		batch := deviceTokens[i:end]
		batchResult, err := c.sendMulticastBatch(ctx, batch, title, body, data)
		if err != nil {
			return nil, err
		}

		result.SuccessCount += batchResult.SuccessCount
		result.FailureCount += batchResult.FailureCount
		result.FailedTokens = append(result.FailedTokens, batchResult.FailedTokens...)
	}

	return result, nil
}

// sendMulticastBatch sends a multicast message to a batch of tokens
func (c *Client) sendMulticastBatch(ctx context.Context, tokens []string, title, body string, data map[string]string) (*MulticastResult, error) {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				Sound:       "default",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := c.messaging.SendEachForMulticast(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	result := &MulticastResult{
		SuccessCount: response.SuccessCount,
		FailureCount: response.FailureCount,
		FailedTokens: make([]string, 0),
	}

	// Collect failed tokens
	for i, resp := range response.Responses {
		if !resp.Success {
			result.FailedTokens = append(result.FailedTokens, tokens[i])
		}
	}

	return result, nil
}

// SendToTopic sends a push notification to a topic
func (c *Client) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error {
	if !c.IsInitialized() {
		log.Println("FCM client not initialized, skipping topic notification")
		return nil
	}

	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
	}

	_, err := c.messaging.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	return nil
}

// SubscribeToTopic subscribes tokens to a topic
func (c *Client) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	if !c.IsInitialized() {
		return nil
	}

	_, err := c.messaging.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return nil
}

// UnsubscribeFromTopic unsubscribes tokens from a topic
func (c *Client) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	if !c.IsInitialized() {
		return nil
	}

	_, err := c.messaging.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}

	return nil
}

// MulticastResult represents the result of a multicast send operation
type MulticastResult struct {
	SuccessCount int
	FailureCount int
	FailedTokens []string
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
