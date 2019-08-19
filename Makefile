.PHONY: build 

browser: browser_za

browser_za:
	go build -v -o build/block-browser
clean:
	rm -rf build
test:
	go test ./job
