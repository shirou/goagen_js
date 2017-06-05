gen_js:
	cd example && goagen gen --pkg-path=github.com/shirou/goagen_js -d github.com/shirou/goagen_js/example/design

gen_js_flow:
	cd example && goagen gen --pkg-path=github.com/shirou/goagen_js -d github.com/shirou/goagen_js/example/design -- --target flow

gen_js_type:
	cd example && goagen gen --pkg-path=github.com/shirou/goagen_js -d github.com/shirou/goagen_js/example/design -- --target type --genout ts


build:
	rm -f example/*.go
	cd example && goagen app -d github.com/shirou/goagen_js/example/design
	cd example && goagen main -d github.com/shirou/goagen_js/example/design
	cd example && go build
