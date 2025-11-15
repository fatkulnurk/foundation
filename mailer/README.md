# Mailer Package

The `mailer` package provides a flexible and extensible email sending solution for Go applications. It supports multiple email delivery methods including SMTP and AWS SES (Simple Email Service).

> **Note:** ⚠️ This package is still under development and may undergo API changes.

## Features

- Multiple email delivery providers:
  - SMTP via [go-mail](https://github.com/wneessen/go-mail)
  - AWS SES (Simple Email Service) via AWS SDK v2
- Rich email content support:
  - Plain text emails
  - HTML emails
  - Mixed content (both text and HTML)
  - File attachments
- Flexible recipient management:
  - To, CC, and BCC recipients
- Customizable sender information
- Raw message building for advanced use cases

## Usage

### Interface

The package defines a common interface `IMailer` that all email delivery implementations must satisfy:

```go
type IMailer interface {
	SendMail(ctx context.Context, msg InputSendMail) (*OutputSendMail, error)
}
```

### Creating an SMTP Mailer

```go
// Create SMTP client
smtpClient, err := mailer.NewSmtp(&config.SMTP{
    Host:             "smtp.example.com",
    Port:             587,
    Username:         "user@example.com",
    Password:         "password",
    AuthType:         1, // Use appropriate auth type
    WithTLSPortPolicy: 2, // Use appropriate TLS policy
})
if err != nil {
    // Handle error
}

// Create SMTP mailer with default sender
smtpMailer := mailer.NewSMTPMailer(smtpClient, "sender@example.com", "Sender Name")
```

### Creating an AWS SES Mailer

```go
// Create SES client
sesClient, err := mailer.NewSESClient(&config.SES{})
if err != nil {
    // Handle error
}

// Create SES mailer with default sender
sesMailer := mailer.NewSESMailer(sesClient, "sender@example.com", "Sender Name")
```

### Sending an Email

```go
output, err := mailer.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Hello World",
    TextMessage: "This is a plain text message",
    HtmlMessage: "<h1>Hello World</h1><p>This is an HTML message</p>",
    Destination: mailer.Destination{
        ToAddresses:  []string{"recipient@example.com"},
        CcAddresses:  []string{"cc@example.com"},
        BccAddresses: []string{"bcc@example.com"},
    },
    Attachments: []mailer.Attachment{
        {
            Content:  fileBytes,
            Name:     "document.pdf",
            MimeType: "application/pdf",
        },
    },
    // Optional: override default sender
    Sender: &mailer.Sender{
        FromAddress: "custom@example.com",
        FromName:    "Custom Sender",
    },
})
```

### Raw Message Building

For advanced use cases, you can build raw email messages:

```go
rawMessage := mailer.NewRawMessage().
    SetSubject("Hello World").
    SetTextMessage("This is a plain text message").
    SetHtmlMessage("<h1>Hello World</h1>").
    SetSender(mailer.Sender{
        FromAddress: "sender@example.com",
        FromName:    "Sender Name",
    }).
    SetDestination(mailer.Destination{
        ToAddresses: []string{"recipient@example.com"},
    }).
    SetAttachments([]mailer.Attachment{
        {
            Content:  fileBytes,
            Name:     "document.pdf",
            MimeType: "application/pdf",
        },
    })

buffer, err := rawMessage.Build(context.Background())
// Use the raw message buffer
```

## Implementation Details

- `mailer.go`: Defines the core interfaces and data structures
- `builder.go`: Provides functionality for building raw email messages
- `smtp.go`: Implements email delivery via SMTP
- `ses.go`: Implements email delivery via AWS SES