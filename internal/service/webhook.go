package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

// WebhookSender отправляет вебхуки
type WebhookSender struct {
	url    string
	client *http.Client
}

func NewWebhookSender(url string, timeout time.Duration) *WebhookSender {
	return &WebhookSender{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *WebhookSender) Send(ctx context.Context, payload domain.WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook responded with status %d", resp.StatusCode)
	}

	return nil
}

// WebhookWorker воркер для отправки вебхуков
type WebhookWorker struct {
	queue         repository.WebhookQueue
	sender        *WebhookSender
	retryAttempts int
	retryDelay    time.Duration
	popTimeout    time.Duration
}

func NewWebhookWorker(
	queue repository.WebhookQueue,
	sender *WebhookSender,
	retryAttempts int,
	retryDelay time.Duration,
) *WebhookWorker {
	return &WebhookWorker{
		queue:         queue,
		sender:        sender,
		retryAttempts: retryAttempts,
		retryDelay:    retryDelay,
		popTimeout:    2 * time.Second,
	}
}

func (w *WebhookWorker) Start(ctx context.Context) {
	log.Println("Webhook worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Webhook worker stopped")
			return
		default:
		}

		job, ok, err := w.queue.Dequeue(ctx, w.popTimeout)
		if err != nil {
			log.Printf("Webhook queue error: %v\n", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if !ok {
			continue
		}

		if err := w.sender.Send(ctx, job.Payload); err != nil {
			attempt := job.Attempt + 1
			if attempt <= w.retryAttempts {
				job.Attempt = attempt
				delay := w.retryDelay * time.Duration(attempt)
				log.Printf("Webhook failed (attempt %d/%d). Retrying in %s: %v\n", attempt, w.retryAttempts, delay, err)
				go func(job domain.WebhookJob, delay time.Duration) {
					time.Sleep(delay)
					if enqueueErr := w.queue.Enqueue(context.Background(), job); enqueueErr != nil {
						log.Printf("Failed to requeue webhook job: %v\n", enqueueErr)
					}
				}(*job, delay)
			} else {
				log.Printf("Webhook permanently failed after %d attempts: %v\n", w.retryAttempts, err)
			}
		}
	}
}
