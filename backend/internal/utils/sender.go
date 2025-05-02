package utils

import (
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"

	"pryanik_studio/internal/config"
	"pryanik_studio/internal/models"
)

// Sender интерфейс для отправки электронных писем
type Sender interface {
	SendOrderConfirmation(order *models.Order) error
	SendContactForm(form *models.ContactFormRequest) error
}

// GomailSender реализация Sender с использованием gomail
type GomailSender struct {
	config config.EmailConfig
	logger *logrus.Logger
}

// NewGomailSender создает новый экземпляр GomailSender
func NewGomailSender(config config.EmailConfig, logger *logrus.Logger) *GomailSender {
	return &GomailSender{
		config: config,
		logger: logger,
	}
}

// EmailTemplate содержит шаблоны для разных языков
type EmailTemplate struct {
	Subject map[string]string
	Body    map[string]string
}

// GetSubject возвращает тему письма на нужном языке
func (t *EmailTemplate) GetSubject(lang string) string {
	if subject, ok := t.Subject[lang]; ok {
		return subject
	}
	// Возвращаем тему на русском языке по умолчанию
	return t.Subject["ru"]
}

// GetBody возвращает тело письма на нужном языке
func (t *EmailTemplate) GetBody(lang string) string {
	if body, ok := t.Body[lang]; ok {
		return body
	}
	// Возвращаем тело на русском языке по умолчанию
	return t.Body["ru"]
}

// FormatCurrency форматирует валюту в зависимости от языка
func FormatCurrency(price float64, currency string, lang string) string {
	// Если валюта не указана, определяем её по языку
	if currency == "" {
		switch lang {
		case "en":
			currency = "USD"
		case "es":
			currency = "EUR"
		default:
			currency = "RUB"
		}
	}

	// Форматируем отображение в зависимости от валюты
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", price)
	case "EUR":
		return fmt.Sprintf("€%.2f", price)
	case "RUB":
		return fmt.Sprintf("%.2f ₽", price)
	default:
		return fmt.Sprintf("%.2f %s", price, currency)
	}
}

// SendOrderConfirmation отправляет уведомление о заказе клиенту и компании
func (s *GomailSender) SendOrderConfirmation(order *models.Order) error {
	// Определяем язык клиента (по умолчанию русский)
	lang := order.Language
	if lang == "" {
		lang = "ru"
	}

	// Создаем шаблоны для клиента
	customerTemplate := createOrderCustomerTemplate(order, lang)

	// Создаем сообщение для клиента
	customerMsg := gomail.NewMessage()
	customerMsg.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.MailFromName, s.config.MailFrom))
	customerMsg.SetHeader("To", order.Email)
	customerMsg.SetHeader("Subject", customerTemplate.GetSubject(lang))
	customerMsg.SetBody("text/html", customerTemplate.GetBody(lang))

	// Создаем шаблон для владельца на том же языке, что и клиент
	ownerTemplate := createOrderOwnerTemplate(order, lang)

	// Создаем сообщение для владельца
	ownerMsg := gomail.NewMessage()
	ownerMsg.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.MailFromName, s.config.MailFrom))
	ownerMsg.SetHeader("To", s.config.CompanyEmail)
	ownerMsg.SetHeader("Subject", ownerTemplate.GetSubject(lang))
	ownerMsg.SetBody("text/html", ownerTemplate.GetBody(lang))

	// Объединяем сообщения и отправляем их
	allMessages := []*gomail.Message{customerMsg, ownerMsg}
	return s.sendEmails(allMessages)
}

// SendContactForm отправляет уведомление о новом сообщении с формы обратной связи
func (s *GomailSender) SendContactForm(form *models.ContactFormRequest) error {
	// Определяем язык клиента (по умолчанию русский)
	lang := form.Language
	if lang == "" {
		lang = "ru"
	}

	// Создаем шаблоны для клиента
	customerTemplate := createContactCustomerTemplate(form, lang)

	// Создаем сообщение для клиента
	customerMsg := gomail.NewMessage()
	customerMsg.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.MailFromName, s.config.MailFrom))
	customerMsg.SetHeader("To", form.Email)
	customerMsg.SetHeader("Subject", customerTemplate.GetSubject(lang))
	customerMsg.SetBody("text/html", customerTemplate.GetBody(lang))

	// Создаем шаблон для владельца на том же языке, что и клиент
	ownerTemplate := createContactOwnerTemplate(form, lang)

	// Создаем сообщение для владельца
	ownerMsg := gomail.NewMessage()
	ownerMsg.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.MailFromName, s.config.MailFrom))
	ownerMsg.SetHeader("To", s.config.CompanyEmail)
	ownerMsg.SetHeader("Subject", ownerTemplate.GetSubject(lang))
	ownerMsg.SetBody("text/html", ownerTemplate.GetBody(lang))

	// Объединяем сообщения и отправляем их
	allMessages := []*gomail.Message{customerMsg, ownerMsg}
	return s.sendEmails(allMessages)
}

