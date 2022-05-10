module github.com/wspowell/snailmail

go 1.18

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/aws/aws-sdk-go v1.42.35
	github.com/caarlos0/env/v6 v6.9.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/google/uuid v1.3.0
	github.com/wspowell/context v0.0.8
	github.com/wspowell/errors v0.3.0
	github.com/wspowell/log v0.0.11
	github.com/wspowell/spiderweb v0.0.9
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e
	gorm.io/driver/mysql v1.2.3
	gorm.io/gorm v1.22.5
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/fasthttp/router v1.4.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/savsgio/gotils v0.0.0-20210907153846-c06938798b52 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.30.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
)

replace (
	github.com/ugorji/go => github.com/ugorji/go v1.2.6
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v1.2.6
	github.com/wspowell/spiderweb => github.com/wspowell/spiderweb v0.0.9-generics-v.0.0.1
)
