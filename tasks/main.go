package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	do "gopkg.in/godo.v2"
	"os"
)

func exist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func connect() *sql.DB {
	db, err := sql.Open("postgres", "port=9456 dbname=orangez sslmode=disable")
	if err != nil {
		panic("Open postgres failed")
	}
	return db
}

func tasks(p *do.Project) {
	p.Task("db-create", nil, func(c *do.Context) {
		if !exist("data/PG_VERSION") {
			c.Run("pg_ctl initdb -D data")
		}
	})

	p.Task("db-start", do.S{"db-create"}, func(c *do.Context) {
		if !exist("data/postmaster.pid") {
			c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' start -w -t 120")
		}
	})

	p.Task("db-stop", nil, func(c *do.Context) {
		if exist("data/postmaster.pid") {
			c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' stop -m fast")
		}
	})

	p.Task("db-destory", do.S{"db-stop"}, func(c *do.Context) {
		if exist("data") {
			c.Run("rm -rf data")
		}
	})

	p.Task("db-init", do.S{"db-start"}, func(c *do.Context) {
		c.Run("createdb -p 9456 orangez")
		db := connect()
		defer db.Close()
		_, err := db.Exec("CREATE TABLE authors ( id int primary key, name varchar(50) not null unique, description varchar(200) );")
		if err != nil {
			panic(fmt.Sprintf("Create table failed : %v", err))
		}
	})

	p.Task("db-console", do.S{"db-start"}, func(c *do.Context) {
		c.Run("psql -p 9456 -d orangez")
	})

	p.Task("db-reinit", do.S{"db-destory", "db-init", "db-restore"}, nil)

	p.Task("db-dump", do.S{"db-start"}, func(c *do.Context) {
		c.Run("pg_dump -p 9456 -d orangez -a -f tasks/dump.sql")
	})

	p.Task("db-restore", do.S{"db-start"}, func(c *do.Context) {
		c.Run("psql --set ON_ERROR_STOP=on -p 9456 -d orangez -f tasks/dump.sql")
	})

	p.Task("run", do.S{"db-start"}, func(c *do.Context) {
		c.Run("go run main.go")
	})

	p.Task("test", do.S{"db-start"}, func(c *do.Context) {
		c.Run("go test")
	})
}

func main() {
	do.Godo(tasks)
}
