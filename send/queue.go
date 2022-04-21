package send

import (
	"io/ioutil"

	"github.com/gotamer/mail/envelop"
)

const dirQueueNew = "/var/spool/queue/new/"

type Queue struct {
	Env *envelop.Envelop
}

func NewQueue() *Queue {
	var o = new(Queue)
	o.Env = envelop.New()
	return o
}

// Sends all emails to the Queue and waits for the process to end, giving proper error feedback.
func (s *Queue) Send() (err error) {
	s.uniqueTo()
	s.Env.Reset()

	for _, to := range s.Env.Tos {
		s.Env.SetTo(to.Name, to.Email)
		s.Env.Create()
		var f = dirQueueNew + s.Env.Id + ".eml"
		err = ioutil.WriteFile(f, s.Env.Bytes(), 0660)
		if err != nil {
			return
		}
		s.Env.Reset()
		var b []byte
		if b, err = envelop.Marshal(s.Env); err == nil {
			f = dirQueueNew + s.Env.Id + ".gob"
			err = ioutil.WriteFile(f, b, 0660)
		}
		if err != nil {
			return err
		}
	}
	return
}

func (o *Queue) uniqueTo() {
	if o.Env.To.Email != "" {
		o.Env.Tos = append(o.Env.Tos, o.Env.To)
	}

	var tos = o.Env.Tos
	o.Env.Tos = nil

	for _, t := range tos {
		if !o.hasEmail(t.Email) {
			o.Env.Tos = append(o.Env.Tos, t)
		}
	}
}
func (o *Queue) hasEmail(email string) (ok bool) {
	for _, t := range o.Env.Tos {
		if t.Email == email {
			return true
		}
	}
	return
}
