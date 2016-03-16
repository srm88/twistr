.PHONY: gofmt

gofmt:
	@for f in `ls src/twistr`; do \
		gofmt src/twistr/$$f >src/twistr/$$f.bak && mv src/twistr/$$f.bak src/twistr/$$f; \
	done
	@gofmt main.go >main.go.bak && mv main.go.bak main.go
