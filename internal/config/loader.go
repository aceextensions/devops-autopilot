package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerName string `yaml:"server_name"`
	Alerts     struct {
		CPU    int `yaml:"cpu"`
		Memory int `yaml:"memory"`
		Disk   int `yaml:"disk"`
	} `yaml:"alerts"`
	Slack struct {
		Enabled    bool   `yaml:"enabled"`
		WebhookURL string `yaml:"webhook_url"`
	} `yaml:"slack"`
	Email struct {
		Enabled  bool   `yaml:"enabled"`
		SMTPHost string `yaml:"smtp_host"`
		SMTPPort int    `yaml:"smtp_port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		To       string `yaml:"to"`
	} `yaml:"email"`
	Monitor struct {
		IntervalSeconds int `yaml:"interval_seconds"`
	} `yaml:"monitor"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if env := os.Getenv("SERVER_NAME"); env != "" {
		cfg.ServerName = env
	}
	if env := os.Getenv("SLACK_WEBHOOK_URL"); env != "" {
		cfg.Slack.WebhookURL = env
	}
	if env := os.Getenv("SMTP_PASSWORD"); env != "" {
		cfg.Email.Password = env
	}

	if cfg.Monitor.IntervalSeconds == 0 {
		cfg.Monitor.IntervalSeconds = 60
	}

	return &cfg, nil
}
