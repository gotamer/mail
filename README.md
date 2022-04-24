gotamer/mail
=================
[![GoDoc Send](https://godoc.org/github.com/gotamer/mail?status.svg)](https://godoc.org/github.com/gotamer/mail)

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

Sending mail with sendmail over smtp
```bash
$ go get github.com/gotamer/mail/cmd/sendmail
# cd /usr/local/bin
# ln -s $GOBIN/sendmail
$ sudo runuser -u mail -- sendmail -t user@example.com -s "My subject" -b "Body of mail"
```

gotamer/mail/send also implements the io.Writer interface.  


gotamer/mail/imap
=================
[![GoDoc Imap](https://godoc.org/github.com/gotamer/mail/imap?status.svg)](https://godoc.org/github.com/gotamer/mail/imap)

gotamer/mail/imap to get mail from an imap server  

