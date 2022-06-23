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
	"unsafe"
)

func Alloc(size uint64) (unsafe.Pointer, error) {
	ptr := C.inaccel_alloc(C.size_t(size))
	if ptr == nil {
		return nil, fmt.Errorf(C.GoString(C.strerror(C.get_errno())))
	}
	return ptr, nil
}

func Free(ptr unsafe.Pointer) {
	C.inaccel_free(ptr)
}

func Realloc(ptr unsafe.Pointer, size uint64) (unsafe.Pointer, error) {
	ptr = C.inaccel_realloc(ptr, C.size_t(size))
	if ptr == nil {
		return nil, fmt.Errorf(C.GoString(C.strerror(C.get_errno())))
	}
	return ptr, nil
}
