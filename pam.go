package main

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#cgo CFLAGS: -Wall
#include <security/pam_appl.h>
#include <stdlib.h>
#include <pam.h>
*/
import "C"

func GetUsername(handle *C.pam_handle_t) string {
	// Retrieve username
	cUsername := C.get_username(handle)
	defer C.free(unsafe.Pointer(cUsername))
	return C.GoString(cUsername)
}

func GetPassword(handle *C.pam_handle_t) string {
	// Retrieve password
	cPassword := C.get_password(handle)
	defer C.free(unsafe.Pointer(cPassword))
	return C.GoString(cPassword)
}