// createOrderCustomerTemplate создает шаблоны письма для клиента о заказе
func createOrderCustomerTemplate(order *models.Order, lang string) *EmailTemplate {
	template := &EmailTemplate{
		Subject: map[string]string{
			"ru": "Подтверждение заказа №" + fmt.Sprintf("%d", order.ID),
			"en": "Order Confirmation #" + fmt.Sprintf("%d", order.ID),
			"es": "Confirmación de pedido #" + fmt.Sprintf("%d", order.ID),
		},
		Body: make(map[string]string),
	}

	// Определяем валюту и форматируем общую стоимость
	currency := GetCurrencyByLanguage(lang)
	formattedTotalCost := FormatCurrency(order.TotalCost, currency, lang)

	// Формируем тело письма на русском
	ruBody := fmt.Sprintf(`
		<h2>Спасибо за ваш заказ №%d!</h2>
		<p>Уважаемый(ая) %s,</p>
		<p>Мы получили ваш заказ и обрабатываем его. Мы свяжемся с вами в ближайшее время для уточнения деталей.</p>
		<h3>Информация о заказе:</h3>
		<p><strong>Сумма заказа:</strong> %s</p>
		<p><strong>Дата заказа:</strong> %s</p>
		`,
		order.ID,
		order.Name,
		formattedTotalCost,
		order.CreatedAt.Format("02.01.2006 15:04"),
	)

	// Формируем тело письма на английском
	enBody := fmt.Sprintf(`
		<h2>Thank you for your order #%d!</h2>
		<p>Dear %s,</p>
		<p>We have received your order and are processing it. We will contact you shortly to confirm the details.</p>
		<h3>Order information:</h3>
		<p><strong>Order amount:</strong> %s</p>
		<p><strong>Order date:</strong> %s</p>
		`,
		order.ID,
		order.Name,
		formattedTotalCost,
		order.CreatedAt.Format("01/02/2006 15:04"),
	)

	// Формируем тело письма на испанском
	esBody := fmt.Sprintf(`
		<h2>¡Gracias por su pedido #%d!</h2>
		<p>Estimado/a %s,</p>
		<p>Hemos recibido su pedido y lo estamos procesando. Nos pondremos en contacto con usted en breve para confirmar los detalles.</p>
		<h3>Información del pedido:</h3>
		<p><strong>Importe del pedido:</strong> %s</p>
		<p><strong>Fecha del pedido:</strong> %s</p>
		`,
		order.ID,
		order.Name,
		formattedTotalCost,
		order.CreatedAt.Format("02/01/2006 15:04"),
	)

	// Добавляем товары, если они есть
	if len(order.Items) > 0 {
		// Русский
		ruItemsSection := "<h3>Товары:</h3><ul>"
		for _, item := range order.Items {
			formattedPrice := FormatCurrency(item.Price, currency, "ru")
			formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, "ru")
			ruItemsSection += fmt.Sprintf("<li>%s - %d шт. x %s = %s</li>",
				item.ProductName,
				item.Quantity,
				formattedPrice,
				formattedTotal,
			)
		}
		ruItemsSection += "</ul>"
		ruBody += ruItemsSection

		// Английский
		enItemsSection := "<h3>Products:</h3><ul>"
		for _, item := range order.Items {
			formattedPrice := FormatCurrency(item.Price, currency, "en")
			formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, "en")
			enItemsSection += fmt.Sprintf("<li>%s - %d pcs x %s = %s</li>",
				item.ProductName,
				item.Quantity,
				formattedPrice,
				formattedTotal,
			)
		}
		enItemsSection += "</ul>"
		enBody += enItemsSection

		// Испанский
		esItemsSection := "<h3>Productos:</h3><ul>"
		for _, item := range order.Items {
			formattedPrice := FormatCurrency(item.Price, currency, "es")
			formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, "es")
			esItemsSection += fmt.Sprintf("<li>%s - %d uds x %s = %s</li>",
				item.ProductName,
				item.Quantity,
				formattedPrice,
				formattedTotal,
			)
		}
		esItemsSection += "</ul>"
		esBody += esItemsSection
	}

	// Завершаем тело письма
	ruBody += `
		<p>Если у вас возникнут вопросы, пожалуйста, свяжитесь с нами по электронной почте или телефону.</p>
		<p>С уважением,<br>Команда Creality Workshop</p>
	`

	enBody += `
		<p>If you have any questions, please contact us by email or phone.</p>
		<p>Best regards,<br>Creality Workshop Team</p>
	`

	esBody += `
		<p>Si tiene alguna pregunta, por favor contáctenos por correo electrónico o teléfono.</p>
		<p>Atentamente,<br>Equipo de Creality Workshop</p>
	`

	template.Body["ru"] = ruBody
	template.Body["en"] = enBody
	template.Body["es"] = esBody

	return template
}

