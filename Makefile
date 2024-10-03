build:
	mkdir -p site
	mkdir -p site/static
	cp static/* site/static/
	go run *.go

deploy:
	sudo cp -rv site/* /usr/share/nginx/html/
