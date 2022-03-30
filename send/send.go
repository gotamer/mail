// sendmail, SMTP, Mail Quene & mail with io.Writer
package send

// Also check out add on github.com/gotamer/mail/quene for non blocking mail sending

const (
	EOL = "\r\n"
)

// Sender is an interface with a Send method, that dispatches a single Mail
type Sender interface {
	Send() error
}

// mail represents an eMail
type mail struct {

	// Sender Name
	FromName string
	// Sender email address
	From string

	// List of recipients
	To  []string
	Cc  []string
	Bcc []string

	// Mail subject as UTF-8 string
	Subject string

	// Headers are the headers
	Headers map[string]string

	// Body provides the actual body of the mail. It has to be UTF-8 encoded,
	// or you must set the Content-Type Header
	Body string
}

// Returns a new mail struct with Headers initialized to an empty map
func NewMail() *mail {
	m := new(mail)
	m.Headers = make(map[string]string)
	m.SetHeader("X-Mailer", "MailTamer @ github.com/gotamer/mail")
	return m
}

// Set from mail To addresses
func (o *mail) SetFrom(name, email string) {
	o.FromName = name
	o.From = email
}

// Set mail To addresses
func (o *mail) SetTo(addresses ...string) {
	o.To = sliceIt(o.To, addresses)
}

// Set mail cc addresses
func (o *mail) SetCc(addresses ...string) {
	o.Cc = sliceIt(o.Cc, addresses)
}

// Set mail bcc addresses
func (o *mail) SetBcc(addresses ...string) {
	o.Bcc = sliceIt(o.Bcc, addresses)
}

// Set custom headers
func (o *mail) SetHeader(k, v string) {
	o.Headers[k] = v
}

func sliceIt(slice, add []string) []string {
	if len(slice) == 0 {
		return add
	}
	for _, a := range add {
		slice = append(slice, a)
	}
	return slice
}
