/*
 * Copyright (c) $2013, Ömer Yildiz. All rights reserved.
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301  USA
 */
package shout

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lshout
#include <stdlib.h>
#include <shout/shout.h>
*/
import "C"

const (
	BUFFER_SIZE = 8192
)

const (
	// See shout.h
	SHOUTERR_SUCCESS     = 0
	SHOUTERR_INSANE      = -1
	SHOUTERR_NOCORRECT   = -2
	SHOUTERR_NOLOGIN     = -3
	SHOUTERR_SOCKET      = -4
	SHOUTERR_MALLOC      = -5
	SHOUTERR_METADATA    = -6
	SHOUTERR_CONNECTED   = -7
	SHOUTERR_UNCONNECTED = -8
	SHOUTERR_UNSUPPORTED = -9
	SHOUTERR_BUSY        = -10
)

const (
	FORMAT_OGG  = 0
	FORMAT_MP3  = 1
	FORMAT_WEBM = 2
)

const (
	PROTOCOL_HTTP       = 0
	PROTOCOL_XAUDIOCAST = 1
	PROTOCOL_ICY        = 2
)

type ShoutError struct {
	Message string
	Code    int
}

func (e ShoutError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

type Shout struct {
	Host     string
	Port     uint
	User     string
	Password string
	Mount    string
	Format   int
	Protocol int

	// wrap the native C struct
	struc *C.struct_shout

	stream chan []byte
}

func init() {
	C.shout_init()
}

func Shutdown() {
	C.shout_shutdown()
}

func Free(s *Shout) {
	C.shout_free(s.struc)
}

func (s *Shout) lazyInit() {
	if s.struc != nil {
		return
	}

	s.struc = C.shout_new()
	s.updateParameters()

	s.stream = make(chan []byte)
}

func (s *Shout) updateParameters() {
	// set hostname
	p := C.CString(s.Host)
	C.shout_set_host(s.struc, p)
	C.free(unsafe.Pointer(p))

	// set port
	C.shout_set_port(s.struc, C.ushort(s.Port))

	// set username
	p = C.CString(s.User)
	C.shout_set_user(s.struc, p)
	C.free(unsafe.Pointer(p))

	// set password
	p = C.CString(s.Password)
	C.shout_set_password(s.struc, p)
	C.free(unsafe.Pointer(p))

	// set mount point
	p = C.CString(s.Mount)
	C.shout_set_mount(s.struc, p)
	C.free(unsafe.Pointer(p))

	// set format
	C.shout_set_format(s.struc, C.uint(s.Format))

	// set protocol
	C.shout_set_protocol(s.struc, C.uint(s.Protocol))
}

func (s *Shout) GetError() string {
	s.lazyInit()
	err := C.shout_get_error(s.struc)
	return C.GoString(err)
}

func (s *Shout) Open() (chan<- []byte, error) {
	s.lazyInit()

	errcode := int(C.shout_open(s.struc))
	if errcode != C.SHOUTERR_SUCCESS {
		return nil, ShoutError{
			Code:    errcode,
			Message: s.GetError(),
		}
	}

	go s.handleStream()

	return s.stream, nil
}

func (s *Shout) Close() error {
	errcode := int(C.shout_close(s.struc))
	if errcode != C.SHOUTERR_SUCCESS {
		return ShoutError{
			Code:    errcode,
			Message: s.GetError(),
		}
	}

	return nil
}

func (s *Shout) send(buffer []byte) error {
	ptr := (*C.uchar)(&buffer[0])
	C.shout_send(s.struc, ptr, C.size_t(len(buffer)))

	errno := int(C.shout_get_errno(s.struc))
	if errno != C.SHOUTERR_SUCCESS {
		fmt.Println("something went wrong: %d", errno)
	}

	C.shout_sync(s.struc)
	return nil
}

func (s *Shout) handleStream() {
	for buf := range s.stream {
		s.send(buf)
	}
}
