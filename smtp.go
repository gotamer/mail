// gotamer/mail
// Simplifies the interface to the go smtp package, and
// creates a pipe for mail queuing.
// It now also implements io.Writer!!
//
package mail

import (
	"fmt"
	"net/smtp"
	//"os"
	"runtime"
	"time"

	"bitbucket.org/gotamer/final"
)

var (
	enumerate uint
	Pipe      = make(map[uint]Smtp)
)

type Smtp struct {
	failcount uint8
	HostName  string
	HostPort  int
	HostPass  string
	FromName  string
	FromAddr  string
	ToAddrs   []string
	Subject   string
	Body      string
}

func init() {
	go loop()
	final.Add(FinalMail)
}

func (s *Smtp) SetHostname(v string) {
	s.HostName = v
}

func (s *Smtp) SetHostport(v int) {
	s.HostPort = v
}

func (s *Smtp) SetPassword(v string) {
	s.HostPass = v
}

func (s *Smtp) SetFromName(v string) {
	s.FromName = v
}

func (s *Smtp) SetFromAddr(v string) {
	s.FromAddr = v
}

// Add multiple comma seperated
func (s *Smtp) SetToAddrs(addresses ...string) {
	for _, a := range addresses {
		s.AddToAddr(a)
	}
}

// Add one at a time
func (s *Smtp) AddToAddr(v string) {
	s.ToAddrs = append(s.ToAddrs, v)
}

func (s *Smtp) SetSubject(v string) {
	s.Subject = v
}

func (s *Smtp) SetBody(v string) {
	s.Body = v
}

// Implementing io.Writer
func (s *Smtp) Write(data []byte) (n int, err error) {
	n = len(data)
	s.Body = string(data)
	err = s.write()
	return
}

// Be aware this only works in long running applications like
// webservers. Be also aware that if you kill the server with Ctrl-c
// you may loose emailes, which are still in the pipe.
func (s Smtp) Send() {
	Pipe[ai()] = s
}

func (s Smtp) SendNow() error {
	return s.write()
}

// Sends an email and waits for the process to end, giving proper error feedback.
func (s *Smtp) write() (err error) {
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.HostName, s.HostPort),
		smtp.PlainAuth("", s.FromAddr, s.HostPass, s.HostName),
		s.FromAddr,
		s.ToAddrs,
		[]byte(fmt.Sprintf("From: %s <%s>\nSubject: %s\r\n\r\n%s\r\n", s.FromName, s.FromAddr, s.Subject, s.Body)),
	)
	return
}

// ai() Auto Increment / enumerate
func ai() uint {
	enumerate = enumerate + 1
	return enumerate
}

func send() {
	go func() {
		for i, s := range Pipe {
			delete(Pipe, i)
			if err := s.write(); err != nil {
				s.failcount++
				if s.failcount < 3 {
					s.Send()
				}
			}
			break // Sending only one per go routine
		}
	}()
	return
}

func loop() {
	var c uint8
	for {
		if c > 2 {
			c = 0
			runtime.Gosched()
		}
		send()
		time.Sleep(1 * time.Second)
		c++
	}
}

func FinalMail() {
	for i, s := range Pipe {
		if err := s.write(); err == nil {
			delete(Pipe, i)
		}
	}
}
