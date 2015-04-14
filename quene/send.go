// Mail Quene for non blocking mail sending
package quene

import (
	"log"
	"runtime"
	"time"

	"bitbucket.org/gotamer/final"
)

var (
	enumerate uint
	Pipe      map[uint]*quene
)

type quene struct {
	failcount  uint8
	failtime   time.Time
	mailsender Sender
}

// Pipe mail through a mail quene via the mail sender interface for smtp or sendmail
func NewPipe(mailer Sender) *quene {
	o := new(quene)
	o.mailsender = mailer
	return o
}

func (o *quene) Send() (err error) {
	if Pipe == nil {
		Pipe = make(map[uint]*quene)
		go o.send()
		final.Add(o.FinalMail)
	}
	Pipe[ai()] = o
	return
}

func PipeLength() int {
	return len(Pipe)
}

func (s *quene) send() {
	for {
		for i, s := range Pipe {
			if time.Since(s.failtime) > (time.Second * 15) {
				if err := s.mailsender.Send(); err != nil {
					s.failcount++
					s.failtime = time.Now()
					log.Printf("Quene sending Pipe # %d, FailCount %d, error %s", i, s.failcount, err)
				} else {
					log.Println("Send success, deleting pipe #  ", i)
					delete(Pipe, i)
				}
			}
			runtime.Gosched()
		}
		if len(Pipe) == 0 {
			Pipe = nil
			break
		}
	}
}

func (s *quene) FinalMail() {
	for i, s := range Pipe {
		if err := s.mailsender.Send(); err == nil {
			log.Println("FinalMail deleting pipe #  ", i)
			delete(Pipe, i)
		}
	}
}

// ai() Auto Increment / enumerate
func ai() uint {
	enumerate = enumerate + 1
	return enumerate
}
