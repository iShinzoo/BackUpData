package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackMessage struct {
	Text string `json:"text"`
}

func New(webhook string) *SlackNotifier {

	return &SlackNotifier{
		webhookURL: webhook,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *SlackNotifier) Notify(ctx context.Context, msg string) error {

	payload := slackMessage{
		Text: msg,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.webhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
