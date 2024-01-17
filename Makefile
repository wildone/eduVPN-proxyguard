PREFIX=/usr/local

.PHONY: all client server lint clean install sloc

all: client server

client: cmd/proxyguard-client/main.go
	go build -o proxyguard-$@ codeberg.org/eduVPN/proxyguard/cmd/proxyguard-client/...

server: cmd/proxyguard-server/main.go
	go build -o proxyguard-$@ codeberg.org/eduVPN/proxyguard/cmd/proxyguard-server/...

lint:
	golangci-lint run ./... -E stylecheck,revive,gocritic

sloc:
	cloc --include-ext=go .

clean:
	rm -f proxyguard

install-client: client
	install -m 0755 -D proxyguard-client $(DESTDIR)$(PREFIX)/sbin/proxyguard-client

install-server: server
	install -m 0755 -D proxyguard-server $(DESTDIR)$(PREFIX)/sbin/proxyguard-server

install: install-client install-server
