//Simple interface to Go smtp
package mail

import (
	"fmt"
	"net/smtp"
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
