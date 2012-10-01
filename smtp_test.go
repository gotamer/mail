package mail

import (
	"gotamer/cfg"
	"testing"
)

func TestSmtp(t *testing.T) {
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
