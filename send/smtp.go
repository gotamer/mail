// Simplifies the interface to the go smtp package, and implements io.Writer!!
package send

import (
	"fmt"
	"net/smtp"
)

type SmtpTamer struct {
	HostName string
	HostPort int
	HostUser string
	HostPass string
	Envelop  *mail
}

// Returns a smtp struct
func NewSMTP(hostname, hostuser, hostpass string) *SmtpTamer {
	o := new(SmtpTamer)
	o.HostName = hostname
	o.HostPort = 587
	o.HostUser = hostuser
	o.HostPass = hostpass
	o.Envelop = NewMail()
	return o
}

// Implementing io.Writer
func (s *SmtpTamer) Write(data []byte) (n int, err error) {
	n = len(data)
	s.Envelop.Body = string(data)
	err = s.Send()
	return
}

// Sends an email and waits for the process to end, giving proper error feedback.
func (s *SmtpTamer) Send() (err error) {
	var format = "From: %s <%s>\r\nSubject: %s\r\n\r\n%s\r\n"
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.HostName, s.HostPort),
		smtp.PlainAuth("", s.HostUser, s.HostPass, s.HostName),
		s.Envelop.From,
		s.Envelop.To,
		[]byte(fmt.Sprintf(format, s.Envelop.FromName, s.Envelop.From, s.Envelop.Subject, s.Envelop.Body)),
	)
	return
}
