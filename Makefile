PREFIX=/usr/local

.PHONY: all lint clean install sloc

proxyguard: cmd/proxyguard/main.go
	go build -o $@ codeberg.org/eduVPN/proxyguard/cmd/proxyguard/...

all: proxyguard

lint:
	golangci-lint run ./... -E stylecheck,revive,gocritic

sloc:
	cloc cmd/proxyguard/*.go

clean:
	rm -f proxyguard

install: all
	install -m 0755 -D proxyguard $(DESTDIR)$(PREFIX)/sbin/proxyguard
