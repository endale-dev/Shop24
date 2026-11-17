package mail

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"time"
)

type Log struct {
	Log *slog.Logger
}

func (l Log) Send(_ context.Context, to, subject, body string) error {
	l.Log.Info("mail", "to", to, "subject", subject, "body", body)
	return nil
}

type SMTP struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func (s SMTP) Send(ctx context.Context, to, subject, body string) error {
	addr := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
	var auth smtp.Auth
	if s.User != "" {
		auth = smtp.PlainAuth("", s.User, s.Password, s.Host)
	}
	msg := "From: " + s.From + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer client.Close()
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12}); err != nil {
			return err
		}
	}
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}
	if err := client.Mail(s.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(msg)); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}
