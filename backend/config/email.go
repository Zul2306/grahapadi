package config

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// LoadEmailConfig loads email configuration from environment
func LoadEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     getEnvStr("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvStr("SMTP_PORT", "587"),
		SMTPUsername: getEnvStr("SMTP_USERNAME", ""),
		SMTPPassword: getEnvStr("SMTP_PASSWORD", ""),
		FromEmail:    getEnvStr("FROM_EMAIL", "noreply@inventory.com"),
		FromName:     getEnvStr("FROM_NAME", "Inventory System"),
	}
}

func getEnvStr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// SendEmail sends an email using SMTP
func (ec *EmailConfig) SendEmail(to, subject, body string) error {
	if ec.SMTPUsername == "" || ec.SMTPPassword == "" {
		// For development, just log the email
		fmt.Printf("\n=== EMAIL (Development Mode) ===\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("Body:\n%s\n", body)
		fmt.Printf("================================\n\n")
		return nil
	}

	// Set up authentication
	auth := smtp.PlainAuth("", ec.SMTPUsername, ec.SMTPPassword, ec.SMTPHost)

	// Compose message
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", ec.FromName, ec.FromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", ec.SMTPHost, ec.SMTPPort)

	// Connect without TLS first
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Start TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         ec.SMTPHost,
	}

	if err = conn.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authenticate
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Set sender
	if err = conn.Mail(ec.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Set recipient
	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send message
	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}
