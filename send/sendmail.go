// Sending email with *nix sendmail
package send

//	"os/exec"

const (
	DEFAULT_SENDMAIL_PATH = "/usr/sbin/sendmail"
)

/*
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

func (s *sendmail) send() error {
	cmd := exec.Command(s.path, "-t", "-i")
	cmd.Stdin = s.Mail.CreateMln()
	return cmd.Run()
}

// Implementing io.Writer
func (s *sendmail) Write(data []byte) (n int, err error) {
	n = len(data)
	s.Mail.Body = string(data)
	err = s.send()
	return
}

func (s *sendmail) Path(p string) {
	s.path = p
}
*/
