package mail

import (
	"fmt"
	"net/smtp"
)

type Smtp struct {
	Hostname string
	HostPort int
	Password string
	FromName string
	FromAddr string
	ToAddrs  []string
	Subject  string
	Body     string
}

func (s *Smtp) SetHostname(v string) {
	s.Hostname = v
}

func (s *Smtp) SetHostport(v int) {
	s.HostPort = v
}

func (s *Smtp) SetPassword(v string) {
	s.Password = v
}

func (s *Smtp) SetFromName(v string) {
	s.FromName = v
}

func (s *Smtp) SetFromAddr(v string) {
	s.FromAddr = v
}

func (s *Smtp) SetToAddrs(addresses ...string) {
	for _, a := range addresses {
		s.ToAddrs = append(s.ToAddrs, a)
	}
}

func (s *Smtp) AddToAddr(v string) {
	s.ToAddrs = append(s.ToAddrs, v)
}

func (s *Smtp) SetSubject(v string) {
	s.Subject = v
}

func (s *Smtp) SetBody(v string) {
	s.Body = v
}

func (s *Smtp) Write() (err error) {

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Hostname, s.HostPort),
		smtp.PlainAuth("", s.FromAddr, s.Password, s.Hostname),
		s.FromAddr,
		s.ToAddrs,
		[]byte(fmt.Sprintf("From: %s <%s>\nSubject: %s\n%s\n", s.FromName, s.FromAddr, s.Subject, s.Body)),
	)
	if err != nil {
		return err
	}
	return nil
}
