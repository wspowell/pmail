build:
	GOOS=linux GOGC=3000 go build .

test: build
	golangci-lint run
	go test ./...

mocks:
	mockery --all --dir ./resources --output ./test/resources/mocks

lambda:
	rm -r ./bin/lambdas

	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/ws/connect/main ./lambdas/ws/connect
	zip -j ./bin/lambdas/ws/connect/main.zip ./bin/lambdas/ws/connect/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-ws-connect --zip-file fileb://./bin/lambdas/ws/connect/main.zip

	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/ws/disconnect/main ./lambdas/ws/disconnect
	zip -j ./bin/lambdas/ws/disconnect/main.zip ./bin/lambdas/ws/disconnect/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-ws-disconnect --zip-file fileb://./bin/lambdas/ws/disconnect/main.zip

	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/authorize/user/main ./lambdas/authorize/user
	zip -j ./bin/lambdas/authorize/user/main.zip ./bin/lambdas/authorize/user/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-authorize-user --zip-file fileb://./bin/lambdas/authorize/user/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/authorize/session/main ./lambdas/authorize/session
	zip -j ./bin/lambdas/authorize/session/main.zip ./bin/lambdas/authorize/session/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-authorize-session --zip-file fileb://./bin/lambdas/authorize/session/main.zip

	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/users/create/main ./lambdas/users/create
	zip -j ./bin/lambdas/users/create/main.zip ./bin/lambdas/users/create/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-users-create --zip-file fileb://./bin/lambdas/users/create/main.zip

	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/users/delete/main ./lambdas/users/delete
	zip -j ./bin/lambdas/users/delete/main.zip ./bin/lambdas/users/delete/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-users-delete --zip-file fileb://./bin/lambdas/users/delete/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/users/get/main ./lambdas/users/get
	zip -j ./bin/lambdas/users/get/main.zip ./bin/lambdas/users/get/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-users-get --zip-file fileb://./bin/lambdas/users/get/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/users/update/main ./lambdas/users/update
	zip -j ./bin/lambdas/users/update/main.zip ./bin/lambdas/users/update/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-users-update --zip-file fileb://./bin/lambdas/users/update/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/mailboxes/create/main ./lambdas/mailboxes/create
	zip -j ./bin/lambdas/mailboxes/create/main.zip ./bin/lambdas/mailboxes/create/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-mailboxes-create --zip-file fileb://./bin/lambdas/mailboxes/create/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/mailboxes/get/main ./lambdas/mailboxes/get
	zip -j ./bin/lambdas/mailboxes/get/main.zip ./bin/lambdas/mailboxes/get/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-mailboxes-get --zip-file fileb://./bin/lambdas/mailboxes/get/main.zip
	
	GOOS=linux GOGC=3000 go build -o ./bin/lambdas/mailboxes/mailboxmail/exchange/main ./lambdas/mailboxes/mailboxmail/exchange
	zip -j ./bin/lambdas/mailboxes/mailboxmail/exchange/main.zip ./bin/lambdas/mailboxes/mailboxmail/exchange/main
	aws lambda update-function-code --function-name arn:aws:lambda:us-east-1:677645055400:function:snailmail-mailbox-mail-exchange --zip-file fileb://./bin/lambdas/mailboxes/mailboxmail/exchange/main.zip

	echo "lambda upload complete"
	