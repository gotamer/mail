// Simplifies the interface to the go smtp package, and implements io.Writer!!
package send

import (
	"fmt"
	"net/smtp"

	"github.com/gotamer/mail/envelop"
)

type SmtpTamer struct {
	HostName string
	HostPort int
	HostUser string
	HostPass string
	*envelop.Envelop
}

// Returns a smtp struct
func NewSMTP(hostname, hostuser, hostpass string) *SmtpTamer {
	o := new(SmtpTamer)
	o.HostName = hostname
	o.HostPort = 587
	o.HostUser = hostuser
	o.HostPass = hostpass
	o.Envelop = envelop.New()
	return o
}

// Implementing io.Writer
func (s *SmtpTamer) Write(body []byte) (n int, err error) {
	n = len(body)
	s.Envelop.SetBody(string(body))
	err = s.Send()
	return
}

// Sends an email and waits for the process to end, giving proper error feedback.
func (s *SmtpTamer) Send() (err error) {
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.HostName, s.HostPort),
		smtp.PlainAuth("", s.HostUser, s.HostPass, s.HostName),
		s.Envelop.From.Email,
		[]string{s.Envelop.To.Email},
		s.Envelop.Bytes(),
	)
	return
}
