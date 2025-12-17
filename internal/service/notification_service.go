package service

import (
	"fmt"
	"net/smtp"
	"sync"

	"Gym-StrongCode/config"
	"Gym-StrongCode/internal/utils"

	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

type NotificationService struct {
	cfg       *config.Config
	emailPool chan *email.Email
	wg        sync.WaitGroup
	stop      chan struct{}
}

type Notification struct {
	To      string
	Subject string
	Body    string
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	if cfg.SMTPHost == "" || cfg.FromEmail == "" {
		utils.GetLogger().Warn("SMTP not configured - notifications will be logged only")
		return &NotificationService{}
	}

	ns := &NotificationService{
		cfg:       cfg,
		emailPool: make(chan *email.Email, 100),
		stop:      make(chan struct{}),
	}

	return ns
}

func (ns *NotificationService) SendNotification(to, subject, body string) {
	logger := utils.GetLogger()

	if ns.cfg == nil || ns.cfg.SMTPHost == "" {
		logger.Info("Email notification (SMTP not configured)",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.String("body", body),
		)
		return
	}

	e := email.NewEmail()
	e.From = ns.cfg.FromEmail
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	select {
	case ns.emailPool <- e:
	default:
		logger.Warn("Email queue full, dropping notification")
	}
}

func (ns *NotificationService) StartWorker() {
	if ns.cfg.SMTPHost == "" {
		return
	}

	ns.wg.Add(1)
	go func() {
		defer ns.wg.Done()
		auth := smtp.PlainAuth("", ns.cfg.SMTPUser, ns.cfg.SMTPPass, ns.cfg.SMTPHost)

		for {
			select {
			case <-ns.stop:
				return
			case e := <-ns.emailPool:
				addr := fmt.Sprintf("%s:%s", ns.cfg.SMTPHost, ns.cfg.SMTPPort)
				if err := e.Send(addr, auth); err != nil {
					utils.GetLogger().Error("Failed to send email", zap.Error(err))
				} else {
					utils.GetLogger().Info("Email sent", zap.String("to", e.To[0]), zap.String("subject", e.Subject))
				}
			}
		}
	}()
}

func (ns *NotificationService) StopWorker() {
	if ns.cfg.SMTPHost == "" {
		return
	}
	close(ns.stop)
	ns.wg.Wait()
}
