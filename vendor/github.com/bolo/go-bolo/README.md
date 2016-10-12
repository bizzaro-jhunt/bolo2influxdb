# Overview

*go-bolo* is a golang library designed to make it easier to stream and parse
[bolo](https://github.com/bolo/bolo) messages from inside Go applications.

# Using

TO make use of this library, perform the usual `go get github.com/bolo/go-bolo`,
and import it into your project and start connecting!:

```
package main

import "github.com/bolo/go-bolo"

func main() {
	// Stream bolo data
	pduChan, errChan, err := bolo.Connect("tcp://10.10.10.10:2997")
}
```

Bolo messages will be streamed onto `pduChan` as they come in, and any
errors will come over the `errChan`. If there was a problem setting up
the connection, it is returned via `err`.


# Docs

https://godoc.org/github/bolo/go-bolo
