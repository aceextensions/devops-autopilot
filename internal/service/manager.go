package service

import (
	"fmt"
	"os"
	"time"

	"github.com/aceextensions/devops-autopilot/internal/alert"
	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/aceextensions/devops-autopilot/internal/notifier"
	"github.com/aceextensions/devops-autopilot/internal/report"
	"github.com/kardianos/service"
)

type program struct {
	exit    chan struct{}
	cfgPath string
}

func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	cfg, err := config.Load(p.cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		return
	}

	ticker := time.NewTicker(time.Duration(cfg.Monitor.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	slackNs := notifier.NewSlack(cfg)
	emailNs := notifier.NewEmail(cfg)

	// run initial check
	p.check(cfg, slackNs, emailNs)

	for {
		select {
		case <-p.exit:
			return
		case <-ticker.C:
			p.check(cfg, slackNs, emailNs)
		}
	}
}

func (p *program) check(cfg *config.Config, slackNs *notifier.SlackNotifier, emailNs *notifier.EmailNotifier) {
	sys, err := collector.CollectSystem(cfg.ServerName)
	if err != nil {
		fmt.Printf("Error collecting system stats: %v\n", err)
		return
	}

	doc, _ := collector.CollectDocker()
	alerts := alert.CheckThresholds(cfg, sys)

	if len(alerts) > 0 {
		rpt := report.Build(sys, doc, alerts)
		_ = slackNs.Send(rpt)
		_ = emailNs.Send(fmt.Sprintf("Alerts on %s", cfg.ServerName), rpt)
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	return nil
}

func Manage(action string, cfgPath string) error {
	svcConfig := &service.Config{
		Name:        "devops-autopilot",
		DisplayName: "DevOps Autopilot",
		Description: "Monitoring and alerting daemon",
		Arguments:   []string{"--config", cfgPath, "service", "run"},
	}

	prg := &program{cfgPath: cfgPath}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}

	if action == "run" {
		return s.Run()
	}

	if action == "status" {
		status, err := s.Status()
		if err != nil {
			return err
		}
		switch status {
		case service.StatusRunning:
			fmt.Println("Service is running")
		case service.StatusStopped:
			fmt.Println("Service is stopped")
		default:
			fmt.Println("Service status unknown")
		}
		return nil
	}

	return service.Control(s, action)
}
