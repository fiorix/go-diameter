all: helloworld server.crt

deps:
	go get -u github.com/fiorix/go-diameter/base

helloworld:
	go build -v -o helloworld

server.crt:
	@echo KEY
	openssl genrsa -des3 -out server.key 4096
	@echo
	@echo CSR: Set \"Common Name\" to the hostname. e.g. localhost
	openssl req -new -key server.key -out server.csr
	@echo
	@echo Remove passphrase
	cp server.key orig.server.key
	openssl rsa -in orig.server.key -out server.key
	@echo
	openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
	chmod 600 orig.server.key server.key server.csr server.crt

clean:
	rm -f helloworld orig.server.key server.key server.csr server.crt

.PHONY: helloworld
