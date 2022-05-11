MODULE := pam_aad_oidc

default: module

module: test
	go build -buildmode=c-shared -o ${MODULE}.so

test: *.go
	go test -cover

clean:
	go clean
	-rm -f ${MODULE}.so ${MODULE}.h

.PHONY: test module clean