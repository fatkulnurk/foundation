package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/fatkulnurk/gostarter/pkg/config"
	"github.com/fatkulnurk/gostarter/pkg/logging"

	"github.com/wneessen/go-mail"
)

func NewSmtp(cfg *config.SMTP) (*mail.Client, error) {
	// Deliver the mails via SMTP
	client, err := mail.NewClient(cfg.Host,
		mail.WithSMTPAuth(mail.SMTPAuthType(cfg.AuthType)),
		mail.WithTLSPortPolicy(mail.TLSPolicy(cfg.WithTLSPortPolicy)),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
		mail.WithPort(cfg.Port),
	)

	if err != nil {
		logging.Error(context.Background(), fmt.Sprintf("Error creating SMTP client: %s", err))
		return nil, err
	}

	return client, nil
}

type SMTPMailer struct {
	client                *mail.Client
	fromAddress, fromName string
}

func NewSMTPMailer(client *mail.Client, fromAddress, FromName string) Mailer {
	return &SMTPMailer{
		client:      client,
		fromName:    FromName,
		fromAddress: fromAddress,
	}
}

func (s SMTPMailer) SendMail(ctx context.Context, msg InputSendMail) (*OutputSendMail, error) {
	senderName := s.fromName
	senderAddress := s.fromAddress
	if msg.Sender != nil && msg.Sender.FromName != "" {
		senderName = msg.Sender.FromName
	}

	if msg.Sender != nil && msg.Sender.FromAddress != "" {
		senderAddress = msg.Sender.FromAddress
	}

	if msg.Destination.ToAddresses == nil || len(msg.Destination.ToAddresses) == 0 {
		return nil, errors.New("destination can't be empty")
	}

	message := mail.NewMsg()
	if err := message.FromFormat(senderName, senderAddress); err != nil {
		logging.Error(context.Background(), fmt.Sprintf("failed to set FROM address: %s", err))
		return nil, err
	}
	if err := message.To(msg.Destination.ToAddresses...); err != nil {
		logging.Error(context.Background(), fmt.Sprintf("failed to set TO address: %s", err))
		return nil, err
	}

	if msg.Destination.CcAddresses != nil {
		err := message.Cc(msg.Destination.CcAddresses...)
		if err != nil {
			return nil, err
		}
	}

	if msg.Destination.BccAddresses != nil {
		err := message.Bcc(msg.Destination.BccAddresses...)
		if err != nil {
			return nil, err
		}
	}

	message.Subject(msg.Subject)
	if msg.HtmlMessage != "" {
		message.SetBodyString(mail.TypeTextHTML, msg.HtmlMessage)
	} else {
		message.SetBodyString(mail.TypeTextPlain, msg.TextMessage)
	}

	if msg.Attachments != nil {
		for _, attachment := range msg.Attachments {
			message.AttachReadSeeker(attachment.Name, bytes.NewReader(attachment.Content))
		}
	}

	if err := s.client.DialAndSend(message); err != nil {
		logging.Error(context.Background(), fmt.Sprintf("failed to deliver mail: %s", err))
		return nil, err
	}

	return &OutputSendMail{}, nil
}
