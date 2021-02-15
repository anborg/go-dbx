### Database structuring in golang

https://www.alexedwards.net/blog/organising-database-access

mkdir bookstore && cd bookstore
$ mkdir models
$ touch main.go models/models.go
GO111MODULE=on go mod init org.bookstore


Program Structure: 

dbconnection.go 

- Responsible for getting conncetion for any db type pg/ora/mysql/mssql
- Exposes a common struct Datasource
- Default properties for connection 

models.go
- All Domain objects, like Book
- All Repo db, like BookRepo
- 

main.go
- Read Datasource from config. 
- Get connection from dbconnection.go
- Initialize BookModel with db-repo obj
- for webservice, wire "/books" to bookHandler(w,r)
- listen :3000

TODO
- test oracle - DONE
- test pg - DONE
- test mysql 
- test mssql
- 