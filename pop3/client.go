package pop3

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"strconv"
	"strings"
)

// Debug mode will add the Client Info and Error information
var Debug bool

type MessageList struct {
	// Non unique id reported by the server
	ID int

	// Unique id reported by the server
	UID string

	// Size of the message
	Size int
}

// Client for POP3 with message list, which represents the metadata returned by the server for a
// messages stored in the maildrop.
type Client struct {
	conn *Connection

	// server support for POP3 CAPA command
	CapaCAPA bool

	// server support for POP3 TOP command
	CapaTOP bool

	// server support for POP3 UIDL command
	CapaUIDL bool

	// Count of messages
	Count int

	// Size of all messages
	Size int

	List []MessageList

	// List of info messages
	Info []string

	// List of errors
	Errors []error
}

// Dial opens new connection and creates a new POP3 client.
func Dial(addr string) (c *Client, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		err = fmt.Errorf("failed to dial: %w", err)
		return
	}
	return NewClient(conn)
}

// DialTLS opens new TLS connection and creates a new POP3 client.
func DialTLS(addr string) (c *Client, err error) {
	var conn *tls.Conn
	if conn, err = tls.Dial("tcp", addr, nil); err != nil {
		err = fmt.Errorf("failed to dial TLS: %w", err)
		return
	}
	return NewClient(conn)
}

// NewClient creates a new POP3 client.
func NewClient(conn net.Conn) (c *Client, err error) {
	c = &Client{
		conn: NewConnection(conn),
	}
	c.addInfo("Connected to POP3 server")
	var line string
	// Make sure we receive the server greeting
	if line, err = c.conn.ReadLine(); err != nil {
		err = fmt.Errorf("failed to read line: %w", err)
		c.addError(err)
		return
	}

	c.addInfo(fmt.Sprintf("Greeting: %s", line))

	if strings.Fields(line)[0] != "+OK" {
		err = fmt.Errorf("server did not response with +OK: %s", line)
		c.addError(err)
		return nil, err
	}
	return c, nil
}

