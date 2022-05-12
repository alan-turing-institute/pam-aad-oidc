package main

import (
	"strings"
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

	// Load username and password from the pam_handle
	username := GetUsername(handle)
	password := GetPassword(handle)

	// Load config path from the config=XXX module argument
	configPath := ""
	for _, option := range ToArray(argc, argv) {
		tokens := strings.Split(option, "=")
		if tokens[0] == "config" {
			configPath = tokens[1]
		}
	}

	// ValidateCredentials returns an int corresponding to one of pam_types.h
	return C.int(ValidateCredentials(configPath, username, password))
}
