package stream

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const maxRetries = 4

var retryDelays = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	4 * time.Second,
	8 * time.Second,
}

// RetryHandler handles retry logic with exponential backoff
type RetryHandler struct {
	client        *redis.Client
	deadLetterKey string
}

// NewRetryHandler creates a new retry handler
func NewRetryHandler(client *redis.Client, deadLetterKey string) *RetryHandler {
	return &RetryHandler{
		client:        client,
		deadLetterKey: deadLetterKey,
	}
}

// Retries with exponential backoff
func (r *RetryHandler) RetryWithBackoff(ctx context.Context, fn func() error, streamID string, fields map[string]interface{}) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := fn(); err == nil {
			return nil // Success
		} else {
			lastErr = err
		}

		if attempt < maxRetries-1 {
			delay := retryDelays[attempt]
			log.Warn().
				Err(lastErr).
				Int("attempt", attempt+1).
				Dur("delay", delay).
				Str("stream_id", streamID).
				Msg("Retrying after delay")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	// All retries failed, send to death queue
	log.Error().
		Err(lastErr).
		Str("stream_id", streamID).
		Msg("All retries failed, sending to death queue")

	return r.sendToDeathQueue(ctx, streamID, fields, lastErr)
}

func (r *RetryHandler) sendToDeathQueue(ctx context.Context, streamID string, fields map[string]interface{}, err error) error {
	fields["_error"] = err.Error()
	fields["_stream_id"] = streamID
	fields["_failed_at"] = time.Now().Unix()

	args := make([]interface{}, 0, len(fields)*2+1)
	args = append(args, "*")
	for k, v := range fields {
		args = append(args, k, v)
	}

	_, err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.deadLetterKey,
		Values: fields,
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to send to death queue: %w", err)
	}

	log.Info().
		Str("stream_id", streamID).
		Str("dead_letter_key", r.deadLetterKey).
		Msg("Message sent to death queue")

	return nil
}
