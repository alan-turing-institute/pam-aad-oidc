MODULE = pam_aad_oidc
PREFIX = /usr/local

default: build

build:
	go build -buildmode=c-shared -o ${MODULE}.so

install:
	install -D -m 644 ${MODULE}.so ${DESTDIR}${PREFIX}/lib/x86_64-linux-gnu/security/${MODULE}.so

test:
	go test -cover

clean:
	go clean
	rm -f ${MODULE}.so ${MODULE}.h

.PHONY: build install test clean
