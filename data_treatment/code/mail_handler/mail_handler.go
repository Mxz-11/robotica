package mail_handler

import (
	"data_treatment/log_handler"

	"github.com/wneessen/go-mail"
)

func SendMail(sender string, senderPassword string, receptors []string, smtpServer string, subject string, body string) {
	client, err := mail.NewClient(smtpServer, mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(sender), mail.WithPassword(senderPassword))
	if err != nil {
		log_handler.Error("Error creating SMTP client: %s", err)
		return
	}

	m := mail.NewMsg()
	m.From(sender)
	if len(receptors) == 0 {
		log_handler.Error("No receptors provided")
		return
	}
	for _, r := range receptors {
		m.To(r)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)

	if err := client.DialAndSend(m); err != nil {
		log_handler.Error("Error sending email: %s", err)
		return
	}
	log_handler.Success("Mail sent correctly")
}
