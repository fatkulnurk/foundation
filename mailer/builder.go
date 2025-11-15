package mailer

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

func joinEmails(emails []string) string {
	var buf bytes.Buffer
	for i, email := range emails {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(email)
	}
	return buf.String()
}

type InputBuildRawMessage struct {
	// Email subject
	Subject string

	// Email's plain text content, with header text/plain; charset="UTF-8"
	TextMessage string

	// Email's HTML content, with header text/html; charset="UTF-8"
	HtmlMessage string

	// Sender using for sent email, required
	Sender Sender

	// Email's recipient details
	Destination *Destination

	// Attachments to be added to the email
	Attachments []Attachment

	// Boundary for MIME parts
	Boundary string
}

type RawMessage struct {
	input InputBuildRawMessage
}

func NewRawMessage() *RawMessage {
	return &RawMessage{}
}

func (r *RawMessage) SetSubject(subject string) *RawMessage {
	r.input.Subject = subject
	return r
}

func (r *RawMessage) SetTextMessage(textMessage string) *RawMessage {
	r.input.TextMessage = textMessage
	return r
}

func (r *RawMessage) SetHtmlMessage(htmlMessage string) *RawMessage {
	r.input.HtmlMessage = htmlMessage
	return r
}

func (r *RawMessage) SetSender(sender Sender) *RawMessage {
	r.input.Sender = sender
	return r
}

func (r *RawMessage) SetDestination(destination Destination) *RawMessage {
	r.input.Destination = &destination
	return r
}

func (r *RawMessage) SetAttachments(attachments []Attachment) *RawMessage {
	r.input.Attachments = attachments
	return r
}

func (r *RawMessage) SetBoundary(boundary string) *RawMessage {
	r.input.Boundary = boundary
	return r
}

func (r *RawMessage) Build(ctx context.Context) (*bytes.Buffer, error) {
	return buildRawMessage(ctx, r.input)
}

func buildRawMessage(ctx context.Context, i InputBuildRawMessage) (*bytes.Buffer, error) {
	// MIME Message
	var rawMessage bytes.Buffer
	if i.Boundary == "" {
		i.Boundary = fmt.Sprintf("MAIN-BOUNDARY-EMAIL-%d", time.Now().Year())
	}

	mainBoundary := i.Boundary
	if i.Destination == nil || i.Destination.ToAddresses == nil {
		return nil, errors.New("destination can't be empty")
	}

	// header from formatnya "Nama <email>"
	fromSender := fmt.Sprintf("%s <%s>", i.Sender.FromName, i.Sender.FromAddress)

	// Header Utama
	rawMessage.WriteString(fmt.Sprintf("From: %s\r\n", fromSender))
	rawMessage.WriteString(fmt.Sprintf("To: %s\r\n", joinEmails(i.Destination.ToAddresses)))
	if i.Destination.CcAddresses != nil {
		rawMessage.WriteString(fmt.Sprintf("Cc: %s\r\n", joinEmails(i.Destination.CcAddresses)))
	}
	rawMessage.WriteString(fmt.Sprintf("Subject: %s\r\n", i.Subject))
	rawMessage.WriteString("MIME-Version: 1.0\r\n")

	// Content-Type berdasarkan jenis konten
	hasBody := i.TextMessage != "" || i.HtmlMessage != ""
	hasAttachments := len(i.Attachments) > 0

	if hasAttachments || (hasBody && (i.TextMessage != "" && i.HtmlMessage != "")) {
		rawMessage.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", mainBoundary)) // wajib jika tidak mau dianggap spam
	} else if i.HtmlMessage != "" {
		rawMessage.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		rawMessage.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}
	rawMessage.WriteString("\r\n")

	// Body Content
	if hasBody {
		// Jika ada attachment atau dual content, gunakan boundary
		if hasAttachments || (i.TextMessage != "" && i.HtmlMessage != "") {
			rawMessage.WriteString(fmt.Sprintf("--%s\r\n", mainBoundary))

			// Handle dual content (text + HTML)
			if i.TextMessage != "" && i.HtmlMessage != "" {
				altBoundary := "ALT-BOUNDARY-456" // wajib dikasih boundary
				rawMessage.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", altBoundary))

				// Text Part
				rawMessage.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
				rawMessage.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
				rawMessage.WriteString(i.TextMessage + "\r\n\r\n")

				// HTML Part
				rawMessage.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
				rawMessage.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
				rawMessage.WriteString(i.HtmlMessage + "\r\n\r\n")

				rawMessage.WriteString(fmt.Sprintf("--%s--\r\n\r\n", altBoundary))
			} else {
				// Single content type
				if i.TextMessage != "" {
					rawMessage.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
					rawMessage.WriteString(i.TextMessage + "\r\n\r\n")
				} else {
					rawMessage.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
					rawMessage.WriteString(i.HtmlMessage + "\r\n\r\n")
				}
			}
		} else {
			// Tanpa attachment & single content type
			if i.TextMessage != "" {
				rawMessage.WriteString(i.TextMessage + "\r\n")
			} else {
				rawMessage.WriteString(i.HtmlMessage + "\r\n")
			}
		}
	}

	// Attachments
	for _, attachment := range i.Attachments {
		rawMessage.WriteString(fmt.Sprintf("--%s\r\n", mainBoundary))
		rawMessage.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.Name))
		rawMessage.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", attachment.MimeType, attachment.Name))
		rawMessage.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

		encoder := base64.NewEncoder(base64.StdEncoding, &rawMessage)
		encoder.Write(attachment.Content)
		encoder.Close()
		rawMessage.WriteString("\r\n")
	}

	// Tutup Boundary email
	if hasAttachments || (hasBody && (i.TextMessage != "" && i.HtmlMessage != "")) {
		rawMessage.WriteString(fmt.Sprintf("--%s--\r\n", mainBoundary))
	}

	return &rawMessage, nil
}
