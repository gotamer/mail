// Light weight email envelop creater. You can use this for sendmail, smtp, or to create .eml files
// Works only for text emails, no html, attachments, mime etc...
package envelop

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const EOL = "\r\n"

type Address struct {
	Name, Email string
}

type Envelop struct {
	Id          string
	From        Address
	To          Address
	Tos         []Address // List of recipients
	Subject     string
	subjectLong string
	Body        string
	b           bytes.Buffer
	ok          bool
}

func New() *Envelop {
	var env = &Envelop{}
	env.Id = sRand(10, 16, false)
	return env
}

func NewAll(from, to, subject, body string) *Envelop {
	var env = New()
	if from != "" {
		env.SetFrom("", from)
	}
	env.SetTo("", to)
	env.SetSubject(subject)
	env.SetBody(body)
	env.ok = true
	return env
}

func (env *Envelop) String() string {
	return env.b.String()
}

func (env *Envelop) Bytes() []byte {
	return env.b.Bytes()
}

func (env *Envelop) Reset() {
	env.b.Reset()
}

func (env *Envelop) Create() {
	env.Id = sRand(10, 16, false)
	env.AddHeaderLine("MIME-Version", "1.0")
	env.AddHeaderLine("Content-Type", "text/plain; charset=UTF-8")
	env.AddHeaderLine("Message-ID", fmt.Sprintf("<%s@%s>", env.Id, env.hostname()))
	//env.AddHeaderLine("Date", time.Now().UTC().Format(time.ANSIC))
	env.AddHeaderLine("Date", time.Now().UTC().Format(time.RFC1123Z))
	env.AddHeaderLine("Return-Path", env.From.Email)
	env.AddHeaderLine("From", fmt.Sprintf("%s <%s>", env.From.Name, env.From.Email))
	env.AddHeaderLine("To", fmt.Sprintf("%s <%s>", env.To.Name, env.To.Email))
	env.AddHeaderLine("Subject", env.Subject)
	if env.subjectLong != "" {
		env.b.WriteString(EOL)
		env.b.WriteString(EOL)
		env.b.WriteString(env.subjectLong)
	}
	env.b.WriteString(EOL)
	env.b.WriteString(EOL)
	env.b.WriteString(env.Body)
	env.b.WriteString(EOL)
}

// Set From mail addresses
// You can skip this if you are sending via an email provider suck as gmail
func (env *Envelop) SetFrom(name, email string) {
	if name == "" {
		name = strings.Split(email, "@")[0]
	}
	env.From.Name = name
	env.From.Email = email
}

// Set To mail addresses
func (env *Envelop) SetTo(name, email string) {
	env.To.Name = name
	env.To.Email = email
	if name == "" {
		env.To.Name = strings.Split(email, "@")[0]
	}
}

// Add To addresses (you can call this multiple times to add multiple receipiens only when you use the queue)
func (o *Envelop) AddTo(name, addr string) {
	var a Address
	a.Name = name
	a.Email = addr
	o.Tos = append(o.Tos, a)
}

// Set subject line
func (env *Envelop) SetSubject(subject string) {
	subject = strings.TrimSpace(subject)
	if len(subject) > 76 {
		env.subjectLong = subject
		env.Subject = subject[:68] + "..."
	} else {
		env.Subject = subject
	}
}

// Set body text
// This must be set last, the order of the other Set... are not importent
func (env *Envelop) SetBody(body string) {
	env.Body = strings.TrimSpace(body)
}

// Add custom header (should start with X- such as X-MyHeader) or overwrite the internal Standard headers
func (env *Envelop) AddHeaderLine(token, text string) {
	env.b.WriteString(token)
	env.b.WriteString(": ")
	env.b.WriteString(text)
	env.b.WriteString(EOL)
}

func (env *Envelop) hostname() string {
	host, err := os.Hostname()
	if err != nil {
		host = "example.org"
	}
	return host
}

// generates a random string
func sRand(min, max int, readable bool) string {
	var length int
	var charLength int
	var char []rune

	if min < max {
		length = min + rand.Intn(max-min)
	} else if max < min {
		length = max + rand.Intn(min-max)
	} else {
		length = min
	}

	if readable {
		char = []rune("ABCDEFHJLMNQRTUVWXYZabcefghijkmnopqrtuvwxyz23479")
	} else {
		char = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}
	charLength = len(char)
	buf := make([]rune, length)
	for i := range buf {
		buf[i] = char[rand.Intn(charLength)]
	}
	return string(buf)
}

// Marshal takes a populated struct and returns a []byte in GOB format
func Marshal(o *Envelop) (b []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(o)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// Unmarshal takes an empty struct as well as the []byte returned by Marshal(),
// and populates the given empty struct from the []byte.
func Unmarshal(b []byte, o *Envelop) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write(b)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(o)

	return
}
