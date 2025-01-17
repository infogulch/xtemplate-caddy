package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/infogulch/xtemplate-caddy"
)

func main() {
	caddycmd.Main()
}
