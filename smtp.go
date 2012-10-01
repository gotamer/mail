// gotamer/mail  
// Simplifies the interface to the go smtp package, and
// creates a pipe for mail queuing.   
// See dos at http://www.robotamer.com/html/GoTamer/Mail.html
// 
// Please remember to call following from your main 
// until I figure out how to create a real final() method
//      defer mail.final()
package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

var (
	enumerate uint
	Pipe      = make(map[uint]Smtp)
)

type Smtp struct {
	HostName string
	HostPort int
	HostPass string
	FromName string
	FromAddr string
	ToAddrs  []string
	Subject  string
	Body     string
}

func init() {
	go loop()
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
		s.ToAddrs = append(s.ToAddrs, a)
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

func (s Smtp) Send() {
	Pipe[ai()] = s
}

func (s *Smtp) Write() (err error) {
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.HostName, s.HostPort),
		smtp.PlainAuth("", s.FromAddr, s.HostPass, s.HostName),
		s.FromAddr,
		s.ToAddrs,
		[]byte(fmt.Sprintf("From: %s <%s>\nSubject: %s\n%s\n", s.FromName, s.FromAddr, s.Subject, s.Body)),
	)
	if err != nil {
		return err
	}
	return nil
}

// ai() Auto Increment / enumerate 
func ai() uint {
	enumerate = enumerate + 1
	return enumerate
}

func loop() {
	for {
		send()
		sleep(2)
	}
}

func send() {
	go func() {
		//println("In go routine ")
		fmt.Printf("Checking %g\n", Pipe)

		for i, s := range Pipe {
			fmt.Printf("Sending %d - %s\n", i, s)
			if err := s.Write(); err == nil {
				delete(Pipe, i)
			}
			break // Sending only one per go routine
		}
	}()
	return
}

// call following from main to empty the pipe
// defer mail.final()
func final() {
	for i, s := range Pipe {
		if err := s.Write(); err == nil {
			delete(Pipe, i)
		}
	}
	os.Exit(0)
}

func sleep(sec time.Duration) {
	time.Sleep(time.Second * sec)
}
