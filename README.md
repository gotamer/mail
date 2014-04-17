gotamer/mail
============

Now gotamer/mail also implements the io.Writer interface.
 
```go
	package main

	import (
	    "fmt"
	    
	    "bitbucket.org/gotamer/mail"
	)

	func main() {
		s := new(mail.Smtp)
		s.SetHostname("smtp.gmail.com")
		s.SetHostport(587)
		s.SetFromName("GoTamer")
		s.SetFromAddr("xxxx@gmail.com")
		s.SetPassword("*********")
		s.SetToAddrs("one@example.com", "two@example.com")
		s.AddToAddr("three@example.com")
		s.SetSubject("GoTamer test smtp mail")
		s.SetBody("The Writer below should replace this default line")
		if _, err := fmt.Fprintf(s, "Hello, smtp mail writer\n"); err != nil {
			fmt.Println(err)
		}
	}
```

#### A note on the host:  
Go SMTP does not allow to connect to SMPT servers with a self signed certs, you will get an error like following:

	x509: certificate signed by unknown authority

The way I got around that is by using [CAcert][1]. [CAcert][1] provides FREE digital certificates.

### Links
 * [Pkg Documantation](http://go.pkgdoc.org/bitbucket.org/gotamer/mail "GoTamer Mail Pkg Documentation")
 * [Repository](https://bitbucket.org/gotamer/mail "GoTamer Mail Repository")


[1]: http://www.cacert.org  "CA Cert"
	


________________________________________________________

#### The MIT License (MIT)

Copyright © 2012-2013 Dennis T Kaplan <http://www.robotamer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
