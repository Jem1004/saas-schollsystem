package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/school-management/backend/internal/shared/fcm"
	"github.com/school-management/backend/internal/shared/redis"
)

// RetryConfig holds configuration for retry logic
// Requirements: 17.5 - IF FCM delivery fails, THEN THE System SHALL retry with exponential backoff
type RetryConfig struct {
	MaxRetries    int           // Maximum number of retry attempts
	InitialDelay  time.Duration // Initial delay before first retry
	MaxDelay      time.Duration // Maximum delay between retries
	BackoffFactor float64       // Multiplier for exponential backoff
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    5,
		InitialDelay:  1 * time.Second,
		MaxDelay:      5 * time.Minute,
		BackoffFactor: 2.0,
	}
}

// Worker processes notification queue and sends push notifications
// Requirements: 17.1, 17.2, 17.5 - Background queue processing with retry
type Worker struct {
	redisClient *redis.Client
	fcmClient   *fcm.Client
	repo        Repository
	retryConfig RetryConfig
	stopCh      chan struct{}
	wg          sync.WaitGroup
	running     bool
	mu          sync.Mutex
}

// NewWorker creates a new notification worker
func NewWorker(redisClient *redis.Client, fcmClient *fcm.Client, repo Repository) *Worker {
	return &Worker{
		redisClient: redisClient,
		fcmClient:   fcmClient,
		repo:        repo,
		retryConfig: DefaultRetryConfig(),
		stopCh:      make(chan struct{}),
	}
}

// NewWorkerWithConfig creates a new notification worker with custom retry config
func NewWorkerWithConfig(redisClient *redis.Client, fcmClient *fcm.Client, repo Repository, config RetryConfig) *Worker {
	return &Worker{
		redisClient: redisClient,
		fcmClient:   fcmClient,
		repo:        repo,
		retryConfig: config,
		stopCh:      make(chan struct{}),
	}
}

// Start starts the notification worker
// Requirements: 17.2 - WHEN a worker processes the queue, THE System SHALL send notification via FCM
func (w *Worker) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	w.wg.Add(1)
	go w.processLoop()

	log.Println("Notification worker started")
}

// Stop stops the notification worker gracefully
func (w *Worker) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.mu.Unlock()

	close(w.stopCh)
	w.wg.Wait()

	log.Println("Notification worker stopped")
}

// processLoop continuously processes the notification queue
func (w *Worker) processLoop() {
	defer w.wg.Done()

	for {
		select {
		case <-w.stopCh:
			return
		default:
			w.processQueue()
		}
	}
}

// processQueue processes items from the notification queue
// Requirements: 17.1 - THE System SHALL queue the notification in Redis
func (w *Worker) processQueue() {
	ctx := context.Background()

	// Try to dequeue a notification (blocking with timeout)
	data, err := w.redisClient.Dequeue(ctx, redis.NotificationQueue, 5*time.Second)
	if err != nil {
		log.Printf("Error dequeuing notification: %v", err)
		return
	}

	if data == "" {
		return // No item in queue
	}

	var item NotificationQueueItem
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		log.Printf("Error unmarshaling notification queue item: %v", err)
		return
	}

	// Process the notification
	if err := w.processNotification(ctx, &item); err != nil {
		log.Printf("Error processing notification %d: %v", item.NotificationID, err)
		w.handleRetry(ctx, &item, err)
	}
}


// processNotification processes a single notification
// Requirements: 17.2 - THE System SHALL send notification via FCM
func (w *Worker) processNotification(ctx context.Context, item *NotificationQueueItem) error {
	// Get user's active FCM tokens
	tokens, err := w.repo.FindActiveFCMTokensByUserID(ctx, item.UserID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		log.Printf("No active FCM tokens for user %d, skipping push notification", item.UserID)
		return nil
	}

	// Prepare data for FCM
	data := make(map[string]string)
	data["notification_id"] = uintToString(item.NotificationID)
	data["type"] = string(item.Type)

	// Add custom data if present
	if item.Data != nil {
		for k, v := range item.Data {
			if str, ok := v.(string); ok {
				data[k] = str
			}
		}
	}

	// Send to all active tokens
	tokenStrings := make([]string, len(tokens))
	for i, t := range tokens {
		tokenStrings[i] = t.Token
	}

	// Use multicast for efficiency
	result, err := w.fcmClient.SendMulticast(ctx, tokenStrings, item.Title, item.Message, data)
	if err != nil {
		return err
	}

	// Deactivate failed tokens
	if len(result.FailedTokens) > 0 {
		for _, failedToken := range result.FailedTokens {
			if err := w.repo.DeactivateFCMToken(ctx, failedToken); err != nil {
				log.Printf("Error deactivating failed token: %v", err)
			}
		}
	}

	log.Printf("Notification %d sent: %d success, %d failed", item.NotificationID, result.SuccessCount, result.FailureCount)
	return nil
}

// handleRetry handles retry logic for failed notifications
// Requirements: 17.5 - IF FCM delivery fails, THEN THE System SHALL retry with exponential backoff
func (w *Worker) handleRetry(ctx context.Context, item *NotificationQueueItem, originalErr error) {
	item.RetryCount++

	if item.RetryCount > w.retryConfig.MaxRetries {
		log.Printf("Notification %d exceeded max retries (%d), giving up", item.NotificationID, w.retryConfig.MaxRetries)
		return
	}

	// Calculate delay with exponential backoff
	delay := w.calculateBackoff(item.RetryCount)

	log.Printf("Scheduling retry %d for notification %d in %v", item.RetryCount, item.NotificationID, delay)

	// Schedule retry after delay
	go func() {
		time.Sleep(delay)

		// Re-queue the notification
		if err := w.redisClient.Enqueue(ctx, redis.NotificationQueue, item); err != nil {
			log.Printf("Error re-queuing notification %d: %v", item.NotificationID, err)
		}
	}()
}

// calculateBackoff calculates the delay for exponential backoff
// Requirements: 17.5 - Retry with exponential backoff
func (w *Worker) calculateBackoff(retryCount int) time.Duration {
	// Calculate exponential delay: initialDelay * (backoffFactor ^ (retryCount - 1))
	delay := float64(w.retryConfig.InitialDelay) * math.Pow(w.retryConfig.BackoffFactor, float64(retryCount-1))

	// Cap at max delay
	if delay > float64(w.retryConfig.MaxDelay) {
		delay = float64(w.retryConfig.MaxDelay)
	}

	return time.Duration(delay)
}

// ProcessPendingNotifications processes all pending notifications in the queue
// This can be called manually for batch processing
func (w *Worker) ProcessPendingNotifications(ctx context.Context, maxItems int) (int, error) {
	processed := 0

	for i := 0; i < maxItems; i++ {
		data, err := w.redisClient.DequeueNonBlocking(ctx, redis.NotificationQueue)
		if err != nil {
			return processed, err
		}

		if data == "" {
			break // No more items
		}

		var item NotificationQueueItem
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			log.Printf("Error unmarshaling notification queue item: %v", err)
			continue
		}

		if err := w.processNotification(ctx, &item); err != nil {
			log.Printf("Error processing notification %d: %v", item.NotificationID, err)
			w.handleRetry(ctx, &item, err)
		}

		processed++
	}

	return processed, nil
}

// GetQueueLength returns the current length of the notification queue
func (w *Worker) GetQueueLength(ctx context.Context) (int64, error) {
	return w.redisClient.QueueLength(ctx, redis.NotificationQueue)
}

// IsRunning returns true if the worker is currently running
func (w *Worker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// Helper function to convert uint to string
func uintToString(n uint) string {
	return fmt.Sprintf("%d", n)
}
