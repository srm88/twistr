.PHONY: gofmt clean

ts: main.go twistr/*.go
	go build -o ts .

gofmt:
	@for f in `ls twistr`; do \
		gofmt twistr/$$f >twistr/$$f.bak && mv twistr/$$f.bak twistr/$$f; \
	done
	@gofmt main.go >main.go.bak && mv main.go.bak main.go

clean:
	@rm -f twistr/*.bak ts
