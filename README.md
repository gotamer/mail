gotamer/mail/send
=================
[![GoDoc Send](https://godoc.org/github.com/gotamer/mail/send?status.svg)](https://godoc.org/github.com/gotamer/mail/send)

Sending mail with sendmail or smtp

gotamer/mail/send implements the io.Writer interface.  

```bash
$ sendmail -h
  -b string
        Mail body text
  -c string
        Config file location (default "/etc/sendmail.json")
  -d    Debug mode
  -f string
        Mail From email address
  -h    Display help
  -now
        Send message now, message will not be added to Queue
  -run
        Process Queue
  -s string
        Mail subject line
  -t string
        Mail To email address
``` 
gotamer/mail/imap
=================
[![GoDoc Imap](https://godoc.org/github.com/gotamer/mail/imap?status.svg)](https://godoc.org/github.com/gotamer/mail/imap)

gotamer/mail/imap to get mail from an imap server  

