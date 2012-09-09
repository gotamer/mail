package mail

import "testing"

func TestSmtp(t *testing.T) {
	s := new(Smtp)
	s.SetHostname("smtp.gmail.com")
	s.SetHostport(587)
	s.SetFromName("GoTamer")
	s.SetFromAddr("xxxx@gmail.com")
	s.SetPassword("********")
	s.AddToAddr("xxxx@yahoo.com")
	s.SetSubject("GoTamer test mail")
	s.SetBody("Let's see if we get this")
	if err := s.Write(); err != nil {
		t.Error(err)
	}
}
