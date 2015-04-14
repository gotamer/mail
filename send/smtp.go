// Simplifies the interface to the go smtp package, and
// also implements io.Writer!!
package send

import (
	"fmt"
	"net/smtp"
)

type sendsmtp struct {
	HostName string
	HostPort int
	HostUser string
	HostPass string
	Mail     *mail
}

// Returns a smtp struct
func NewSMTP(hostname, hostuser, hostpass string) *sendsmtp {
	o := new(sendsmtp)
	o.HostName = hostname
	o.HostPort = 587
	o.HostUser = hostuser
	o.HostPass = hostpass
	o.Mail = NewMail()
	return o
}

// Implementing io.Writer
func (s *sendsmtp) Write(data []byte) (n int, err error) {
	n = len(data)
	s.Mail.Body = string(data)
	err = s.Send()
	return
}

// Sends an email and waits for the process to end, giving proper error feedback.
func (s *sendsmtp) Send() (err error) {
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.HostName, s.HostPort),
		smtp.PlainAuth("", s.HostUser, s.HostPass, s.HostName),
		s.Mail.From,
		s.Mail.To,
		[]byte(fmt.Sprintf("From: %s <%s>\r\nSubject: %s\r\n\r\n%s\r\n", s.Mail.FromName, s.Mail.From, s.Mail.Subject, s.Mail.Body)),
	)
	return
}
