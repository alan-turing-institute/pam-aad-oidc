// Standard library includes
#include <stdlib.h>
#include <string.h>

// Third party includes
#include <security/pam_appl.h>
#include <security/pam_ext.h>


// Retrieve a username from a PAM handle
char* get_user(pam_handle_t* handle) {
    if (!handle) {
        return NULL;
    }
    const char* user = NULL;
    // pam_get_item outputs into `const void**` so we must cast `user` accordingly
    if (pam_get_item(handle, PAM_USER, (const void**)&user) != PAM_SUCCESS) {
        return NULL;
    }
    return strdup(user);
}


// Retrieve a password from a PAM handle
char* get_password(pam_handle_t* handle) {
    if (!handle) {
        return NULL;
    }
    const char* password = NULL;
    if (pam_get_authtok(handle, PAM_AUTHTOK, &password, NULL) != PAM_SUCCESS) {
        return NULL;
    }
    return strdup(password);
}
