// Sending email with *nix sendmail
package send

import (
	"bytes"
	"os/exec"
	"strings"
)

const (
	DEFAULT_SENDMAIL_PATH = "/usr/sbin/sendmail"
)

type sendmail struct {
	// Path to the sendmail binary. If emtpy, /usr/sbin/sendmail is used
	path string
	Mail *mail
}

// Returns a sendmail struct
func NewSendMail() *sendmail {
	o := new(sendmail)
	o.path = DEFAULT_SENDMAIL_PATH
	o.Mail = NewMail()
	return o
}

func (s *sendmail) Send() error {
	return s.send()
}

// Implementing io.Writer
func (s *sendmail) Write(data []byte) (n int, err error) {
	n = len(data)
	s.Mail.Body = string(data)
	err = s.send()
	return
}

func (s *sendmail) send() error {
	buf := bytes.NewBuffer(nil)
	if len(s.Mail.From) != 0 {
		buf.WriteString("From: " + s.Mail.From + EOL)
	}
	if len(s.Mail.To) != 0 {
		buf.WriteString("To: " + strings.Join(s.Mail.To, ", ") + EOL)
	}
	if len(s.Mail.Cc) != 0 {
		buf.WriteString("Cc: " + strings.Join(s.Mail.Cc, ", ") + EOL)
	}
	if len(s.Mail.Bcc) != 0 {
		buf.WriteString("Bcc: " + strings.Join(s.Mail.Bcc, ", ") + EOL)
	}
	buf.WriteString("Subject: " + s.Mail.Subject + EOL)
	buf.WriteString("MIME-Version: 1.0" + EOL)
	buf.WriteString("Content-Type: text/plain; charset=utf-8" + EOL + EOL)
	buf.WriteString(s.Mail.Body)
	cmd := exec.Command(s.path, "-t", "-i")
	cmd.Stdin = buf
	return cmd.Run()
}

func (s *sendmail) Path(p string) {
	s.path = p
}
