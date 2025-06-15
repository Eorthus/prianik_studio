package utils

import (
	"fmt"
	"pryanik_studio/internal/config"
	"pryanik_studio/internal/models"
	"pryanik_studio/internal/services/email"

	"github.com/sirupsen/logrus"
)

type SendGridSender struct {
	service *email.SendGridService
	config  config.EmailConfig
	logger  *logrus.Logger
}

func NewSendGridSender(config config.EmailConfig, logger *logrus.Logger) *SendGridSender {
	service := email.NewSendGridService(
		config.SendGridAPIKey,
		config.MailFrom,
		config.MailFromName,
	)

	return &SendGridSender{
		service: service,
		config:  config,
		logger:  logger,
	}
}

func (s *SendGridSender) SendOrderConfirmation(order *models.Order) error {
	// Определяем язык клиента (по умолчанию русский)
	lang := order.Language
	if lang == "" {
		lang = "ru"
	}

	// Письмо клиенту
	clientHTML := s.generateOrderClientHTML(order, lang)
	clientText := s.generateOrderClientText(order, lang)
	clientSubject := s.getOrderSubject(lang, order.ID)

	if err := s.service.SendEmail(order.Email, clientSubject, clientHTML, clientText); err != nil {
		s.logger.WithError(err).Error("Ошибка отправки письма клиенту")
		return err
	}

	// Письмо администратору (всегда на том же языке, что и клиент)
	adminHTML := s.generateOrderAdminHTML(order, lang)
	adminText := s.generateOrderAdminText(order, lang)
	adminSubject := s.getAdminOrderSubject(lang, order.ID)

	if err := s.service.SendEmail(s.config.CompanyEmail, adminSubject, adminHTML, adminText); err != nil {
		s.logger.WithError(err).Error("Ошибка отправки письма администратору")
		// Не возвращаем ошибку, так как клиенту письмо уже отправлено
	}

	s.logger.WithField("order_id", order.ID).Info("Уведомления о заказе отправлены")
	return nil
}

func (s *SendGridSender) SendContactForm(form *models.ContactFormRequest) error {
	// Определяем язык клиента (по умолчанию русский)
	lang := form.Language
	if lang == "" {
		lang = "ru"
	}

	// Подтверждение клиенту
	clientHTML := s.generateContactClientHTML(form, lang)
	clientText := s.generateContactClientText(form, lang)
	clientSubject := s.getContactSubject(lang)

	if err := s.service.SendEmail(form.Email, clientSubject, clientHTML, clientText); err != nil {
		s.logger.WithError(err).Error("Ошибка отправки подтверждения клиенту")
		return err
	}

	// Уведомление администратору (на том же языке, что и клиент)
	adminHTML := s.generateContactAdminHTML(form, lang)
	adminText := s.generateContactAdminText(form, lang)
	adminSubject := s.getAdminContactSubject(lang)

	if err := s.service.SendEmail(s.config.CompanyEmail, adminSubject, adminHTML, adminText); err != nil {
		s.logger.WithError(err).Error("Ошибка отправки уведомления администратору")
		// Не возвращаем ошибку, так как клиенту подтверждение уже отправлено
	}

	s.logger.WithField("email", form.Email).Info("Форма обратной связи обработана")
	return nil
}

// Функция для форматирования валюты в зависимости от языка
func (s *SendGridSender) formatCurrency(price float64, lang string) string {
	switch lang {
	case "en":
		return fmt.Sprintf("$%.2f", price)
	case "es":
		return fmt.Sprintf("€%.2f", price)
	default:
		return fmt.Sprintf("%.2f ₽", price)
	}
}

