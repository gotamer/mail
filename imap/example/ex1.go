package main

import (
	//"io/ioutil"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"bitbucket.org/gotamer/cfg"
	"bitbucket.org/gotamer/errors"
	"bitbucket.org/gotamer/mail/imap"
	"bitbucket.org/gotamer/mail/send"
)

const (
	APPNAME = "IMAP"
	PS      = string(os.PathSeparator)
)

var (
	LogFile *os.File

	help    = flag.Bool("h", false, "Display this help")
	cfgfile = flag.String("c", "", "Config file ("+APPNAME+".json)")
)

var Cfg struct {
	Debug  uint8
	LogDir string
	Mail   struct {
		FromName string
		FromAddr string
		LogTo    string // eMail To address
	}
	IMAP struct {
		Hostname string
		Username string
		Password string
		Folder   string
	}
}

func main() {
	setup()

	mb := imap.NewMailbox(Cfg.IMAP.Hostname, Cfg.IMAP.Username, Cfg.IMAP.Password, Cfg.IMAP.Folder)

	mb.Dial()
	defer mb.Client.Logout(30 * time.Second)

	mb.ID("GoTamer")
	mb.NoOp()
	mb.Login()
	mb.Select()

	set := mb.UnSeen()
	e.Debug("IMAP", "set: %v", set)

	if set.Empty() == false {
		mails := make(chan imap.Mail)
		go mb.Body(set, mails)

		for m := range mails {
			if m.Err == nil && m.Mail == nil {
				break
			}
			if m.Err != nil {
				fmt.Println("Mails err: ", m.Err)
				break
			}
			fmt.Println(m.Mail)
		}
	}
}

func setup() {
	flag.Parse()

	if *help == true || *cfgfile == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if err := cfg.Load(*cfgfile, &Cfg); err != nil {
		if err = cfg.Save(*cfgfile, &Cfg); err != nil {
			fmt.Println("cfg.Save error: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("Please edit your config file at:\n\n\t%s\n\n", *cfgfile)
			os.Exit(0)
		}
	}

	e.AppName(APPNAME)

	switch Cfg.Debug {
	case 0:
		e.ENVIRONMENT = e.ENV_PROD
		break
	case 1:
		e.ENVIRONMENT = e.ENV_BETA
		break
	default:
		e.ENVIRONMENT = e.ENV_DEBG
	}

	if e.ENVIRONMENT != e.ENV_DEBG {
		var file string
		if Cfg.LogDir != "" {
			file = fmt.Sprintf("%s%s%s.log", Cfg.LogDir, PS, APPNAME)
		} else {
			file = fmt.Sprintf("%s.log", APPNAME)
		}

		e.FileLogger(file)
		imap.DefaultLogger = log.New(e.LogFile, "[IMAP]", 19)

		// sendmail Logger
		sm := send.NewSendMail()
		sm.Mail.From = Cfg.Mail.FromAddr
		sm.Mail.Subject = fmt.Sprintf("ERROR %s", APPNAME)
		sm.Mail.Body = "Default Error Body"
		sm.Mail.SetTo(Cfg.Mail.LogTo)
		e.MailLogger(sm)
	}
	log.Println("======= INIT =======")
	e.Debug("Main", "Start Up Setup Done")
}
