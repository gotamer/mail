package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gotamer/mail/envelop"
	"github.com/gotamer/mail/send"
)

const (
	EXT_GOB     = ".gob"
	EXT_EML     = ".eml"
	FILE_CONFIG = "/etc/sendmail.json"
	DIR_QUEUE   = "/var/spool/queue/"
)

var Cfg struct {
	Smtp struct {
		Hostname string // smtp.example.com
		Hostport int    // 587
		Username string // usally email address username@example.com
		Password string
	}
}

var (
	Info     = *log.Default()
	hostname = "mail.example.org"
	to       = flag.String("t", "", "Mail To email address")
	from     = flag.String("f", "", "Mail From email address")
	subject  = flag.String("s", "", "Mail subject line")
	body     = flag.String("b", "", "Mail body text")
	cfg_file = flag.String("c", FILE_CONFIG, "Config file location")
	help     = flag.Bool("h", false, "Display help")
	runQueue = flag.Bool("run", false, "Process Queue")
	now      = flag.Bool("now", false, "Send message now, message will not be added to Queue")
	debug    = flag.Bool("d", false, "Debug mode")
)

type env struct {
	env *envelop.Envelop
}

func init() {
	var err error
	log.SetPrefix("sendmail ERR ")
	Info.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.Parse()
	// || *to == "" || *subject == "" || *body == "" || *runQueue == false
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *cfg_file == "" {
		*cfg_file = FILE_CONFIG
	}

	if *debug {
		Info.SetPrefix("sendmail INF ")
		Info.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		Info.SetOutput(ioutil.Discard)
	}

	if err = LoadJson(*cfg_file, &Cfg); err != nil {
		log.Println("cfg.Load error: ", err.Error())
		if err = SaveJson(*cfg_file, &Cfg); err != nil {
			log.Println("cfg.Save error: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("\n\n++++++++++++++++++++++++++++++++++++")
			fmt.Println("\n\tPlease edit your config file at:\n\n\t", *cfg_file)
			fmt.Println("\n\n++++++++++++++++++++++++++++++++++++")
			os.Exit(0)
		}
	}

	if h, err := os.Hostname(); err != nil {
		log.Println(err)
	} else {
		hostname = h
	}
	Info.Println("hostname: ", hostname)
}

func main() {
	var err error
	var e = new(env)

	if *runQueue {
		Info.Println("Runing Queue")
		e.processQueue()
		Info.Println("Queue Finished")
		os.Exit(0)
	}

	e.env = envelop.New()
	if *from != "" {
		e.env.SetFrom("", *from)
	} else {
		e.env.SetFrom(hostname, Cfg.Smtp.Username)
	}
	e.env.SetSubject(*subject)
	e.env.SetBody(*body)
	e.env.SetTo("", *to)

	Info.Printf("smtphost: %s port: %d\n", Cfg.Smtp.Hostname, Cfg.Smtp.Hostport)
	Info.Printf("From: %s To: %v\n", e.env.From, e.env.To)

	if *now {
		e.env.Create()
		if err = e.sendNow(); err != nil {
			log.Printf("Send Now error: %s", err.Error())
		}
	} else {
		if err = e.sendQueue(); err != nil {
			log.Printf(`Send Queue err: %s`, err.Error())
		}
	}
	Info.Println("envelop:\n", e.env.String())
	e.env.Reset()
}

func (e *env) sendNow() (err error) {
	o := send.NewSMTP(Cfg.Smtp.Hostname, Cfg.Smtp.Username, Cfg.Smtp.Password)
	o.Envelop = e.env
	if Cfg.Smtp.Hostport != 0 || Cfg.Smtp.Hostport != 587 {
		o.HostPort = Cfg.Smtp.Hostport
	}
	if err = o.Send(); err == nil {
		Info.Println("mail send via SMTP")
	}
	return
}

func (e *env) sendQueue() (err error) {
	var o = send.NewQueue()
	o.Env = e.env
	err = o.Send()
	if err == nil {
		Info.Println("mail send to Quene")
	}
	return
}

func (e *env) processQueue() {

	var dirNew = path.Join(DIR_QUEUE, "new")
	var dirSend = path.Join(DIR_QUEUE, "send")
	var dirFailed = path.Join(DIR_QUEUE, "failed")

	Info.Println("Create directories")
	if err := os.MkdirAll(dirNew, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirSend, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dirFailed, 0755); err != nil {
		log.Fatal(err)
	}

	Info.Println("Read: ", dirNew)
	files, err := os.ReadDir(dirNew)
	if err != nil {
		log.Fatalf(`Queue Read err: %s`, err.Error())
	}
	for _, file := range files {
		var nameGob = file.Name()
		if filepath.Ext(nameGob) == EXT_GOB {
			var id = nameGob[:len(nameGob)-len(EXT_GOB)]
			var nameEml = id + EXT_EML
			Info.Println("Processing id: ", id)

			var b []byte
			if b, err = os.ReadFile(path.Join(dirNew, nameGob)); err != nil {
				log.Fatalf(`Queue Read File err: %s`, err.Error())
			}
			e.env = envelop.New()
			if err = envelop.Unmarshal(b, e.env); err != nil {
				log.Fatalf(`Queue Unmarshal err: %s`, err.Error())
			}
			e.env.Create()
			if err = e.sendNow(); err != nil {
				Info.Println("Failed id: ", id)
				log.Println(err)
				if err = os.Rename(path.Join(dirNew, nameGob), path.Join(dirFailed, nameGob)); err != nil {
					log.Println(err)
				}
				if err = os.Rename(path.Join(dirNew, nameEml), path.Join(dirFailed, nameEml)); err != nil {
					log.Println(err)
				}
			} else {
				Info.Println("Moving to send folder: ", id)
				if err = os.Rename(path.Join(dirNew, nameGob), path.Join(dirSend, nameGob)); err != nil {
					log.Println(err)
				}
				if err = os.Rename(path.Join(dirNew, nameEml), path.Join(dirSend, nameEml)); err != nil {
					log.Println(err)
				}
			}
			e.env.Reset()
		}
	}
}

// Load gets your config from the json file,
// and fills your struct with the option
func LoadJson(filename string, o interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(b, &o)
	}
	return err
}

// Save will save your struct to the given filename,
// this is a good way to create a json template
func SaveJson(filename string, o interface{}) error {
	j, err := json.MarshalIndent(&o, "", "\t")
	if err == nil {
		err = ioutil.WriteFile(filename, j, 0660)
	}
	return err
}
