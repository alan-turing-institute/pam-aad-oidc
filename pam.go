package main

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#cgo CFLAGS: -Wall
#include <security/pam_appl.h>
#include <stdlib.h>
#include "pam.h"
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

func ToArray(argc C.int, argv **C.char) []string {
	array := make([]string, 0, argc)
	for i := 0; i < int(argc); i++ {
		cString := C.get_array_item(C.int(i), argv)
		defer C.free(unsafe.Pointer(cString))
		array = append(array, C.GoString(cString))
	}
	return array
}
