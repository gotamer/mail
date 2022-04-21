package send

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gotamer/mail/envelop"
)

func TestQueue(t *testing.T) {
	fmt.Println("TestQueue Started")
	var o = NewQueue()
	o.Env.SetFrom("", "me@example.org")
	o.Env.SetTo("", "you@example.org")
	o.Env.AddTo("", "him@example.org")
	o.Env.SetSubject("test mail")
	o.Env.SetBody("Very long mail about email")
	o.Env.Create()
	fmt.Println("o: ", o)
	fmt.Println(o.Env.String())
	err := o.Send()
	if err != nil {
		t.Fatalf(`Queue Send err: %s`, err.Error())
	} else {
		t.Log("No error")
	}
}

func TestQueueSend(t *testing.T) {
	fmt.Println("TestQueueSend Started")
	files, err := os.ReadDir(dirQueue)

	if err != nil {
		t.Fatalf(`Queue Read err: %s`, err.Error())
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gob" {
			var b []byte
			if b, err = os.ReadFile(dirQueue + file.Name()); err != nil {
				t.Fatalf(`Queue Read File err: %s`, err.Error())
			}
			var s = NewSMTP(hostname, hostuser, hostpass)
			if err = envelop.Unmarshal(b, s.Envelop); err != nil {
				t.Fatalf(`Queue Read File err: %s`, err.Error())
			}
			s.Envelop.Create()
			if err = s.Send(); err != nil {
				t.Fatalf(`Queue Send(): %s`, err.Error())
			}
			s.Envelop.Reset()
		}
	}
}