// Authorization logs into POP3 server with login and password.
func (c *Client) Authorization(user, pass string) (err error) {
	var s string
	if s, err = c.conn.Cmd("%s %s", "USER", user); err != nil {
		err = fmt.Errorf("failed at USER command: %w", err)
		c.addError(err)
		return
	}
	c.addInfo(fmt.Sprintf("User: %s", s))

	if s, err = c.conn.Cmd("%s %s", "PASS", pass); err != nil {
		err = fmt.Errorf("failed at PASS command: %w", err)
		c.addError(err)
		return
	}
	c.addInfo(fmt.Sprintf("Password: %s", s))
	return c.Noop()
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (c *Client) Quit() (err error) {
	var s string
	if s, err = c.conn.Cmd("QUIT"); err != nil {
		err = fmt.Errorf("failed at QUIT: %w", err)
		c.addError(err)
		return
	}
	c.addInfo(fmt.Sprintf("QUIT cmd: %s", s))

	if err = c.conn.Close(); err != nil {
		err = fmt.Errorf("failed to close connection: %w", err)
		c.addError(err)
	} else {
		c.addInfo("connection closed")
	}
	return
}

// Status checks if we are connected to a pop3 server
func (c *Client) Status() bool {
	return c.conn != nil
}

// Noop will do nothing however can prolong the end of a connection.
func (c *Client) Noop() (err error) {
	if _, err = c.conn.Cmd("NOOP"); err != nil {
		err = fmt.Errorf("failed at NOOP: %w", err)
		c.addError(err)
	}
	return
}

// CAPA List Capabilities.
func (c *Client) ListCapabilities() (err error) {
	if _, err = c.conn.Cmd("CAPA"); err != nil {
		c.addError(err)
		return
	}
	lines, err := c.conn.ReadLines()
	if err != nil {
		c.addError(err)
		err = fmt.Errorf("failed at CAPA command: %w", err)
		return
	}
	for i, v := range lines {
		c.addInfo(fmt.Sprintf("CAPA Line %d: %s", i, v))

		switch v {
		case "CAPA":
			c.CapaCAPA = true
		case "TOP":
			c.CapaTOP = true
		case "UIDL":
			c.CapaUIDL = true
		}
	}
	return
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// In the event of an error, all returned numeric values will be 0.
func (c *Client) Stat() (err error) {
	line, err := c.conn.Cmd("STAT")
	if err != nil {
		c.addError(err)
		return
	}

	if len(strings.Fields(line)) != 3 {
		err = fmt.Errorf("invalid response returned from server: %s", line)
		c.addError(err)
		return
	}

	// Number of messages in maildrop
	c.Count, err = strconv.Atoi(strings.Fields(line)[1])
	if err != nil {
		c.addError(err)
		return
	}
	if c.Count == 0 {
		c.addInfo("STAT no messages found")
		return
	}

	// Total size of messages in bytes
	if c.Size, err = strconv.Atoi(strings.Fields(line)[2]); err != nil {
		c.addError(err)
	}
	c.addInfo(fmt.Sprintf("STAT number of messages: %d", c.Count))
	c.addInfo(fmt.Sprintf("STAT messages total bytes: %d", c.Size))
	return
}

// ListAll returns a MessageList object which contains all messages in the maildrop.
func (c *Client) ListAll() (err error) {

	if _, err = c.conn.Cmd("LIST"); err != nil {
		c.addError(err)
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		c.addError(err)
		return
	}

	var NoUIDL bool
	var uids []string
	if _, err = c.conn.Cmd("UIDL"); err != nil {
		c.addError(err)
	} else {
		if uids, err = c.conn.ReadLines(); err != nil {
			c.addError(err)
		}
	}
	if len(uids) != len(lines) {
		NoUIDL = true
		c.addInfo("UIDL not available")
	}

	for i, v := range lines {
		var id int
		id, err = strconv.Atoi(strings.Fields(v)[0])
		if err != nil {
			c.addError(err)
			return
		}

		var uid string
		if NoUIDL {
			uid = ""
		} else {
			uid = strings.Fields(uids[i])[1]
		}

		var size int
		if size, err = strconv.Atoi(strings.Fields(v)[1]); err != nil {
			c.addError(err)
		}

		c.List = append(c.List, MessageList{id, uid, size})
		//fmt.Println("id: ", id, " UID: ", uid, " Size: ", size)
	}
	return
}

// TOP Return headers.
func (c *Client) Top() (headers []string, err error) {
	if _, err = c.conn.Cmd("TOP 1 0"); err != nil {
		err = fmt.Errorf("cmd error from server: %w", err)
		c.addError(err)
	} else if headers, err = c.conn.ReadLines(); err != nil {
		err = fmt.Errorf("TOP ReadLine: %w", err)
		c.addError(err)
	}
	return
}

// Rset will unmark any messages that have being marked for deletion in
// the current session.
func (c *Client) Rset() (err error) {
	if _, err := c.conn.Cmd("RSET"); err != nil {
		err = fmt.Errorf("failed at RSET: %w", err)
		c.addError(err)
	}
	return
}

// RetrRaw downloads the given message and returns it as []byte object.
func (c *Client) RetrRaw(msg int) (message []byte, err error) {
	if _, err = c.conn.Cmd("RETR %d", msg); err != nil {
		err = fmt.Errorf("failed at RETR command: %w", err)
		c.addError(err)
		return
	}
	message, err = c.conn.ReadDot()
	if err != nil {
		err = fmt.Errorf("failed to read message: %w", err)
		c.addError(err)
	}
	return
}

// Retr downloads the given message and returns it as a mail.Message object.
func (c *Client) Retr(msg int) (message *mail.Message, err error) {
	if _, err = c.conn.Cmd("RETR %d", msg); err != nil {
		err = fmt.Errorf("failed at RETR command: %w", err)
		c.addError(err)
		return
	}
	message, err = mail.ReadMessage(c.conn.Reader.DotReader())
	if err != nil {
		err = fmt.Errorf("failed to read message: %w", err)
		c.addError(err)
	}
	return
}

// Dele will delete the given message from the maildrop.
// Changes will only take affect after the Quit command is issued.
func (c *Client) Dele(msg int) (err error) {
	if _, err := c.conn.Cmd("%s %d", "DELE", msg); err != nil {
		err = fmt.Errorf("failed at DELE: %w", err)
		c.addError(err)
	}
	return
}

func (c *Client) addError(err error) {
	if Debug {
		c.Errors = append(c.Errors, err)
	}
}
func (c *Client) addInfo(info string) {
	if Debug {
		c.Info = append(c.Info, info)
	}
}
