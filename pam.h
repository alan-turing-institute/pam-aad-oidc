#ifndef PAM_AAD_OIDC__PAM_H
#define PAM_AAD_OIDC__PAM_H

// Retrieve a username from a PAM handle
char* get_username(pam_handle_t* handle);

// Retrieve a password from a PAM handle
char* get_password(pam_handle_t* handle);

// Get the specified item from an array
char *get_array_item(int i, char **argv);

#endif
