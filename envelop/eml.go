package envelop

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

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

func NewEnvelop() *envelop {
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

// Set mail To addresses
func (o *envelop) AddTo(name, addr string) {
	a := Address{name, addr}
	o.To = sliceIt(o.To, a)
}

// Set mail To addresses
func (o *envelop) AddCc(name, addr string) {
	a := Address{name, addr}
	o.Cc = sliceIt(o.Cc, a)
}

// Set mail To addresses
func (o *envelop) AddBcc(name, addr string) {
	a := Address{name, addr}
	o.Bcc = sliceIt(o.Bcc, a)
}

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
	env.AddHeader("Message-ID", fmt.Sprintf("<%s@mailtamer.tamer.pw>", sRand(10, 16, false)))

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
