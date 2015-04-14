package send

import "testing"

var (
	googleuser = "xxxx@gmail.com"
	googlepass = "***"
	otheruser  = "xxx@xxxx.net"
)

func TestSendmail(t *testing.T) {

	sm := NewSendMail()
	sm.m.From = "gotamer"
	sm.m.Subject = "Sendmail Test"
	sm.m.Body = "This is a Test from Go Sendmail"
	sm.m.SetTo(otheruser)
	if err := sm.Send(); err != nil {
		t.Errorf("Send error: %s", err)
	}
}

func TestSMTP(t *testing.T) {

	sm := NewSMTP("smtp.gmail.com", googleuser, googlepass)
	sm.m.FromName = "GoTamer"
	sm.m.From = googleuser
	sm.m.Subject = "SMTP sendmail Test"
	sm.m.Body = "This is a Test from Go Sendmail via SMTP"
	sm.m.SetTo(otheruser)
	if err := sm.Send(); err != nil {
		t.Errorf("Send error: %s", err)
	}
}
