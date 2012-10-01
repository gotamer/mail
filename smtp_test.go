package mail

import (
	"gotamer/cfg"
	"testing"
)

func TestSmtpWrite(t *testing.T) {
	cfg.APPL = "mail"
	cfg.Name = "smtp"
	s := new(Smtp)
	if err := cfg.Get(&s); err != nil {
		t.Errorf("mail smtp error: %v\n", err)
	}
	s.SetSubject("GoTamer test mail")
	s.SetBody("Let's see if we get this")
	if err := s.Write(); err != nil {
		t.Error(err)
	}
}

func TestSmtpSend(t *testing.T) {
	cfg.APPL = "mail"
	cfg.Name = "smtp"
	s := new(Smtp)
	if err := cfg.Get(&s); err != nil {
		t.Errorf("mail smtp error: %v\n", err)
	}
	s.SetSubject("GoTamer test Send 1")
	s.SetBody("Let's see if we get this")
	s2 := new(Smtp)
	if err := cfg.Get(&s2); err != nil {
		t.Errorf("mail smtp error: %v\n", err)
	}
	s2.SetSubject("GoTamer test Send 2")
	s2.SetBody("Let's see if we get this")
	s.Send()
	s2.Send()
	sleep(60)
}
