package mailer

type ConfigSMTP struct {
	Host              string
	Port              int
	Username          string
	Password          string
	AuthType          string // one of => CRAM-MD5, CUSTOM, LOGIN, LOGIN-NOENC, NOAUTH, PLAIN, PLAIN-NOENC, XOAUTH2, SCRAM-SHA-1, SCRAM-SHA-1-PLUS, SCRAM-SHA-256, SCRAM-SHA-256-PLUS, SCRAM-SHA-384, SCRAM-SHA-384-PLUS, SCRAM-SHA-512, SCRAM-SHA-512-PLUS, AUTODISCOVER
	WithTLSPortPolicy int    // one of => 0 = Mandatory, 1 = Opportunistic, 2 = no tls
}

type ConfigSES struct {
	Region string
}
