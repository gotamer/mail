// High level IMAP access library
package imap

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"bitbucket.org/gotamer/errors"
	"bitbucket.org/gotamer/imap"
)

const (
	PS = string(os.PathSeparator)
)

var (
	DefaultLogger  = imap.DefaultLogger
	DefaultLogMask = imap.DefaultLogMask
)

type mailbox struct {
	Hostname string
	Username string
	Password string
	Folder   string
	Client   *imap.Client
}

// Mail wraps the email responses and possible errors.
type Mail struct {
	Mail *mail.Message
	Err  error
}

func init() {
	imap.DefaultLogMask = imap.LogConn | imap.LogRaw
}

// Set the mailbox settings
func NewMailbox(host, user, pass, folder string) *mailbox {
	mb := new(mailbox)
	mb.Hostname = host
	mb.Username = user
	mb.Password = pass
	mb.Folder = folder
	return mb
}

// Connect to the imap server
func (mb *mailbox) Dial() {
	var err error
	if strings.HasSuffix(mb.Hostname, ":993") {
		mb.Client, err = imap.DialTLS(mb.Hostname, nil)
	} else {
		mb.Client, err = imap.Dial(mb.Hostname)
	}
	e.Fail(err)
}

// IMAP NOOP
func (mb *mailbox) NoOp() {
	ReportOK(mb.Client.Noop())
}

// Sets the Application ID
func (mb *mailbox) ID(id string) {
	if mb.Client.Caps["ID"] {
		ReportOK(mb.Client.ID("name", id))
	}
}

// @todo set State error
func (mb *mailbox) Login() {
	if mb.Client.State() == imap.Login {
		defer mb.Client.SetLogMask(mb.Sensitive("LOGIN"))
		ReportOK(mb.Client.Login(mb.Username, mb.Password))
	}
}

func (mb *mailbox) Sensitive(action string) imap.LogMask {
	mask := mb.Client.SetLogMask(imap.LogConn)
	hide := imap.LogCmd | imap.LogRaw
	if mask&hide != 0 {
		mb.Client.Logln(imap.LogConn, "Raw logging disabled during", action)
	}
	mb.Client.SetLogMask(mask &^ hide)
	return mask
}

// Selects the mailbox
func (mb *mailbox) Select() {
	ReportOK(mb.Client.Select(mb.Folder, true))
}

// Gets the message UID for not seen messages
func (mb *mailbox) UnSeen() *imap.SeqSet {

	set, _ := imap.NewSeqSet("")
	unseenset, _ := imap.NewSeqSet("")

	if mb.Client.Mailbox == nil {
		return unseenset
	}

	set.AddRange(1, mb.Client.Mailbox.Messages)
	cmd, err := mb.Client.Fetch(set, "UID", "FLAGS")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for cmd.InProgress() {
		// Wait for the next response (timeout)
		err := mb.Client.Recv(time.Second + 3)
		if err != nil {
			fmt.Errorf("client.Recv err: %s", err)
		}

		// Process command data
		for _, rsp := range cmd.Data {

			uid := rsp.MessageInfo().UID

			flags := imap.AsFlagSet(rsp.MessageInfo().Attrs["FLAGS"])
			if flags["\\Seen"] == false {
				unseenset.AddNum(uid)
			}
		}
		cmd.Data = nil

		// Process unilateral server data
		for _, rsp := range mb.Client.Data {
			fmt.Println("Server data:", rsp)
		}
		mb.Client.Data = nil
	}
	return unseenset
}

// Channels the Body of selected messages
func (mb *mailbox) Body(set *imap.SeqSet, mails chan Mail) {

	// nothing to request!
	if set.Empty() {
		mails <- Mail{}
		return
	}

	// "UID", "FLAGS", "INTERNALDATE", "RFC822.SIZE", "BODY[]"
	cmd, err := mb.Client.UIDFetch(set, "UID", "BODY[]")
	if err != nil {
		mails <- Mail{Err: err}
		return
	}

	for cmd.InProgress() {
		// Wait for the next response (timeout)
		err := mb.Client.Recv(time.Second + 3)
		if err != nil {
			mails <- Mail{Err: err}
			return
		}

		// Process command data
		for _, rsp := range cmd.Data {
			var email *mail.Message
			body := imap.AsBytes(rsp.MessageInfo().Attrs["BODY[]"])
			if err == nil {
				email, err = mail.ReadMessage(bytes.NewReader(body))
			}
			mails <- Mail{email, err}
		}
		cmd.Data = nil
	}
	mails <- Mail{}
}

//####### IMAP Error #############
type errorIMAP struct {
	cmd  *imap.Command
	text string
}

// Returns a IMAP Error.
func NewImapErr(cmd *imap.Command, err error) error {
	return &errorIMAP{cmd, err.Error()}
}

// Returns a IMAP error as a formated string
func (e *errorIMAP) Error() string {
	name := "Unknown"
	if e.cmd != nil {
		name = e.cmd.Name(true)
	}
	return fmt.Sprintf("\n--- %s ---\n%v\n\n", name, e.text)
}

func ReportOK(cmd *imap.Command, err error) *imap.Command {
	var ierr error
	var rsp *imap.Response
	if cmd == nil {
		ierr = NewImapErr(cmd, err)
	} else if err == nil {
		rsp, err = cmd.Result(imap.OK)
	}
	if err != nil {
		ierr = NewImapErr(cmd, err)
	}
	e.Check(ierr)

	if e.ENVIRONMENT == e.ENV_BETA {
		c := cmd.Client()
		fmt.Printf("--- %s ---\n"+
			"%d command response(s), %d unilateral response(s)\n"+
			"%s %s\n\n",
			cmd.Name(true), len(cmd.Data), len(c.Data), rsp.Status, rsp.Info)
		c.Data = nil
	}
	return cmd
}
