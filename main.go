package main

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#cgo CFLAGS: -Wall
#include <security/pam_appl.h>
#include <stdlib.h>

char *get_password(pam_handle_t *handle);
char *get_user(pam_handle_t *handle);
*/
import "C"

func main() {
	fmt.Println("Testing main function. To be removed.")
}

//export pam_sm_authenticate
func pam_sm_authenticate(handle *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	/*
		In this function we will ask the username and the password with pam_get_user() and pam_get_authtok().
		We will then decide if the user is authenticated
	*/

	// Retrieve username
	cUsername := C.get_user(handle)
	defer C.free(unsafe.Pointer(cUsername))

	cPassword := C.get_password(handle)
	defer C.free(unsafe.Pointer(cPassword))

	username := C.GoString(cUsername)
	fmt.Println("Username is: " + username)

	password := C.GoString(cPassword)
	fmt.Println("Password is: " + password)

	return C.PAM_SUCCESS
}
