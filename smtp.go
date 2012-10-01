// gotamer/mail  
// Simplifies the interface to the go smtp package, and
// creates a pipe for mail queuing.   
// See docs at http://www.robotamer.com/html/GoTamer/Mail.html
// 
// Please remember to call following from your main() 
// until I figure out how to create a real final() method that compliments init()
//      defer mail.final()
package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
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
	//go signalCatcher()
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

// Be aware this only works in long running applications like
// webservers. Be also aware that if you kill the server with Ctrl-c 
// you may loose emailes, which are still in the pipe.
// They will be lost for good! Use Write() if you need to be sure! 
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

func send() {
	go func() {
		for i, s := range Pipe {
			delete(Pipe, i)
			if err := s.Write(); err != nil {
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
	for {
		runtime.Gosched()
		send()
		sleep(1)
	}
}

func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	Final()
	os.Exit(0)
}

func Final() {
	for i, s := range Pipe {
		if err := s.Write(); err == nil {
			delete(Pipe, i)
		}
	}
}

func sleep(sec time.Duration) {
	time.Sleep(time.Second * sec)
}
