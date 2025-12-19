package service

import (
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"Gym_StrongCode/config"
	"Gym_StrongCode/internal/utils"

	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

type NotificationService struct {
	cfg       *config.Config
	emailPool chan *email.Email
	wg        sync.WaitGroup
	stop      chan struct{}
	logger    *zap.Logger
}

type Notification struct {
	To      string
	Subject string
	Body    string
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	ns := &NotificationService{
		cfg:    cfg,
		logger: utils.GetLogger(),
	}

	if cfg.SMTPHost == "" || cfg.FromEmail == "" {
		ns.logger.Warn("SMTP not configured - notifications will be logged only")
		return ns
	}

	ns.emailPool = make(chan *email.Email, 100)
	ns.stop = make(chan struct{})

	return ns
}

func (ns *NotificationService) SendNotification(to, subject, body string) error {

	if to == "" || subject == "" {
		ns.logger.Warn("Invalid notification parameters",
			zap.String("to", to),
			zap.String("subject", subject),
		)
		return fmt.Errorf("to and subject are required")
	}

	if ns.cfg.SMTPHost == "" {
		ns.logger.Info("Notification logged (no SMTP)",
			zap.String("to", to),
			zap.String("subject", subject),
		)
		return nil
	}

	e := email.NewEmail()
	e.From = ns.cfg.FromEmail
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	select {
	case ns.emailPool <- e:
		ns.logger.Debug("Email queued", zap.String("to", to))
	case <-time.After(5 * time.Second):
		ns.logger.Warn("Email queue full, dropping notification", zap.String("to", to))
		return fmt.Errorf("email queue full")
	}

	return nil
}

func (ns *NotificationService) StartWorker() {
	if ns.cfg.SMTPHost == "" {
		return
	}

	ns.wg.Add(1)
	go func() {
		defer ns.wg.Done()

		auth := smtp.PlainAuth("", ns.cfg.SMTPUser, ns.cfg.SMTPPass, ns.cfg.SMTPHost)
		addr := fmt.Sprintf("%s:%s", ns.cfg.SMTPHost, ns.cfg.SMTPPort)

		conn, err := smtp.Dial(addr)
		if err != nil {
			ns.logger.Error("Failed to connect to SMTP server",
				zap.String("host", ns.cfg.SMTPHost),
				zap.String("port", ns.cfg.SMTPPort),
				zap.Error(err),
			)
			return
		}
		conn.Close()

		for {
			select {
			case <-ns.stop:
				ns.logger.Info("Email worker stopping")
				return

			case e := <-ns.emailPool:
				if e == nil {
					continue
				}

				if err := ns.sendEmailWithRetry(e, addr, auth); err != nil {
					ns.logger.Error("Failed to send email after retries",
						zap.String("to", e.To[0]),
						zap.String("subject", e.Subject),
						zap.Error(err),
					)
				} else {
					ns.logger.Info("Email sent successfully",
						zap.String("to", e.To[0]),
						zap.String("subject", e.Subject),
					)
				}
			}
		}
	}()
}

func (ns *NotificationService) sendEmailWithRetry(e *email.Email, addr string, auth smtp.Auth) error {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := e.Send(addr, auth); err == nil {
			return nil
		} else {
			lastErr = err
			ns.logger.Warn("Email send attempt failed",
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", maxRetries),
				zap.String("to", e.To[0]),
				zap.Error(err),
			)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}
		}
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

func (ns *NotificationService) StopWorker() {
	if ns.cfg.SMTPHost == "" {
		return
	}

	select {
	case <-ns.stop:

		return
	default:
		close(ns.stop)
	}

	done := make(chan struct{})
	go func() {
		ns.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		ns.logger.Info("Email worker stopped successfully")
	case <-time.After(10 * time.Second):
		ns.logger.Warn("Timeout waiting for email worker to stop")
	}
}
