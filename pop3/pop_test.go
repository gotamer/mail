package pop3_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gotamer/pop3"
	"github.com/joho/godotenv"
)

var c *pop3.Client

var (
	HostName string
	Username string
	Password string
)

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
		os.Exit(1)
	}
	HostName = os.Getenv("GOTAMER_POP3_HOSTNAME")
	Username = os.Getenv("GOTAMER_POP3_USERNAME")
	Password = os.Getenv("GOTAMER_POP3_USERPASS")
}

func TestConnect(t *testing.T) {
	if HostName == "" {
		t.Log("Set env vars! ", HostName)
		t.Fail()
		return
	}
	var err error
	if c, err = pop3.DialTLS(HostName); err != nil {
		t.Log(err)
	}
	// Authenticate with the server
	if err = c.Authorization(Username, Password); err != nil {
		t.Log(err)
	}

	if err = c.ListCapabilities(); err != nil {
		t.Log(err)
	}

	if err = c.Stat(); err != nil {
		t.Log(err)
	}

	if err = c.ListAll(); err != nil {
		t.Log(err)
	}

	fmt.Println("Capa: ", c.CapaCAPA)
	fmt.Println("Count: ", c.Count)

	for _, v := range c.List {
		fmt.Printf("id: %d UID: %s Size: %d\n", v.ID, v.UID, v.Size)
	}

	if err = c.Quit(); err != nil {
		t.Log(err)
	}
}
