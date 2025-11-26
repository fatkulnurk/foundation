# Mailer Package

The `mailer` package provides a flexible and extensible email sending solution for Go applications. It supports multiple email delivery methods including SMTP and AWS SES (Simple Email Service).

> **Note:**  This package is still under development and may undergo API changes.

## Table of Contents

- [Features](#features)
- [Usage](#usage)
- [Implementation Details](#implementation-details)
- [Extending](#extending)

---

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

---

## Extending

You can create custom mailer implementations by implementing the Mailer interface.

### Custom Mailer Implementation

```go
type Mailer interface {
    Send(ctx context.Context, message Message) error
}
```

### Example: SendGrid Mailer

```go
type SendGridMailer struct {
    apiKey string
    client *sendgrid.Client
}

func NewSendGridMailer(apiKey string) *SendGridMailer {
    return &SendGridMailer{
        apiKey: apiKey,
        client: sendgrid.NewSendClient(apiKey),
    }
}

func (m *SendGridMailer) Send(ctx context.Context, message mailer.Message) error {
    from := mail.NewEmail(message.FromName, message.FromAddress)
    subject := message.Subject
    to := mail.NewEmail("", message.ToAddresses[0])
    
    var content *mail.Content
    if message.HTMLBody != "" {
        content = mail.NewContent("text/html", message.HTMLBody)
    } else {
        content = mail.NewContent("text/plain", message.TextBody)
    }
    
    msg := mail.NewV3MailInit(from, subject, to, content)
    
    response, err := m.client.Send(msg)
    if err != nil {
        return err
    }
    
    if response.StatusCode >= 400 {
        return fmt.Errorf("sendgrid error: %d", response.StatusCode)
    }
    
    return nil
}
```

### Example: Mailgun Mailer

```go
type MailgunMailer struct {
    domain string
    apiKey string
    mg     *mailgun.MailgunImpl
}

func NewMailgunMailer(domain, apiKey string) *MailgunMailer {
    return &MailgunMailer{
        domain: domain,
        apiKey: apiKey,
        mg:     mailgun.NewMailgun(domain, apiKey),
    }
}

func (m *MailgunMailer) Send(ctx context.Context, message mailer.Message) error {
    msg := m.mg.NewMessage(
        fmt.Sprintf("%s <%s>", message.FromName, message.FromAddress),
        message.Subject,
        message.TextBody,
        message.ToAddresses...,
    )
    
    if message.HTMLBody != "" {
        msg.SetHtml(message.HTMLBody)
    }
    
    // Add CC
    for _, cc := range message.CcAddresses {
        msg.AddCC(cc)
    }
    
    // Add BCC
    for _, bcc := range message.BccAddresses {
        msg.AddBCC(bcc)
    }
    
    // Add attachments
    for _, att := range message.Attachments {
        msg.AddBufferAttachment(att.Name, att.Content)
    }
    
    _, _, err := m.mg.Send(ctx, msg)
    return err
}
```

### Example: Queue-based Mailer

```go
type QueueMailer struct {
    queue queue.Queue
}

func NewQueueMailer(q queue.Queue) *QueueMailer {
    return &QueueMailer{queue: q}
}

func (m *QueueMailer) Send(ctx context.Context, message mailer.Message) error {
    // Enqueue email for async processing
    _, err := m.queue.Enqueue(ctx, "email:send", message,
        queue.MaxRetry(3),
        queue.Timeout(30*time.Second),
    )
    return err
}
```

### Example: Multi-provider Mailer with Fallback

```go
type MultiMailer struct {
    providers []mailer.Mailer
}

func NewMultiMailer(providers ...mailer.Mailer) *MultiMailer {
    return &MultiMailer{providers: providers}
}

func (m *MultiMailer) Send(ctx context.Context, message mailer.Message) error {
    var lastErr error
    
    for i, provider := range m.providers {
        err := provider.Send(ctx, message)
        if err == nil {
            log.Printf("Email sent successfully via provider %d", i)
            return nil
        }
        
        log.Printf("Provider %d failed: %v", i, err)
        lastErr = err
    }
    
    return fmt.Errorf("all providers failed, last error: %w", lastErr)
}
```

---