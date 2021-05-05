build:
	go build ./...

lambda:
	go build -o users_create ./lambdas/users/create
	go build -o users_delete ./lambdas/users/delete
	go build -o users_get ./lambdas/users/get
	go build -o users_update ./lambdas/users/update
	
	zip snailmail_lambdas.zip users_create users_delete users_get users_update
	
	rm users_create users_delete users_get users_update

mocks:
	mockery --all --dir ./resources --output ./test/resources/mocks