// GetCurrencyByLanguage возвращает валюту по умолчанию для языка
func GetCurrencyByLanguage(lang string) string {
	switch lang {
	case "en":
		return "USD"
	case "es":
		return "EUR"
	default:
		return "RUB"
	}
}

// createOrderOwnerTemplate создает шаблон письма для владельца о заказе
func createOrderOwnerTemplate(order *models.Order, lang string) *EmailTemplate {
	template := &EmailTemplate{
		Subject: map[string]string{
			"ru": "Новый заказ №" + fmt.Sprintf("%d", order.ID),
			"en": "New Order #" + fmt.Sprintf("%d", order.ID),
			"es": "Nuevo Pedido #" + fmt.Sprintf("%d", order.ID),
		},
		Body: make(map[string]string),
	}

	// Определяем валюту по языку
	currency := GetCurrencyByLanguage(lang)

	// Форматируем общую стоимость в валюте
	formattedTotalCost := FormatCurrency(order.TotalCost, currency, lang)

	// Форматируем тело письма в зависимости от языка
	switch lang {
	case "ru":
		template.Body["ru"] = fmt.Sprintf(`
			<h2>Новый заказ №%d</h2>
			<h3>Информация о клиенте:</h3>
			<p><strong>Имя:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Телефон:</strong> %s</p>
			<p><strong>Комментарий:</strong> %s</p>
			<h3>Информация о заказе:</h3>
			<p><strong>Сумма заказа:</strong> %s</p>
			<p><strong>Дата заказа:</strong> %s</p>
			`,
			order.ID,
			order.Name,
			order.Email,
			order.Phone,
			order.Comment,
			formattedTotalCost,
			order.CreatedAt.Format("02.01.2006 15:04"),
		)
	case "en":
		template.Body["en"] = fmt.Sprintf(`
			<h2>New Order #%d</h2>
			<h3>Customer Information:</h3>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Phone:</strong> %s</p>
			<p><strong>Comment:</strong> %s</p>
			<h3>Order Information:</h3>
			<p><strong>Order Amount:</strong> %s</p>
			<p><strong>Order Date:</strong> %s</p>
			`,
			order.ID,
			order.Name,
			order.Email,
			order.Phone,
			order.Comment,
			formattedTotalCost,
			order.CreatedAt.Format("01/02/2006 15:04"),
		)
	case "es":
		template.Body["es"] = fmt.Sprintf(`
			<h2>Nuevo Pedido #%d</h2>
			<h3>Información del Cliente:</h3>
			<p><strong>Nombre:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Teléfono:</strong> %s</p>
			<p><strong>Comentario:</strong> %s</p>
			<h3>Información del Pedido:</h3>
			<p><strong>Importe del Pedido:</strong> %s</p>
			<p><strong>Fecha del Pedido:</strong> %s</p>
			`,
			order.ID,
			order.Name,
			order.Email,
			order.Phone,
			order.Comment,
			formattedTotalCost,
			order.CreatedAt.Format("02/01/2006 15:04"),
		)
	}

	// Добавляем товары, если они есть
	if len(order.Items) > 0 {
		if _, ok := template.Body["ru"]; ok && lang == "ru" {
			ruItemsSection := "<h3>Товары:</h3><ul>"
			for _, item := range order.Items {
				formattedPrice := FormatCurrency(item.Price, currency, lang)
				formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, lang)

				ruItemsSection += fmt.Sprintf("<li>%s (ID: %d) - %d шт. x %s = %s</li>",
					item.ProductName,
					item.ProductID,
					item.Quantity,
					formattedPrice,
					formattedTotal,
				)
			}
			ruItemsSection += "</ul>"
			template.Body["ru"] += ruItemsSection
		}

		if _, ok := template.Body["en"]; ok && lang == "en" {
			enItemsSection := "<h3>Products:</h3><ul>"
			for _, item := range order.Items {
				formattedPrice := FormatCurrency(item.Price, currency, lang)
				formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, lang)

				enItemsSection += fmt.Sprintf("<li>%s (ID: %d) - %d pcs x %s = %s</li>",
					item.ProductName,
					item.ProductID,
					item.Quantity,
					formattedPrice,
					formattedTotal,
				)
			}
			enItemsSection += "</ul>"
			template.Body["en"] += enItemsSection
		}

		if _, ok := template.Body["es"]; ok && lang == "es" {
			esItemsSection := "<h3>Productos:</h3><ul>"
			for _, item := range order.Items {
				formattedPrice := FormatCurrency(item.Price, currency, lang)
				formattedTotal := FormatCurrency(float64(item.Quantity)*item.Price, currency, lang)

				esItemsSection += fmt.Sprintf("<li>%s (ID: %d) - %d uds x %s = %s</li>",
					item.ProductName,
					item.ProductID,
					item.Quantity,
					formattedPrice,
					formattedTotal,
				)
			}
			esItemsSection += "</ul>"
			template.Body["es"] += esItemsSection
		}
	}

	return template
}

