PREFIX=/usr/local

.PHONY: all fmt check vet clean install sloc

proxyguard: cmd/proxyguard/main.go
	go build -o $@ codeberg.org/eduVPN/proxyguard/cmd/proxyguard/...

all: proxyguard

fmt:
	gofmt -s -w cmd/proxyguard/*.go

check:
	# https://staticcheck.io/
	staticcheck codeberg.org/eduVPN/proxyguard/cmd/proxyguard/...

vet:
	go vet codeberg.org/eduVPN/proxyguard/cmd/proxyguard/...

sloc:
	cloc cmd/proxyguard/*.go

clean:
	rm -f proxyguard

install: all
	install -m 0755 -D proxyguard $(DESTDIR)$(PREFIX)/sbin/proxyguard
