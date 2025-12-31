package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// OutboxEventStatus represents the status of an outbox event
type OutboxEventStatus string

const (
	OutboxEventStatusPending   OutboxEventStatus = "pending"
	OutboxEventStatusPublished OutboxEventStatus = "published"
	OutboxEventStatusFailed    OutboxEventStatus = "failed"
)

// IsValid checks if the outbox event status is valid
func (s OutboxEventStatus) IsValid() bool {
	switch s {
	case OutboxEventStatusPending, OutboxEventStatusPublished, OutboxEventStatusFailed:
		return true
	}
	return false
}

// OutboxEvent represents an event to be published (Event Outbox Pattern)
type OutboxEvent struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	AggregateID uint              `gorm:"index;not null" json:"aggregate_id"`
	EventType   string            `gorm:"type:varchar(100);not null" json:"event_type"`
	Payload     string            `gorm:"type:jsonb;not null" json:"payload"`
	Status      OutboxEventStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	RetryCount  int               `gorm:"default:0" json:"retry_count"`
	CreatedAt   time.Time         `json:"created_at"`
	PublishedAt *time.Time        `json:"published_at"`
}

// TableName specifies the table name for OutboxEvent
func (OutboxEvent) TableName() string {
	return "outbox_events"
}

// Validate validates the outbox event data
func (o *OutboxEvent) Validate() error {
	if o.AggregateID == 0 {
		return errors.New("aggregate_id is required")
	}
	if strings.TrimSpace(o.EventType) == "" {
		return errors.New("event_type is required")
	}
	if strings.TrimSpace(o.Payload) == "" {
		return errors.New("payload is required")
	}
	// Validate JSON payload
	var js json.RawMessage
	if err := json.Unmarshal([]byte(o.Payload), &js); err != nil {
		return errors.New("payload must be valid JSON")
	}
	return nil
}

// MarkAsPublished marks the event as published
func (o *OutboxEvent) MarkAsPublished() {
	o.Status = OutboxEventStatusPublished
	now := time.Now()
	o.PublishedAt = &now
}

// MarkAsFailed marks the event as failed
func (o *OutboxEvent) MarkAsFailed() {
	o.Status = OutboxEventStatusFailed
	o.RetryCount++
}

// CanRetry checks if the event can be retried
func (o *OutboxEvent) CanRetry(maxRetries int) bool {
	return o.Status == OutboxEventStatusFailed && o.RetryCount < maxRetries
}

// ResetForRetry resets the event status for retry
func (o *OutboxEvent) ResetForRetry() {
	o.Status = OutboxEventStatusPending
}

// SetPayload sets the payload from a map
func (o *OutboxEvent) SetPayload(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	o.Payload = string(jsonData)
	return nil
}

// GetPayload retrieves the payload as a map
func (o *OutboxEvent) GetPayload() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(o.Payload), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
