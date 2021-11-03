build:
	go build ./...

lambda:
	go build -o authorize_user ./lambdas/authorize/user

	go build -o users_create ./lambdas/users/create
	go build -o users_delete ./lambdas/users/delete
	go build -o users_get ./lambdas/users/get
	go build -o users_update ./lambdas/users/update

	go build -o mailboxes_create ./lambdas/mailboxes/create
	go build -o mailboxes_get ./lambdas/mailboxes/get

	go build -o mailbox_mail_exchange ./lambdas/mailboxes/mailboxmail/exchange
	
	zip snailmail_lambdas.zip authorize_user users_create users_delete users_get users_update mailboxes_create mailboxes_get mailbox_mail_exchange
	
	rm users_create authorize_user users_delete users_get users_update mailboxes_create mailboxes_get mailbox_mail_exchange

mocks:
	mockery --all --dir ./resources --output ./test/resources/mocks