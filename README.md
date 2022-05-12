# pam-aad-oidc

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![Go version](https://img.shields.io/github/go-mod/go-version/alan-turing-institute/pam-aad-oidc.svg)](https://github.com/alan-turing-institute/pam-aad-oidc)
[![Go Report Card example](https://goreportcard.com/badge/github.com/alan-turing-institute/pam-aad-oidc)](https://goreportcard.com/report/github.com/alan-turing-institute/pam-aad-oidc)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

# Project summary

PAM module connecting to AzureAD for user authentication using OpenID Connect/OAuth2.

This code is based on code from [`pam-keycloak-oidc`](https://github.com/zhaow-de/pam-keycloak-oidc) and [`pam-ussh`](https://github.com/uber/pam-ussh).

# Installation

## Configure Azure Active Directory

1. Create a new `App Registration` in your Azure Active Directory.

   - Set the name to whatever you choose (in this example we will use `pam-aad-oidc`)
   - Set access to `Accounts in this organizational directory only`.
   - Set `Redirect URI` to `Public client/native (mobile & desktop)` with a value of `urn:ietf:wg:oauth:2.0:oob`

2. Under `Certificates & secrets` add a `New client secret`

   - Set the description to `Secret for PAM authentication`
   - Set the expiry time to whatever is relevant for your use-case
   - You must **record the value** of this secret at creation time, as it will not be visible later.

3. Under `API permissions`:
   - Ensure that `Microsoft Graph > User.Read` is enabled
   - Select this and click the `Grant admin consent` button (otherwise manual consent is needed from each user)

## Configure local client

1. Either download the latest precompiled binary from `https://github.com/alan-turing-institute/pam-aad-oidc/releases` or compile the code for your own machine.

2. Install the binary in `/lib/x86_64-linux-gnu/security/` or the equivalent for your system

3. Create a `TOML` configuration file in a sensible location (for example `/etc/pam-aad-oidc.toml`) with the following structure:

   ```toml
   # Tenant ID for this AzureAD
   tenant-id="07e4545b-d4e1-e60f-63ab-32a64c0e9346"

   # The Application (client) ID for your registered app
   client-id="0831d551-06ed-db79-d1f3-20a45f0279ae"

   # The (time-limited) client secret generated for this application above
   client-secret="jbi58~72en43pqpdvwg6enb8r0ml3-hq-0ip2s9c"

   # Microsoft.Graph scope to be requested. Unless there is a particular reason not to, use 'user.read'.
   scope="user.read"

   # Name of AAD group that authenticated users must belong to
   group-name="Allowed PAM users"

   # Default domain for AAD users. This will be appended to any users not in `username@domain` format.
   domain="mydomain.onmicrosoft.com"
   ```

4. Add configuration lines to the relevant PAM module, referencing the `TOML` file you wrote above.
   For example, for testing purposes you can add the following to `/etc/pam.d/test`

   ```none
   auth    required     pam_aad_oidc.so config=/etc/pam-aad-oidc.toml
   ```

5. Install `pamtester` in order to test the module.

   ```shell
   # With the password for `myusername` in the file `password.secret`
   > cat password.secret | pamtester test myusername authenticate
   ```

   You should see the message: `[myusername] Authentication succeeded`

# FAQ

## Does this handle MFA?

No. PAM only supports username and password, without the possibility of including a third factor.
The [`pam-keycloak-oidc`](https://github.com/zhaow-de/pam-keycloak-oidc) project includes support for TOTP where the OTP code is embedded into the username or password.
As AzureAD supports several kinds of MFA apart from TOTP we have chosen to leave MFA to other dedicated PAM modules.
**Note** This means that you must _not_ have AzureAD Conditional Access policies applying to this application which enforce the use of MFA.

## Why Go?

The original projects that this work was based off were both written in Go.
A compiled language is needed in order to produce shared libraries for use by PAM.
A high-level language is needed in order to use libraries for handling http requests and JWTs.

## Contributing

If you find this project useful but lacking in some respect, we hope you'll consider contributing back to it.

The easiest way to get involved is by [opening an issue](https://github.com/alan-turing-institute/pam-aad-oidc/issues/new) if you find a bug or have a request for a new feature.

If you'd like to help us tackle some of the technical challenges we follow a standard GitHub contribution process.
Please find or submit an issue and then submit a pull request (PR) that addresses it.