// Генерация HTML для заказа клиенту с локализацией
func (s *SendGridSender) generateOrderClientHTML(order *models.Order, lang string) string {
	var title, greeting, message, orderInfoTitle, orderNumber, totalAmount, orderDate, itemsTitle, questionsText, regards string

	switch lang {
	case "en":
		title = fmt.Sprintf("Order Confirmation #%d", order.ID)
		greeting = fmt.Sprintf("Dear %s,", order.Name)
		message = "We have received your order and are processing it. We will contact you shortly to confirm the details."
		orderInfoTitle = "Order Information:"
		orderNumber = "Order Number:"
		totalAmount = "Total Amount:"
		orderDate = "Order Date:"
		itemsTitle = "Ordered Items:"
		questionsText = "If you have any questions, please contact us:"
		regards = "Best regards,<br><strong>Prianik Studio Team</strong>"
	case "es":
		title = fmt.Sprintf("Confirmación de pedido #%d", order.ID)
		greeting = fmt.Sprintf("Estimado/a %s,", order.Name)
		message = "Hemos recibido su pedido y lo estamos procesando. Nos pondremos en contacto con usted en breve para confirmar los detalles."
		orderInfoTitle = "Información del Pedido:"
		orderNumber = "Número de Pedido:"
		totalAmount = "Importe Total:"
		orderDate = "Fecha del Pedido:"
		itemsTitle = "Productos Pedidos:"
		questionsText = "Si tiene alguna pregunta, por favor contáctenos:"
		regards = "Atentamente,<br><strong>Equipo de Prianik Studio</strong>"
	default: // ru
		title = fmt.Sprintf("Подтверждение заказа №%d", order.ID)
		greeting = fmt.Sprintf("Уважаемый(ая) %s,", order.Name)
		message = "Мы получили ваш заказ и обрабатываем его. Мы свяжемся с вами в ближайшее время для уточнения деталей."
		orderInfoTitle = "Информация о заказе:"
		orderNumber = "Номер заказа:"
		totalAmount = "Общая сумма:"
		orderDate = "Дата заказа:"
		itemsTitle = "Заказанные товары:"
		questionsText = "Если у вас возникнут вопросы, пожалуйста, свяжитесь с нами:"
		regards = "С уважением,<br><strong>Команда Prianik Studio</strong>"
	}

	itemsHTML := ""
	if len(order.Items) > 0 {
		itemsHTML = fmt.Sprintf("<h3>%s</h3><ul>", itemsTitle)
		for _, item := range order.Items {
			total := item.Price * float64(item.Quantity)
			formattedPrice := s.formatCurrency(item.Price, lang)
			formattedTotal := s.formatCurrency(total, lang)

			var itemText string
			switch lang {
			case "en":
				itemText = fmt.Sprintf("%s - %d pcs × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			case "es":
				itemText = fmt.Sprintf("%s - %d uds × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			default:
				itemText = fmt.Sprintf("%s - %d шт. × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			}

			itemsHTML += fmt.Sprintf("<li>%s</li>", itemText)
		}
		itemsHTML += "</ul>"
	}

	dateFormat := order.CreatedAt.Format("02.01.2006 15:04")
	if lang == "en" {
		dateFormat = order.CreatedAt.Format("01/02/2006 15:04")
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">%s</h2>
		
		<p>%s</p>
		
		<p>%s</p>
		
		<div style="background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
			<h3 style="margin-top: 0;">%s</h3>
			<p><strong>%s</strong> %d</p>
			<p><strong>%s</strong> %s</p>
			<p><strong>%s</strong> %s</p>
			%s
		</div>
		
		<p>%s</p>
		<ul>
			<li>Email: %s</li>
			<li>Telefone: +506 8415 4807</li>
		</ul>
		
		<p>%s</p>
	</div>
</body>
</html>
	`, title, title, greeting, message, orderInfoTitle, orderNumber, order.ID,
		totalAmount, s.formatCurrency(order.TotalCost, lang), orderDate, dateFormat,
		itemsHTML, questionsText, s.config.CompanyEmail, regards)
}

// Генерация текста для заказа клиенту с локализацией
func (s *SendGridSender) generateOrderClientText(order *models.Order, lang string) string {
	var title, greeting, message, orderNumber, totalAmount, orderDate, itemsTitle, questionsText, regards string

	switch lang {
	case "en":
		title = fmt.Sprintf("Order Confirmation #%d", order.ID)
		greeting = fmt.Sprintf("Dear %s,", order.Name)
		message = "We have received your order and are processing it. We will contact you shortly."
		orderNumber = "Order Number:"
		totalAmount = "Total Amount:"
		orderDate = "Order Date:"
		itemsTitle = "Ordered Items:"
		questionsText = "If you have any questions:"
		regards = "Best regards,\nPrianik Studio Team"
	case "es":
		title = fmt.Sprintf("Confirmación de pedido #%d", order.ID)
		greeting = fmt.Sprintf("Estimado/a %s,", order.Name)
		message = "Hemos recibido su pedido y lo estamos procesando. Nos pondremos en contacto con usted en breve."
		orderNumber = "Número de Pedido:"
		totalAmount = "Importe Total:"
		orderDate = "Fecha del Pedido:"
		itemsTitle = "Productos Pedidos:"
		questionsText = "Si tiene alguna pregunta:"
		regards = "Atentamente,\nEquipo de Prianik Studio"
	default: // ru
		title = fmt.Sprintf("Подтверждение заказа №%d", order.ID)
		greeting = fmt.Sprintf("Уважаемый(ая) %s,", order.Name)
		message = "Мы получили ваш заказ и обрабатываем его. Мы свяжемся с вами в ближайшее время."
		orderNumber = "Номер заказа:"
		totalAmount = "Общая сумма:"
		orderDate = "Дата заказа:"
		itemsTitle = "Заказанные товары:"
		questionsText = "Если у вас есть вопросы:"
		regards = "С уважением,\nКоманда Prianik Studio"
	}

	itemsText := ""
	if len(order.Items) > 0 {
		itemsText = fmt.Sprintf("\n%s\n", itemsTitle)
		for _, item := range order.Items {
			total := item.Price * float64(item.Quantity)
			formattedPrice := s.formatCurrency(item.Price, lang)
			formattedTotal := s.formatCurrency(total, lang)

			var itemText string
			switch lang {
			case "en":
				itemText = fmt.Sprintf("- %s - %d pcs × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			case "es":
				itemText = fmt.Sprintf("- %s - %d uds × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			default:
				itemText = fmt.Sprintf("- %s - %d шт. × %s = %s", item.ProductName, item.Quantity, formattedPrice, formattedTotal)
			}

			itemsText += itemText + "\n"
		}
	}

	dateFormat := order.CreatedAt.Format("02.01.2006 15:04")
	if lang == "en" {
		dateFormat = order.CreatedAt.Format("01/02/2006 15:04")
	}

	return fmt.Sprintf(`%s

%s

%s

Информация о заказе:
- %s %d
- %s %s
- %s %s
%s
%s
- Email: %s
- Телефон: +506 8415 4807

%s
	`, title, greeting, message, orderNumber, order.ID, totalAmount,
		s.formatCurrency(order.TotalCost, lang), orderDate, dateFormat,
		itemsText, questionsText, s.config.CompanyEmail, regards)
}

// Генерация HTML для администратора с локализацией
func (s *SendGridSender) generateOrderAdminHTML(order *models.Order, lang string) string {
	var title, customerInfo, orderInfo string

	switch lang {
	case "en":
		title = fmt.Sprintf("New Order #%d", order.ID)
		customerInfo = "Customer Information:"
		orderInfo = "Order Information:"
	case "es":
		title = fmt.Sprintf("Nuevo Pedido #%d", order.ID)
		customerInfo = "Información del Cliente:"
		orderInfo = "Información del Pedido:"
	default: // ru
		title = fmt.Sprintf("Новый заказ №%d", order.ID)
		customerInfo = "Информация о клиенте:"
		orderInfo = "Информация о заказе:"
	}

	itemsHTML := ""
	if len(order.Items) > 0 {
		var itemsTitle string
		switch lang {
		case "en":
			itemsTitle = "Products:"
		case "es":
			itemsTitle = "Productos:"
		default:
			itemsTitle = "Товары:"
		}

		itemsHTML = fmt.Sprintf("<h3>%s</h3><ul>", itemsTitle)
		for _, item := range order.Items {
			total := item.Price * float64(item.Quantity)
			formattedPrice := s.formatCurrency(item.Price, lang)
			formattedTotal := s.formatCurrency(total, lang)

			itemsHTML += fmt.Sprintf(`<li>%s (ID: %d) - %d × %s = %s</li>`,
				item.ProductName, item.ProductID, item.Quantity, formattedPrice, formattedTotal)
		}
		itemsHTML += "</ul>"
	}

	return fmt.Sprintf(`
<h2>%s</h2>
<h3>%s</h3>
<p><strong>Имя:</strong> %s</p>
<p><strong>Email:</strong> %s</p>
<p><strong>Телефон:</strong> %s</p>
<p><strong>Комментарий:</strong> %s</p>

<h3>%s</h3>
<p><strong>Сумма заказа:</strong> %s</p>
<p><strong>Дата заказа:</strong> %s</p>
%s
	`, title, customerInfo, order.Name, order.Email, order.Phone, order.Comment,
		orderInfo, s.formatCurrency(order.TotalCost, lang),
		order.CreatedAt.Format("02.01.2006 15:04"), itemsHTML)
}

func (s *SendGridSender) generateOrderAdminText(order *models.Order, lang string) string {
	return fmt.Sprintf("Новый заказ №%d от %s (%s) на сумму %s",
		order.ID, order.Name, order.Email, s.formatCurrency(order.TotalCost, lang))
}

// Контактная форма клиенту с локализацией
func (s *SendGridSender) generateContactClientHTML(form *models.ContactFormRequest, lang string) string {
	var title, greeting, message, regards string

	switch lang {
	case "en":
		title = "Your message has been received"
		greeting = fmt.Sprintf("Dear %s,", form.Name)
		message = "We have received your message and are processing it. We will contact you shortly."
		regards = "Best regards,<br><strong>Prianik Studio Team</strong>"
	case "es":
		title = "Su mensaje ha sido recibido"
		greeting = fmt.Sprintf("Estimado/a %s,", form.Name)
		message = "Hemos recibido su mensaje y lo estamos procesando. Nos pondremos en contacto con usted en breve."
		regards = "Atentamente,<br><strong>Equipo de Prianik Studio</strong>"
	default: // ru
		title = "Ваше сообщение получено"
		greeting = fmt.Sprintf("Уважаемый(ая) %s,", form.Name)
		message = "Мы получили ваше сообщение и обрабатываем его. Мы свяжемся с вами в ближайшее время."
		regards = "С уважением,<br><strong>Команда Prianik Studio</strong>"
	}

	return fmt.Sprintf(`
<h2>%s</h2>
<p>%s</p>
<p>%s</p>
<p>%s</p>
	`, title, greeting, message, regards)
}

func (s *SendGridSender) generateContactClientText(form *models.ContactFormRequest, lang string) string {
	switch lang {
	case "en":
		return fmt.Sprintf("Dear %s, we have received your message and will contact you shortly.", form.Name)
	case "es":
		return fmt.Sprintf("Estimado/a %s, hemos recibido su mensaje y nos pondremos en contacto con usted en breve.", form.Name)
	default:
		return fmt.Sprintf("Уважаемый(ая) %s, мы получили ваше сообщение и свяжемся с вами в ближайшее время.", form.Name)
	}
}

// Контактная форма админу с локализацией
func (s *SendGridSender) generateContactAdminHTML(form *models.ContactFormRequest, lang string) string {
	var title, fromLabel, phoneLabel, messageLabel string

	switch lang {
	case "en":
		title = "New message from the contact form"
		fromLabel = "From:"
		phoneLabel = "Phone:"
		messageLabel = "Message:"
	case "es":
		title = "Nuevo mensaje del formulario de contacto"
		fromLabel = "De:"
		phoneLabel = "Teléfono:"
		messageLabel = "Mensaje:"
	default: // ru
		title = "Новое сообщение с формы обратной связи"
		fromLabel = "От:"
		phoneLabel = "Телефон:"
		messageLabel = "Сообщение:"
	}

	return fmt.Sprintf(`
<h2>%s</h2>
<p><strong>%s</strong> %s (%s)</p>
<p><strong>%s</strong> %s</p>
<p><strong>%s</strong> %s</p>
	`, title, fromLabel, form.Name, form.Email, phoneLabel, form.Phone, messageLabel, form.Message)
}

func (s *SendGridSender) generateContactAdminText(form *models.ContactFormRequest, lang string) string {
	switch lang {
	case "en":
		return fmt.Sprintf("New message from %s (%s), phone: %s. Message: %s", form.Name, form.Email, form.Phone, form.Message)
	case "es":
		return fmt.Sprintf("Nuevo mensaje de %s (%s), teléfono: %s. Mensaje: %s", form.Name, form.Email, form.Phone, form.Message)
	default:
		return fmt.Sprintf("Новое сообщение от %s (%s), телефон: %s. Сообщение: %s", form.Name, form.Email, form.Phone, form.Message)
	}
}

// Заголовки для заказов
func (s *SendGridSender) getOrderSubject(language string, orderID int64) string {
	switch language {
	case "en":
		return fmt.Sprintf("Order Confirmation #%d", orderID)
	case "es":
		return fmt.Sprintf("Confirmación de pedido #%d", orderID)
	default:
		return fmt.Sprintf("Подтверждение заказа №%d", orderID)
	}
}

func (s *SendGridSender) getAdminOrderSubject(language string, orderID int64) string {
	switch language {
	case "en":
		return fmt.Sprintf("New Order #%d", orderID)
	case "es":
		return fmt.Sprintf("Nuevo Pedido #%d", orderID)
	default:
		return fmt.Sprintf("Новый заказ №%d", orderID)
	}
}

// Заголовки для контактной формы
func (s *SendGridSender) getContactSubject(language string) string {
	switch language {
	case "en":
		return "We received your message"
	case "es":
		return "Hemos recibido su mensaje"
	default:
		return "Мы получили ваше сообщение"
	}
}

func (s *SendGridSender) getAdminContactSubject(language string) string {
	switch language {
	case "en":
		return "New message from the contact form"
	case "es":
		return "Nuevo mensaje del formulario de contacto"
	default:
		return "Новое сообщение с формы обратной связи"
	}
}
