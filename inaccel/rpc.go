package inaccel

/*
#cgo LDFLAGS: -lcoral-api

#include <errno.h>
#include <inaccel/coral.h>
#include <string.h>

static int get_errno() {
	return errno;
}
*/
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

func Submit(request *Request) chan error {
	cresponse := C.inaccel_response_create()
	if cresponse == nil {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}

	if error := C.inaccel_submit(request.c, cresponse); error != 0 {
		errsv := C.get_errno()

		C.inaccel_response_release(cresponse)

		panic(fmt.Errorf(C.GoString(C.strerror(errsv))))
	}

	response := make(chan error)

	go func() {
		defer close(response)

		response <- func() error {
			if error := C.inaccel_response_wait(cresponse); error == -1 {
				errsv := C.get_errno()

				C.inaccel_response_release(cresponse)

				panic(fmt.Errorf(C.GoString(C.strerror(errsv))))
			} else if error != 0 {
				n := C.inaccel_response_snprint(nil, 0, cresponse)
				if n < 0 {
					errsv := C.get_errno()

					C.inaccel_response_release(cresponse)

					panic(fmt.Errorf(C.GoString(C.strerror(errsv))))
				}

				s := make([]byte, n+1)
				if C.inaccel_response_snprint((*C.char)(reflect.ValueOf(s).UnsafePointer()), C.size_t(len(s)), cresponse) != n {
					errsv := C.get_errno()

					C.inaccel_response_release(cresponse)

					panic(fmt.Errorf(C.GoString(C.strerror(errsv))))
				}

				C.inaccel_response_release(cresponse)

				return fmt.Errorf(string(s))
			}

			C.inaccel_response_release(cresponse)

			return nil
		}()
	}()

	return response
}

type Request struct {
	c     C.inaccel_request
	index uint32
}

func NewRequest(accelerator string) *Request {
	request := &Request{
		C.inaccel_request_create((*C.char)(reflect.ValueOf([]byte(accelerator)).UnsafePointer())),
		0,
	}
	if request.c == nil {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}
	runtime.SetFinalizer(request, func(request *Request) {
		C.inaccel_request_release(request.c)
	})
	return request
}

func (request *Request) ArgArray(size uint64, value unsafe.Pointer, index ...uint32) *Request {
	var _index uint32
	if index == nil {
		_index = request.index
	} else {
		_index = index[0]
	}

	if error := C.inaccel_request_arg_array(request.c, C.size_t(size), value, C.uint(_index)); error != 0 {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}

	if index == nil {
		request.index++
	}

	return request
}

func (request *Request) ArgScalar(size uint64, value unsafe.Pointer, index ...uint32) *Request {
	var _index uint32
	if index == nil {
		_index = request.index
	} else {
		_index = index[0]
	}

	if error := C.inaccel_request_arg_scalar(request.c, C.size_t(size), value, C.uint(_index)); error != 0 {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}

	if index == nil {
		request.index++
	}

	return request
}

func (request Request) String() string {
	n := C.inaccel_request_snprint(nil, 0, request.c)
	if n < 0 {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}

	s := make([]byte, n+1)
	if C.inaccel_request_snprint((*C.char)(reflect.ValueOf(s).UnsafePointer()), C.size_t(len(s)), request.c) != n {
		panic(fmt.Errorf(C.GoString(C.strerror(C.get_errno()))))
	}

	return string(s)
}
