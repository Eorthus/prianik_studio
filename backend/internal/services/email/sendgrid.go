package email

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridService struct {
	client *sendgrid.Client
	from   *mail.Email
}

func NewSendGridService(apiKey, fromEmail, fromName string) *SendGridService {
	return &SendGridService{
		client: sendgrid.NewSendClient(apiKey),
		from:   mail.NewEmail(fromName, fromEmail),
	}
}

func (s *SendGridService) SendEmail(to, subject, htmlContent, textContent string) error {
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(s.from, subject, toEmail, textContent, htmlContent)

	// Добавляем категорию для аналитики
	message.AddCategories("prianik-studio")

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("ошибка отправки через SendGrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid вернул ошибку: %d - %s", response.StatusCode, response.Body)
	}

	return nil
}
