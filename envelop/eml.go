// Light weight email envelop creater. You can use this for sendmail, smtp, or to create .eml files
// Works only for text emails, no html, attachments, mime etc...
package envelop

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var MessageIdLink = "gotamer.tamer.pw"

type Address struct {
	Name, Address string
}

type envelop struct {
	FromName string
	FromAddr string

	// List of recipients
	To  []Address
	Cc  []Address
	Bcc []Address

	Subject string
	Body    string
	b       bytes.Buffer
	ok      bool
}

func New() *envelop {
	return new(envelop)
}

func (env *envelop) run() {
	env.addHeader()
	env.addFooter()
	env.ok = true
}

func (env *envelop) String() string {
	if env.ok == false {
		env.addHeader()
		env.addFooter()
		env.ok = true
	}
	return env.b.String()
}

func (env *envelop) Bytes() []byte {
	if env.ok == false {
		env.addHeader()
		env.addFooter()
		env.ok = true
	}
	return env.b.Bytes()
}

// Set from mail To addresses
func (o *envelop) SetFrom(name, email string) {
	o.FromName = name
	o.FromAddr = email
}

// Add To addresses (called Addxx, since you can call this multiple times to add multiple receipiens)
func (o *envelop) AddTo(name, addr string) {
	a := Address{name, addr}
	o.To = sliceIt(o.To, a)
}

// Add Cc addresses
func (o *envelop) AddCc(name, addr string) {
	a := Address{name, addr}
	o.Cc = sliceIt(o.Cc, a)
}

// Add Bcc addresses
func (o *envelop) AddBcc(name, addr string) {
	a := Address{name, addr}
	o.Bcc = sliceIt(o.Bcc, a)
}

// Add custom header (should start with X- such as X-MyHeader)
func (env *envelop) AddHeader(token, text string) {
	env.b.WriteString(token)
	env.b.WriteString(": ")
	env.b.WriteString(text)
	env.b.WriteString("\r\n")
}

func sliceIt(slice []Address, add Address) []Address {
	slice = append(slice, add)
	return slice
}

func (env *envelop) addHeader() {
	env.AddHeader("MIME-Version", "1.0")
	env.AddHeader("Content-Type", "text/plain; charset=UTF-8")
	env.AddHeader("Message-ID", fmt.Sprintf("<%s@%s>", sRand(10, 16, false), MessageIdLink))

	for _, v := range env.To {
		if v.Name == "" {
			v.Name = strings.Split(v.Address, "@")[0]
		}
		env.AddHeader("To", fmt.Sprintf("%s <%s>", v.Name, v.Address))
	}
	for _, v := range env.Cc {
		if v.Name == "" {
			v.Name = strings.Split(v.Address, "@")[0]
		}
		env.AddHeader("Cc", fmt.Sprintf("%s <%s>", v.Name, v.Address))
	}
	for _, v := range env.Bcc {
		if v.Name == "" {
			v.Name = strings.Split(v.Address, "@")[0]
		}
		env.AddHeader("Bcc", fmt.Sprintf("%s <%s>", v.Name, v.Address))
	}
	env.AddHeader("Date", time.Now().UTC().Format(time.ANSIC))
	if env.FromName == "" {
		env.FromName = strings.Split(env.FromAddr, "@")[0]
	}
	env.AddHeader("From", fmt.Sprintf("%s <%s>", env.FromName, env.FromAddr))
	env.AddHeader("Subject", env.Subject)
}
func (env *envelop) addFooter() {
	env.b.WriteString("\r\n")
	env.b.WriteString("\r\n")
	env.b.WriteString(env.Body)
	env.b.WriteString("\r\n")
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
