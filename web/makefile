build:
	docker build -t snailmail-web .

run: build
	docker run -d --mount type=bind,source="$(shell pwd)"/snailmail-web,target=/usr/share/nginx/html/ -p 8000:90 snailmail-web