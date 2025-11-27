# Mailer - Email Sending Package for Go

Module for sending emails with support for multiple providers (SMTP, AWS SES) and rich content features.

## Table of Contents

- [What is Mailer?](#what-is-mailer)
- [Module Contents](#module-contents)
- [How to Use](#how-to-use)
- [Configuration](#configuration)
- [Real-World Example](#real-world-example)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)
- [Installation](#installation)
- [Dependencies](#dependencies)
- [Extending](#extending)
- [See Also](#see-also)

---

## What is Mailer?

Mailer is a flexible email sending package that provides a unified interface for sending emails through different providers. It abstracts away the complexity of different email services while providing rich features like HTML emails, attachments, and multiple recipients.

**Think of it like:**
- A universal adapter for different email services
- A way to send emails without worrying about provider-specific APIs
- A builder for complex email messages with attachments

**Use cases:**
- Sending transactional emails (welcome, password reset, etc.)
- Sending marketing emails
- Sending emails with attachments (invoices, reports, etc.)
- Multi-provider email delivery with fallback

## Module Contents

### 1. **mailer.go** - Core Interface
Defines the main interface and data structures:
- `Mailer` interface - Common interface for all email providers
- `InputSendMail` - Input structure for sending emails
- `OutputSendMail` - Output structure with message ID
- `Sender`, `Destination`, `Attachment` - Supporting structures

### 2. **smtp.go** - SMTP Implementation
SMTP email delivery using [go-mail](https://github.com/wneessen/go-mail):
- `NewSmtp` - Create SMTP client
- `NewSMTPMailer` - Create SMTP mailer instance
- Support for TLS, authentication, and attachments

### 3. **ses.go** - AWS SES Implementation
AWS SES email delivery using AWS SDK v2:
- `NewSESClient` - Create SES client
- `NewSESMailer` - Create SES mailer instance
- Support for raw messages and attachments

### 4. **builder.go** - Raw Message Builder
Build raw MIME email messages:
- `NewRawMessage` - Create message builder
- Fluent API for building complex emails
- Support for multipart messages and attachments

### 5. **config.go** - Configuration
Configuration helpers:
- `LoadSMTPConfig` - Load SMTP config from environment
- `LoadSESConfig` - Load SES config from environment

## How to Use

### Basic Usage with SMTP

```go
import "github.com/fatkulnurk/foundation/mailer"

func main() {
    // Create SMTP client
    smtpClient, err := mailer.NewSmtp(&mailer.SMTPConfig{
        Host:              "smtp.gmail.com",
        Port:              587,
        Username:          "user@example.com",
        Password:          "password",
        AuthType:          "PLAIN",
        WithTLSPortPolicy: 0, // Mandatory TLS
    })
    if err != nil {
        panic(err)
    }

    // Create mailer
    m := mailer.NewSMTPMailer(smtpClient, "sender@example.com", "Sender Name")

    // Send email
    output, err := m.SendMail(context.Background(), mailer.InputSendMail{
        Subject:     "Hello World",
        TextMessage: "This is a plain text message",
        Destination: mailer.Destination{
            ToAddresses: []string{"recipient@example.com"},
        },
    })
}
```

### Basic Usage with AWS SES

```go
// Create SES client
sesClient, err := mailer.NewSESClient(&mailer.SESConfig{
    Region: "us-west-2",
})
if err != nil {
    panic(err)
}

// Create mailer
m := mailer.NewSESMailer(sesClient, "sender@example.com", "Sender Name")

// Send email
output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Hello World",
    HtmlMessage: "<h1>Hello World</h1>",
    Destination: mailer.Destination{
        ToAddresses: []string{"recipient@example.com"},
    },
})
```

### Sending HTML Email

```go
output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Welcome!",
    HtmlMessage: "<h1>Welcome to our service!</h1><p>Thank you for signing up.</p>",
    Destination: mailer.Destination{
        ToAddresses: []string{"user@example.com"},
    },
})
```

### Sending Email with Both Text and HTML

```go
output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Newsletter",
    TextMessage: "This is the plain text version",
    HtmlMessage: "<h1>Newsletter</h1><p>This is the HTML version</p>",
    Destination: mailer.Destination{
        ToAddresses: []string{"user@example.com"},
    },
})
```

### Sending Email with Attachments

```go
// Read file
fileBytes, err := os.ReadFile("invoice.pdf")
if err != nil {
    panic(err)
}

output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Your Invoice",
    HtmlMessage: "<h1>Invoice Attached</h1>",
    Destination: mailer.Destination{
        ToAddresses: []string{"customer@example.com"},
    },
    Attachments: []mailer.Attachment{
        {
            Content:  fileBytes,
            Name:     "invoice.pdf",
            MimeType: "application/pdf",
        },
    },
})
```

### Multiple Recipients (To, CC, BCC)

```go
output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Team Update",
    HtmlMessage: "<h1>Important Update</h1>",
    Destination: mailer.Destination{
        ToAddresses:  []string{"user1@example.com", "user2@example.com"},
        CcAddresses:  []string{"manager@example.com"},
        BccAddresses: []string{"admin@example.com"},
    },
})
```

### Override Default Sender

```go
output, err := m.SendMail(context.Background(), mailer.InputSendMail{
    Subject:     "Special Notification",
    HtmlMessage: "<h1>From Support Team</h1>",
    Destination: mailer.Destination{
        ToAddresses: []string{"user@example.com"},
    },
    Sender: &mailer.Sender{
        FromAddress: "support@example.com",
        FromName:    "Support Team",
    },
})
```

### Using Raw Message Builder

```go
rawMessage := mailer.NewRawMessage().
    SetSubject("Complex Email").
    SetTextMessage("Plain text version").
    SetHtmlMessage("<h1>HTML version</h1>").
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
    }).
    SetBoundary("CUSTOM-BOUNDARY")

buffer, err := rawMessage.Build(context.Background())
if err != nil {
    panic(err)
}

// Use the raw message buffer
fmt.Println(buffer.String())
```

## Configuration

### SMTP Configuration

```go
// Manual configuration
smtpConfig := &mailer.SMTPConfig{
    Host:              "smtp.gmail.com",
    Port:              587,
    Username:          "user@example.com",
    Password:          "password",
    AuthType:          "PLAIN",
    WithTLSPortPolicy: 0, // 0 = Mandatory, 1 = Opportunistic, 2 = No TLS
}

// Load from environment variables
smtpConfig := mailer.LoadSMTPConfig()
// Reads: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, 
//        SMTP_AUTH_TYPE, SMTP_WITH_TLS_PORT_POLICY
```

### SES Configuration

```go
// Manual configuration
sesConfig := &mailer.SESConfig{
    Region: "us-west-2",
}

// Load from environment variables
sesConfig := mailer.LoadSESConfig()
// Reads: SES_REGION
```

### SMTP Auth Types

Available authentication types:
- `PLAIN` - Plain authentication
- `LOGIN` - Login authentication
- `CRAM-MD5` - CRAM-MD5 authentication
- `XOAUTH2` - OAuth2 authentication
- `SCRAM-SHA-1`, `SCRAM-SHA-256`, `SCRAM-SHA-512` - SCRAM authentication
- `NOAUTH` - No authentication

### TLS Policies

- `0` - **Mandatory**: Always use TLS (recommended)
- `1` - **Opportunistic**: Use TLS if available
- `2` - **No TLS**: Don't use TLS (not recommended)

## Real-World Example

### Transactional Email Service

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/fatkulnurk/foundation/mailer"
)

type EmailService struct {
    mailer mailer.Mailer
}

func NewEmailService() (*EmailService, error) {
    // Load config from environment
    smtpConfig := mailer.LoadSMTPConfig()
    
    // Create SMTP client
    smtpClient, err := mailer.NewSmtp(smtpConfig)
    if err != nil {
        return nil, err
    }
    
    // Create mailer
    m := mailer.NewSMTPMailer(
        smtpClient,
        "noreply@example.com",
        "Example App",
    )
    
    return &EmailService{mailer: m}, nil
}

func (s *EmailService) SendWelcomeEmail(ctx context.Context, userEmail, userName string) error {
    html := fmt.Sprintf(`
        <h1>Welcome, %s!</h1>
        <p>Thank you for signing up for our service.</p>
        <p>Get started by visiting your dashboard.</p>
    `, userName)
    
    _, err := s.mailer.SendMail(ctx, mailer.InputSendMail{
        Subject:     "Welcome to Example App",
        HtmlMessage: html,
        Destination: mailer.Destination{
            ToAddresses: []string{userEmail},
        },
    })
    
    return err
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) error {
    html := fmt.Sprintf(`
        <h1>Password Reset Request</h1>
        <p>Click the link below to reset your password:</p>
        <a href="https://example.com/reset?token=%s">Reset Password</a>
        <p>This link will expire in 1 hour.</p>
    `, resetToken)
    
    _, err := s.mailer.SendMail(ctx, mailer.InputSendMail{
        Subject:     "Password Reset Request",
        HtmlMessage: html,
        Destination: mailer.Destination{
            ToAddresses: []string{userEmail},
        },
    })
    
    return err
}

func (s *EmailService) SendInvoiceEmail(ctx context.Context, userEmail string, invoicePDF []byte) error {
    _, err := s.mailer.SendMail(ctx, mailer.InputSendMail{
        Subject:     "Your Invoice",
        HtmlMessage: "<h1>Invoice Attached</h1><p>Please find your invoice attached.</p>",
        Destination: mailer.Destination{
            ToAddresses: []string{userEmail},
        },
        Attachments: []mailer.Attachment{
            {
                Content:  invoicePDF,
                Name:     "invoice.pdf",
                MimeType: "application/pdf",
            },
        },
    })
    
    return err
}

func main() {
    emailService, err := NewEmailService()
    if err != nil {
        log.Fatal(err)
    }
    
    // Send welcome email
    err = emailService.SendWelcomeEmail(
        context.Background(),
        "user@example.com",
        "John Doe",
    )
    if err != nil {
        log.Printf("Failed to send welcome email: %v", err)
    }
}
```

## Best Practices

### 1. Use Environment Variables for Configuration

```go
// Good - load from environment
smtpConfig := mailer.LoadSMTPConfig()

// Avoid - hardcoded credentials
smtpConfig := &mailer.SMTPConfig{
    Username: "user@example.com",
    Password: "hardcoded-password", // Don't do this!
}
```

### 2. Always Provide Both Text and HTML Versions

```go
// Good - both versions
output, err := m.SendMail(ctx, mailer.InputSendMail{
    Subject:     "Newsletter",
    TextMessage: "Plain text for email clients that don't support HTML",
    HtmlMessage: "<h1>Rich HTML content</h1>",
    Destination: mailer.Destination{
        ToAddresses: []string{"user@example.com"},
    },
})
```

### 3. Handle Errors Gracefully

```go
output, err := m.SendMail(ctx, mailer.InputSendMail{
    Subject:     "Important Email",
    HtmlMessage: "<h1>Content</h1>",
    Destination: mailer.Destination{
        ToAddresses: []string{"user@example.com"},
    },
})
if err != nil {
    log.Printf("Failed to send email: %v", err)
    // Implement retry logic or fallback
    return err
}

log.Printf("Email sent successfully. Message ID: %s", *output.MessageID)
```

### 4. Use Proper MIME Types for Attachments

```go
// Good - correct MIME types
attachments := []mailer.Attachment{
    {Content: pdfBytes, Name: "doc.pdf", MimeType: "application/pdf"},
    {Content: imageBytes, Name: "image.png", MimeType: "image/png"},
    {Content: csvBytes, Name: "data.csv", MimeType: "text/csv"},
}
```

## Common Patterns

### Multi-Provider with Fallback

```go
type MultiProviderMailer struct {
    primary   mailer.Mailer
    fallback  mailer.Mailer
}

func NewMultiProviderMailer(primary, fallback mailer.Mailer) *MultiProviderMailer {
    return &MultiProviderMailer{
        primary:  primary,
        fallback: fallback,
    }
}

func (m *MultiProviderMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    // Try primary provider
    output, err := m.primary.SendMail(ctx, msg)
    if err == nil {
        return output, nil
    }
    
    log.Printf("Primary mailer failed: %v, trying fallback", err)
    
    // Try fallback provider
    return m.fallback.SendMail(ctx, msg)
}
```

### Email Template System

```go
type EmailTemplate struct {
    subject  string
    htmlTmpl *template.Template
    textTmpl *template.Template
}

func (t *EmailTemplate) Render(data interface{}) (string, string, error) {
    var htmlBuf, textBuf bytes.Buffer
    
    if err := t.htmlTmpl.Execute(&htmlBuf, data); err != nil {
        return "", "", err
    }
    
    if err := t.textTmpl.Execute(&textBuf, data); err != nil {
        return "", "", err
    }
    
    return htmlBuf.String(), textBuf.String(), nil
}

func SendTemplatedEmail(ctx context.Context, m mailer.Mailer, tmpl *EmailTemplate, to string, data interface{}) error {
    html, text, err := tmpl.Render(data)
    if err != nil {
        return err
    }
    
    _, err = m.SendMail(ctx, mailer.InputSendMail{
        Subject:     tmpl.subject,
        HtmlMessage: html,
        TextMessage: text,
        Destination: mailer.Destination{
            ToAddresses: []string{to},
        },
    })
    
    return err
}
```

### Async Email Queue

```go
type EmailQueue struct {
    mailer mailer.Mailer
    queue  chan mailer.InputSendMail
}

func NewEmailQueue(m mailer.Mailer, workers int) *EmailQueue {
    eq := &EmailQueue{
        mailer: m,
        queue:  make(chan mailer.InputSendMail, 100),
    }
    
    // Start workers
    for i := 0; i < workers; i++ {
        go eq.worker()
    }
    
    return eq
}

func (eq *EmailQueue) worker() {
    for msg := range eq.queue {
        _, err := eq.mailer.SendMail(context.Background(), msg)
        if err != nil {
            log.Printf("Failed to send email: %v", err)
            // Implement retry logic
        }
    }
}

func (eq *EmailQueue) Enqueue(msg mailer.InputSendMail) {
    eq.queue <- msg
}
```

### Rate Limited Mailer

```go
type RateLimitedMailer struct {
    mailer  mailer.Mailer
    limiter *rate.Limiter
}

func NewRateLimitedMailer(m mailer.Mailer, rps int) *RateLimitedMailer {
    return &RateLimitedMailer{
        mailer:  m,
        limiter: rate.NewLimiter(rate.Limit(rps), rps),
    }
}

func (r *RateLimitedMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    if err := r.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    return r.mailer.SendMail(ctx, msg)
}
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/mailer
```

## Dependencies

- **SMTP**: [github.com/wneessen/go-mail](https://github.com/wneessen/go-mail)
- **AWS SES**: [github.com/aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)
- **Logging**: `github.com/fatkulnurk/foundation/logging`
- **Support**: `github.com/fatkulnurk/foundation/support`

---

## Extending

You can extend the mailer by implementing custom providers or decorators.

### Custom Provider Implementation

```go
type SendGridMailer struct {
    apiKey      string
    fromAddress string
    fromName    string
}

func NewSendGridMailer(apiKey, fromAddress, fromName string) mailer.Mailer {
    return &SendGridMailer{
        apiKey:      apiKey,
        fromAddress: fromAddress,
        fromName:    fromName,
    }
}

func (s *SendGridMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    // Implement SendGrid API call
    // ...
    
    return &mailer.OutputSendMail{
        MessageID: &messageID,
    }, nil
}
```

### Logging Decorator

```go
type LoggingMailer struct {
    mailer mailer.Mailer
    logger *log.Logger
}

func NewLoggingMailer(m mailer.Mailer, logger *log.Logger) mailer.Mailer {
    return &LoggingMailer{
        mailer: m,
        logger: logger,
    }
}

func (l *LoggingMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    l.logger.Printf("Sending email to: %v, subject: %s", msg.Destination.ToAddresses, msg.Subject)
    
    output, err := l.mailer.SendMail(ctx, msg)
    
    if err != nil {
        l.logger.Printf("Failed to send email: %v", err)
    } else {
        l.logger.Printf("Email sent successfully. Message ID: %s", *output.MessageID)
    }
    
    return output, err
}
```

### Retry Decorator

```go
type RetryMailer struct {
    mailer     mailer.Mailer
    maxRetries int
    delay      time.Duration
}

func NewRetryMailer(m mailer.Mailer, maxRetries int, delay time.Duration) mailer.Mailer {
    return &RetryMailer{
        mailer:     m,
        maxRetries: maxRetries,
        delay:      delay,
    }
}

func (r *RetryMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    var lastErr error
    
    for i := 0; i < r.maxRetries; i++ {
        output, err := r.mailer.SendMail(ctx, msg)
        if err == nil {
            return output, nil
        }
        
        lastErr = err
        log.Printf("Attempt %d failed: %v", i+1, err)
        
        if i < r.maxRetries-1 {
            time.Sleep(r.delay)
        }
    }
    
    return nil, fmt.Errorf("all %d attempts failed, last error: %w", r.maxRetries, lastErr)
}
```

### Metrics Decorator

```go
type MetricsMailer struct {
    mailer        mailer.Mailer
    sentCount     int64
    failedCount   int64
    totalDuration time.Duration
    mu            sync.Mutex
}

func NewMetricsMailer(m mailer.Mailer) *MetricsMailer {
    return &MetricsMailer{
        mailer: m,
    }
}

func (m *MetricsMailer) SendMail(ctx context.Context, msg mailer.InputSendMail) (*mailer.OutputSendMail, error) {
    start := time.Now()
    output, err := m.mailer.SendMail(ctx, msg)
    duration := time.Since(start)
    
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.totalDuration += duration
    if err != nil {
        m.failedCount++
    } else {
        m.sentCount++
    }
    
    return output, err
}

func (m *MetricsMailer) GetMetrics() (sent, failed int64, avgDuration time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    total := m.sentCount + m.failedCount
    if total > 0 {
        avgDuration = m.totalDuration / time.Duration(total)
    }
    
    return m.sentCount, m.failedCount, avgDuration
}
```

---

## See Also

- `mailer.go` - Core interface and types
- `smtp.go` - SMTP implementation
- `ses.go` - AWS SES implementation
- `builder.go` - Raw message builder