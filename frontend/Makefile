fmt:
	./node_modules/.bin/prettier --write --trailing-comma all --no-semi --single-quote src/**/*.js

dist:
	rm -rf dist
	mkdir dist dist/static
	./node_modules/.bin/webpack --config webpack.config.prod.js -p
	cp build/* dist/static
	cp index.html dist
	rsync -a public/ dist

.PHONY: dist
