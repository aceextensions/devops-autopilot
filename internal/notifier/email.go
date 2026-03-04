// Package notifier wraps github.com/aceextensions/notifier/email for use
// inside the devops-autopilot application.
package notifier

import (
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/aceextensions/notifier/email"
)

// EmailNotifier wraps the aceextensions email client.
type EmailNotifier struct {
	client *email.Client
	cfg    *config.Config
}

// NewEmail creates an Email notifier. Returns nil client if email is disabled in config.
func NewEmail(cfg *config.Config) *EmailNotifier {
	return &EmailNotifier{
		client: email.New(email.Config{
			SMTPHost: cfg.Email.SMTPHost,
			SMTPPort: cfg.Email.SMTPPort,
			Username: cfg.Email.Username,
			Password: cfg.Email.Password,
			From:     cfg.Email.Username,
			To:       []string{cfg.Email.To},
		}),
		cfg: cfg,
	}
}

// Send dispatches an email. No-op if email is disabled in config.
func (e *EmailNotifier) Send(subject, body string) error {
	if !e.cfg.Email.Enabled {
		return nil
	}
	return e.client.Send(subject, body)
}
