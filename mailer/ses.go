package mailer

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/fatkulnurk/gostarter/pkg/config"
	"github.com/fatkulnurk/gostarter/pkg/logging"
)

func NewSESClient(cfg *config.SES) (*sesv2.Client, error) {
	awscfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		logging.Error(context.Background(), fmt.Sprintf("unable to load SDK config, %v", err))
		return nil, err
	}

	return sesv2.NewFromConfig(awscfg), nil
}

type SESMailer struct {
	client                *sesv2.Client
	fromAddress, fromName string
}

func NewSESMailer(client *sesv2.Client, fromAddress string, fromName string) Mailer {
	return &SESMailer{
		client:      client,
		fromAddress: fromAddress,
		fromName:    fromName,
	}
}

func (s *SESMailer) SendMail(ctx context.Context, msg InputSendMail) (*OutputSendMail, error) {
	fromEmailAddress := s.fromAddress
	if msg.Sender != nil && msg.Sender.FromAddress != "" {
		fromEmailAddress = msg.Sender.FromAddress
	}
	fromName := s.fromName
	if msg.Sender != nil && msg.Sender.FromName != "" {
		fromName = msg.Sender.FromName
	}

	if msg.Destination.ToAddresses == nil {
		return nil, errors.New("destination can't be empty")
	}

	fromSender := fmt.Sprintf("%s <%s>", fromName, fromEmailAddress)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(fromSender),
		Destination: &types.Destination{
			ToAddresses:  msg.Destination.ToAddresses,
			CcAddresses:  msg.Destination.CcAddresses,
			BccAddresses: msg.Destination.BccAddresses,
		},
		Content: nil,
	}

	if msg.Attachments != nil {
		rawMessage, err := buildRawMessage(ctx, InputBuildRawMessage{
			Subject:     msg.Subject,
			TextMessage: msg.TextMessage,
			HtmlMessage: msg.HtmlMessage,
			Sender: Sender{
				FromAddress: fromEmailAddress,
				FromName:    fromName,
			},
			Destination: &msg.Destination,
			Attachments: msg.Attachments,
			Boundary:    msg.Boundary,
		})
		if err != nil {
			return nil, err
		}
		input.Content = &types.EmailContent{
			Raw: &types.RawMessage{
				Data: rawMessage.Bytes(),
			},
		}
	} else {
		input.Content = &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: &msg.TextMessage,
					},
					Html: &types.Content{
						Data: &msg.HtmlMessage,
					},
				},
				Subject: &types.Content{
					Data: &msg.Subject,
				},
			},
		}
	}

	res, err := s.client.SendEmail(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &OutputSendMail{MessageID: res.MessageId}, err
}
