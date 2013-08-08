libshout.go
===========

(Incomplete) Go binding for libshout 2.x


Example
=======

```go
package main

import (
	"flag"
	"github.com/systemfreund/libshout.go"
	"os"
	"io"
)

// Setup some command line flags
var (
	hostname = flag.String("host", "localhost", "shoutcast server name")
	port = flag.Uint("port", 8000, "shoutcast server source port")
	user = flag.String("user", "source", "source user name")
	password = flag.String("password", "", "source password")
	mount = flag.String("mountpoint", "/stream.mp3", "mountpoint")
	filename = flag.String("file", "", "file to send to shoutcast")
) 

func main() {
	flag.Parse()

	// Setup libshout parameters
	s := shout.Shout{
		Host:     *hostname,
		Port:     *port,
		User:     *user,
		Password: *password,
		Mount:    *mount,
		Format:   shout.FORMAT_MP3,
		Protocol: shout.PROTOCOL_HTTP,
	}

	// Open the file
	//
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	// Create a channel where we can send the data
	//
	stream, err := s.Open()
	if err != nil {
		panic(err)
	}

	defer s.Close()

	buffer := make([]byte, shout.BUFFER_SIZE)
	for {
		// Read from file
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF { panic(err)	}
		if n == 0 { break }

		// Send to shoutcast server
		stream <- buffer
	}

	// done
}
```
