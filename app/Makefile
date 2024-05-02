TEMPL_FILES = $(shell find templates/ -type f -name '*.templ' -print)
OUT=./assets/dist/styles.css

dev: ./assets/dist/styles.css templ
	air --build.cmd "go build -o bin/app ./cmd/**" --build.bin "./bin/app"

strip:
	strip ./bin/app
	upx ./bin/app

./assets/dist/styles.css: $(TEMPL_FILES)
	npx tailwindcss -i ./assets/tailwind.css -o ./assets/dist/css/styles.css

templ: $(TEMPL_FILES)
	templ generate templates

templDev: $(TEMPL_FILES)
	templ generate -watch templates

tw:
	npx tailwindcss -i ./assets/tailwind.css -o ./assets/dist/css/styles.css

twDev:
	npx tailwindcss -i ./assets/tailwind.css -o ./assets/dist/css/styles.css -w

clean:
	@find ./assets/dist -type f -name "*.gz" | xargs -I {} rm -f {}
	@rm -f ./assets/dist/css/styles.css
	@rm -f ./bin/app
