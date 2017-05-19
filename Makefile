gen_js:
	cd example && goagen gen --pkg-path=github.com/shirou/goagen_js -d github.com/shirou/goagen_js/example/design

build:
	rm example/*.go
	cd example && goagen app -d github.com/shirou/goagen_js/example/design
	cd example && goagen main -d github.com/shirou/goagen_js/example/design
	cd example && go build
