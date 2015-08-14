package main

import (
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	p.Task("dbinit", nil, func(c *do.Context) {
		c.Run("pg_ctl initdb -D data")
	})

	p.Task("dbstart", nil, func(c *do.Context) {
		c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' start -w -t 120")
	})

	p.Task("dbstop", nil, func(c *do.Context) {
		c.Run("pg_ctl -D data -o '--config-file=postgresql.conf' stop -m fast")
	})
}

func main() {
	do.Godo(tasks)
}
