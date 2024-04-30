package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/infogulch/xtemplate-caddy"
	_ "github.com/infogulch/xtemplate/providers"
)

func main() {
	caddycmd.Main()
}
