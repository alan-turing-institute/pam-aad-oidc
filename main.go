package main

import (
	"fmt"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#cgo CFLAGS: -Wall
#include <security/pam_appl.h>
#include <stdlib.h>
*/
import "C"

func main() {}

//export pam_sm_authenticate
func pam_sm_authenticate(handle *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	/*
		In this function we will ask the username and the password with pam_get_user() and pam_get_authtok().
		We will then decide if the user is authenticated
	*/
	username := GetUsername(handle)
	fmt.Println("Username is: " + username)

	password := GetPassword(handle)
	fmt.Println("Password is: " + password)

	return C.PAM_SUCCESS
}
