package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/smtp"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"strings"
	"time"
)

type IMailService interface {
	SendUsersEmails(e models.UsersEmails) error
	SendVerificationCode(to string, data models.VerifyEmail) error
	SendEmail(e models.Email) error
}

type MailService struct {
	UserService *UserService

	email    string
	password string
	host     string
	port     string
}

func NewMailService(userService *UserService) *MailService {
	return &MailService{
		UserService: userService,
		email:       config.App.MailUser,
		password:    config.App.MailPass,
		host:        config.App.MailHost,
		port:        config.App.MailPort,
	}
}

func (m *MailService) SendUsersEmails(e models.UsersEmails) error {
	if e.Preview {
		user, err := m.UserService.GetUserDetails(0)
		if err != nil {
			return err
		}
		slog.Info("Sending preview email")
		template, err := checkEmailType(m, e, user.Username)
		if err != nil {
			return err
		}
		m.SendEmail(models.Email{
			To:      []string{user.Email},
			Subject: e.Subject,
			HTML:    template,
		})
		return nil
	}
	users, err := m.UserService.GetAllByIndices(e.UserIDs)
	if err != nil {
		return err
	}
	if len(users) != len(e.UserIDs) || len(users) == 0 {
		return fmt.Errorf("Some or all users not found")
	}

	for _, user := range users {
		template, err := checkEmailType(m, e, user.Username)
		if err != nil {
			return err
		}
		m.SendEmail(models.Email{
			To:      []string{user.Email},
			Subject: e.Subject,
			HTML:    template,
		})
	}
	return nil
}

func (m *MailService) SendVerificationCode(to string, data models.VerifyEmail) error {
	tem, err := utils.GetTemplate(string(utils.EmailVerificationCodeEn), data)
	if err != nil {
		return err
	}
	return m.SendEmail(models.Email{
		To:      []string{to},
		Subject: "Verify your email",
		HTML:    tem,
	})
}

func (m *MailService) SendEmail(e models.Email) error {

	recipients := append([]string{}, e.To...)
	recipients = append(recipients, e.Cc...)
	recipients = append(recipients, e.Bcc...)

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	boundary := fmt.Sprintf("boundary-%d", time.Now().UnixNano())

	headers := []string{
		"From: " + m.email,
		"To: " + strings.Join(e.To, ","),
		"Subject: " + e.Subject,
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		"Message-ID: <" + fmt.Sprint(time.Now().UnixNano()) + "@sea.uofk.edu>",
		"Content-Type: multipart/alternative; boundary=\"" + boundary + "\"",
	}

	if len(e.Cc) > 0 {
		headers = append(headers, "Cc: "+strings.Join(e.Cc, ","))
	}
	if len(e.ReplyTo) > 0 {
		headers = append(headers, "Reply-To: "+e.ReplyTo)
	}

	body := ""

	if e.Text != "" {
		body += "--" + boundary + "\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: 7bit\r\n\r\n" +
			e.Text + "\r\n\r\n"
	}

	if e.HTML != "" {
		body += "--" + boundary + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: 7bit\r\n\r\n" +
			e.HTML + "\r\n\r\n"
	}

	if body == "" {
		return fmt.Errorf("no body provided")
	}
	message := []byte(
		strings.Join(headers, "\r\n") + "\r\n\r\n" +
			body +
			"--" + boundary + "--\r\n",
	)

	auth := smtp.PlainAuth("", m.email, m.password, m.host)

	err := smtp.SendMail(
		m.host+":"+m.port,
		auth,
		m.email,
		recipients,
		message,
	)
	return err
}

// ======== HELPERS ========

func checkEmailType(m *MailService, e models.UsersEmails, username string) (string, error) {
	var template string
	var err error
	switch e.Type {
	case models.TECHNICAL:
		var data models.TechnicalEmail
		if err := json.Unmarshal(e.Data, &data); err != nil {
			return "", err
		}
		data.Message = strings.ReplaceAll(data.Message, "\n", "<br>")
		template, err = utils.GetEmailTechnicalTemplate(models.TechnicalEmailTemplate{
			TechnicalEmail: data,
			Username:       username,
			Year:           time.Now().Year(),
		})
		if err != nil {
			return "", err
		}

	case models.CERTIFICATE:
		var data models.CertificateEmailData
		if err := json.Unmarshal(e.Data, &data); err != nil {
			return "", err
		}
		template, err = utils.GetEnglishCertificateTemplate(data)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("invalid email type")
	}

	return template, nil
}
