# Kulinari (Server)
A recipe management program

## Installation

### 1. Install Go

If you are using a Linux distribution you can usually find a package in your repositories. If not or if you are using a different operating system you can download Go at [https://golang.org/dl/](https://golang.org/dl/)

### 2. Download dependencies and build binary
(for Bash shell only)

```
go get github.com/gorilla/mux && \
go get github.com/gorilla/securecookie && \
go get github.com/go-sql-driver/mysql && \
go get gopkg.in/russross/blackfriday.v2 && \
go get github.com/schlagma/kulinari-server && \
cd $GOPATH/src/github.com/schlagma/kulinari-server && \
go build main.go api.go categories.go cookie.go data.go login.go recipe-create.go recipe-delete.go recipe-edit.go recipe-share.go recipe.go settings.go shopping.go && \
cp $GOPATH/src/github.com/schlagma/kulinari-server/config/config.sample.json $GOPATH/src/github.com/schlagma/kulinari-server/config/config.json && \
./main
```

After you have set up your database and inserted the data into the config.json file you can access Kulinari at localhost:2805.
