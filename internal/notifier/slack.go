// Package notifier wraps github.com/aceextensions/notifier/slack for use
// inside the devops-autopilot application.
package notifier

import (
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/aceextensions/notifier/slack"
)

// SlackNotifier wraps the aceextensions slack client.
type SlackNotifier struct {
	client *slack.Client
	cfg    *config.Config
}

// NewSlack creates a Slack notifier. Returns nil client if Slack is disabled in config.
func NewSlack(cfg *config.Config) *SlackNotifier {
	return &SlackNotifier{
		client: slack.New(cfg.Slack.WebhookURL),
		cfg:    cfg,
	}
}

// Send dispatches a Slack message. No-op if slack is disabled in config.
func (s *SlackNotifier) Send(message string) error {
	if !s.cfg.Slack.Enabled {
		return nil
	}
	return s.client.Send("DevOps Alert", message)
}
