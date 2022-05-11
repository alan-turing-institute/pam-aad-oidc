#ifndef PAM_AAD_OIDC__PAM_H
#define PAM_AAD_OIDC__PAM_H

// Retrieve a username from a PAM handle
char* get_username(pam_handle_t* handle);

// Retrieve a password from a PAM handle
char* get_password(pam_handle_t* handle);

#endif
