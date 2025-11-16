package mailer

import "github.com/fatkulnurk/foundation/support"

type ConfigSMTP struct {
	Host              string
	Port              int
	Username          string
	Password          string
	AuthType          string // one of => CRAM-MD5, CUSTOM, LOGIN, LOGIN-NOENC, NOAUTH, PLAIN, PLAIN-NOENC, XOAUTH2, SCRAM-SHA-1, SCRAM-SHA-1-PLUS, SCRAM-SHA-256, SCRAM-SHA-256-PLUS, SCRAM-SHA-384, SCRAM-SHA-384-PLUS, SCRAM-SHA-512, SCRAM-SHA-512-PLUS, AUTODISCOVER
	WithTLSPortPolicy int    // one of => 0 = Mandatory, 1 = Opportunistic, 2 = no tls
}

func NewConfigSMTP() *ConfigSMTP {
	return &ConfigSMTP{
		Host:              support.GetEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:              support.GetIntEnv("SMTP_PORT", 587),
		Username:          support.GetEnv("SMTP_USERNAME", ""),
		Password:          support.GetEnv("SMTP_PASSWORD", ""),
		AuthType:          support.GetEnv("SMTP_AUTH_TYPE", "PLAIN"),
		WithTLSPortPolicy: support.GetIntEnv("SMTP_WITH_TLS_PORT_POLICY", 0),
	}
}

type ConfigSES struct {
	Region string
}

func NewConfigSES() *ConfigSES {
	return &ConfigSES{
		Region: support.GetEnv("SES_REGION", "us-east-1"),
	}
}
