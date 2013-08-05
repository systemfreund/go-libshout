package shout 

import "net/url"

/*
#cgo LDFLAGS: -lshout
#include <shout/shout.h>
*/
import "C"

type Shout struct {
	// hostname or IP of icecast server 
	Host url.URL

	// port of the icecast server
	Port int

	// login password for the server
	Password string

	// server protocol to use
	Protocol int

	// type of data being sent
	Format int

	// audio encoding parameters
	// TODO: util_dict *audio_info

	// user-agent to use when doing HTTP login
	Useragent string

	// mountpoint for this stream
	Mount string

	// name of the stream
	Name string

	// homepage of the stream
	Url string

	// genre of the stream
	Genre string

	// description of the stream
	Description string

	// username to use for HTTP auth
	User string

	// is this stream private?
	Public bool


}

func Init() {
	C.shout_init();
}

func Shutdown() {
	C.shout_shutdown();
}

func New() *C.struct_shout {
	return C.shout_new()
}