// createContactCustomerTemplate создает шаблоны письма для клиента о форме обратной связи
func createContactCustomerTemplate(form *models.ContactFormRequest, lang string) *EmailTemplate {
	template := &EmailTemplate{
		Subject: map[string]string{
			"ru": "Мы получили ваше сообщение",
			"en": "We have received your message",
			"es": "Hemos recibido su mensaje",
		},
		Body: make(map[string]string),
	}

	// Формируем тело письма на русском
	ruBody := fmt.Sprintf(`
		<h2>Ваше сообщение получено</h2>
		<p>Уважаемый(ая) %s,</p>
		<p>Мы получили ваше сообщение и обрабатываем его. Мы свяжемся с вами в ближайшее время.</p>
		<p>С уважением,<br>Команда Creality Workshop</p>
	`, form.Name)

	// Формируем тело письма на английском
	enBody := fmt.Sprintf(`
		<h2>Your message has been received</h2>
		<p>Dear %s,</p>
		<p>We have received your message and are processing it. We will contact you shortly.</p>
		<p>Best regards,<br>Creality Workshop Team</p>
	`, form.Name)

	// Формируем тело письма на испанском
	esBody := fmt.Sprintf(`
		<h2>Su mensaje ha sido recibido</h2>
		<p>Estimado/a %s,</p>
		<p>Hemos recibido su mensaje y lo estamos procesando. Nos pondremos en contacto con usted en breve.</p>
		<p>Atentamente,<br>Equipo de Creality Workshop</p>
	`, form.Name)

	template.Body["ru"] = ruBody
	template.Body["en"] = enBody
	template.Body["es"] = esBody

	return template
}

// createContactOwnerTemplate создает шаблон письма для владельца о форме обратной связи
func createContactOwnerTemplate(form *models.ContactFormRequest, lang string) *EmailTemplate {
	template := &EmailTemplate{
		Subject: map[string]string{
			"ru": "Новое сообщение с формы обратной связи",
			"en": "New message from the contact form",
			"es": "Nuevo mensaje del formulario de contacto",
		},
		Body: make(map[string]string),
	}

	// Форматируем тело письма в зависимости от языка
	switch lang {
	case "ru":
		template.Body["ru"] = fmt.Sprintf(`
			<h2>Новое сообщение с формы обратной связи</h2>
			<h3>Информация о клиенте:</h3>
			<p><strong>Имя:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Телефон:</strong> %s</p>
			<h3>Сообщение:</h3>
			<p>%s</p>
		`,
			form.Name,
			form.Email,
			form.Phone,
			form.Message,
		)
	case "en":
		template.Body["en"] = fmt.Sprintf(`
			<h2>New message from the contact form</h2>
			<h3>Customer Information:</h3>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Phone:</strong> %s</p>
			<h3>Message:</h3>
			<p>%s</p>
		`,
			form.Name,
			form.Email,
			form.Phone,
			form.Message,
		)
	case "es":
		template.Body["es"] = fmt.Sprintf(`
			<h2>Nuevo mensaje del formulario de contacto</h2>
			<h3>Información del Cliente:</h3>
			<p><strong>Nombre:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Teléfono:</strong> %s</p>
			<h3>Mensaje:</h3>
			<p>%s</p>
		`,
			form.Name,
			form.Email,
			form.Phone,
			form.Message,
		)
	}

	return template
}

// sendEmails отправляет одно или несколько писем
func (s *GomailSender) sendEmails(messages []*gomail.Message) error {
	// Проверяем, что настройки SMTP заданы
	if s.config.SMTPHost == "" || s.config.SMTPUsername == "" || s.config.SMTPPassword == "" {
		s.logger.Warn("Настройки SMTP не заданы, письма не будут отправлены")
		return nil
	}

	// Настраиваем диалер для подключения к SMTP-серверу
	dialer := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUsername, s.config.SMTPPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Открываем соединение
	sender, err := dialer.Dial()
	if err != nil {
		s.logger.WithError(err).Error("Ошибка подключения к SMTP-серверу")
		return fmt.Errorf("ошибка подключения к SMTP-серверу: %w", err)
	}
	defer sender.Close()

	// Отправляем каждое сообщение
	for _, msg := range messages {
		if err := gomail.Send(sender, msg); err != nil {
			s.logger.WithError(err).Error("Ошибка отправки письма")
			return fmt.Errorf("ошибка отправки письма: %w", err)
		}
		s.logger.Info("Письмо успешно отправлено")
	}

	return nil
}
