package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	do "gopkg.in/godo.v2"
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

	p.Task("db-checkpid", nil, func(c *do.Context) {
		if !exist("/tmp/postgres-9456.lock") && exist("data/postmaster.pid") {
			c.Run("rm data/postmaster.pid")
		}
	})

	p.Task("db-start", do.S{"db-create", "db-checkpid"}, func(c *do.Context) {
		if !exist("data/postmaster.pid") {
			c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' start -w -t 120")
			c.Run("touch /tmp/postgres-9456.lock")
		}
	})

	p.Task("db-stop", do.S{"db-checkpid"}, func(c *do.Context) {
		if exist("data/postmaster.pid") {
			c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' stop -m fast")
			c.Run("rm /tmp/postgres-9456.lock")
		}
	})

	p.Task("db-restart", do.S{"db-stop", "db-start"}, nil)

	p.Task("db-destory", do.S{"db-stop"}, func(c *do.Context) {
		if exist("data") {
			c.Run("rm -rf data")
		}
	})

	p.Task("db-init", do.S{"db-start"}, func(c *do.Context) {
		c.Run("createdb -p 9456 orangez")
		c.Run("psql --set ON_ERROR_STOP=on -p 9456 -d orangez -f tasks/init.sql")
	})

	p.Task("db-console", do.S{"db-start"}, func(c *do.Context) {
		c.Run("psql -p 9456 -d orangez")
	})

	p.Task("db-reinit", do.S{"db-destory", "db-init", "db-restore"}, nil)

	p.Task("db-dump", do.S{"db-start"}, func(c *do.Context) {
		c.Run("pg_dump -p 9456 -d orangez -T http_sessions -a -f tasks/dump.sql")
	})

	p.Task("db-restore", do.S{"db-start"}, func(c *do.Context) {
		c.Run("psql --set ON_ERROR_STOP=on -p 9456 -d orangez -f tasks/dump.sql")
	})

	p.Task("db-sessions", do.S{"db-start"}, func(c *do.Context) {
		c.Run("pg_dump -p 9456 -d orangez -a -t http_sessions")
	})

	p.Task("cert-init", nil, func(c *do.Context) {
		if !exist("cert/intermediate.cert.pem") {
			c.Run("sudo sh init.sh", do.M{"$in": "cert"})
		}
	})

	p.Task("cert-generate", do.S{"cert-init"}, func(c *do.Context) {
		if !exist("cert/orangez.cert.bundle.pem") {
			c.Run("sh generate.sh 127.0.0.69 orangez", do.M{"$in": "cert"})
		}
	})

	p.Task("cert-regen", do.S{"cert-init"}, func(c *do.Context) {
		c.Run("sh generate.sh 127.0.0.69 orangez", do.M{"$in": "cert"})
	})

	p.Task("prepare", do.S{"db-start", "cert-generate"}, nil)

	p.Task("run", do.S{"prepare"}, func(c *do.Context) {
		fmt.Println("Running...")
		c.Run("go run main.go")
	})

	p.Task("test", do.S{"prepare"}, func(c *do.Context) {
		c.Run("go test -v")
	})
}

func main() {
	do.Godo(tasks)
}
