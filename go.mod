module gib

go 1.17

replace go-image-board => ./

require (
	github.com/disintegration/imageorient v0.0.0-20180920195336-8147d86e83ec
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/satori/go.uuid v1.2.0
	go-image-board v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220307211146-efcb8507fb70
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867
)

require (
	github.com/disintegration/gift v1.1.2 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
