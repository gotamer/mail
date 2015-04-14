package quene

import (
	"testing"
	"time"

	"bitbucket.org/gotamer/mail/send"
)

var (
	googleuser = "xxxx@gmail.com"
	googlepass = "***"
	otheruser  = "xxx@xxxx.net"
)

func TestPipe(t *testing.T) {
	d := time.Now()

	sm := NewSendMail()
	sm.m.From = otheruser
	sm.m.Subject = "Sendmail Quene Test"
	sm.m.Body = "This is a Test from Go Sendmail via the quene pipe"
	sm.m.SetTo(otheruser)
	qm := NewPipe(sm)
	qm.Send()

	smt := NewSMTP("smtp.gmail.com", googleuser, googlepass)
	smt.m.FromName = "GoTamer"
	smt.m.From = googleuser
	smt.m.Subject = "SMTP sendmail quene Test"
	smt.m.Body = "This is a Test from Go Sendmail via SMTP via quene pipe"
	smt.m.SetTo(otheruser)
	qm = NewPipe(smt)
	qm.Send()
	for {
		if PipeLength() == 0 {
			break
		}
		if time.Since(d) > (time.Minute * 2) {
			break
		}
	}
	if PipeLength() != 0 {
		t.Errorf("FinalMail error, pipe not empty. Length: %d", PipeLength())
	}
}
