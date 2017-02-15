# go-wikiclient

This is a Go interface to the [wikifier](https://github.com/cooper/wikifier)
server. It is abstracted so that various types of transports can be added to
communicate with the wikiserver. Currently though only UNIX sockets are
supported.

## Usage

```go
package main

import (
    wikiclient "github.com/cooper/go-wikiclient"
    "log"
    "time"
)

var tr wikiclient.Transport
var sess *wikiclient.Session

func main() {

    // initialize the transport
    tr = wikiclient.NewUnixTransport("/path/to/wikiserver.sock")
    if err := tr.Connect(); err != nil {
        log.Fatal(err)
    }

    // create a session
    sess = &wikiclient.Session{WikiName: "mywiki", WikiPassword: "secret"}
}

func someHTTPHandlerProbably() {
    // create a client, which pairs the transport and session and
    // provides the high-level methods
    c := wikiclient.NewClient(tr, sess, 3 * time.Second)
    c.DisplayPage("some_page")
}
```

## See also

* [__wikifier__](https://github.com/cooper/wikifier)
* [__quiki__](https://github.com/cooper/quiki) - a standalone webserver for
  wikifier built atop go-wikiclient
