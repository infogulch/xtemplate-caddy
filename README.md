# caddy-xtemplate

> [!tip]
>
> [xtemplate][xtemplate] is a html/template-oriented hypertext preprocessor and
> rapid application development web server. To learn about how to use it, read
> the docs at https://github.com/infogulch/xtemplate

caddy-xtemplate adapts [xtemplate][xtemplate] for use in the [Caddy][caddy] web
server by:

1. Registering as a [Caddy module][extending-caddy] named
   [`http.handlers.xtemplate`][http.handlers.xtemplate] which exposes a
   `caddyhttp.MiddlewareHandler` that can serve as a route handler using the
   `xtemplate` handler middleware definition.
3. Adapts Caddyfile configuration to easily configure xtemplate through Caddy.

[xtemplate]: https://github.com/infogulch/xtemplate
[caddy]: https://caddyserver.com/
[extending-caddy]: https://caddyserver.com/docs/extending-caddy
[http.handlers.xtemplate]: https://caddyserver.com/download?package=github.com%2Finfogulch%2Fxtemplate-caddy

## Quickstart

First, [Download Caddy Server with `http.handlers.xtemplate` module][http.handlers.xtemplate], or [build it yourself](#build).

Write your caddy config and use the xtemplate http handler in a route block. See
[Config](#config) for a listing of xtemplate configs. The simplest Caddy config
is:

```Caddy
:8080

route {
    xtemplate
}
```

Place `.html` files in the directory specified by the `xtemplate.templates_dir`
key in your caddy config (default "templates"). The config above would load
templates from the `./templates` directory, relative to the current working
directory.

Run caddy with your config:

```shell
caddy run --config Caddyfile
```

> Caddy is a very capable http server, check out the caddy docs for features you
> may want to layer on top. Examples: set up an auth proxy, caching, rate
> limiting, automatic https, etc

## Config

Here are the xtemplate configs available to a Caddyfile:

```Caddy
xtemplate {
    templates_path <string>                            # The path to the templates directory. Default: "templates".
    template_extension <string>              # File extension to search for to find template files. Default ".html".
    delimiters <Left:string> <Right:string>  # The template action delimiters, default "{{" and "}}".

    watch_template_path <bool>   # Reloads templates if anything in template path changes. Default: true
    watch_context_path <bool>    # Reloads templates if anything in context path changes. Default: false

    funcs_modules <mod1> <mod2>  # A list of caddy modules under the `xtemplate.funcs.*`
                                 # namespace that implement the FuncsProvider interface,
                                 # to add custom funcs to the Template FuncMap.
}
```

> [!NOTE]
>
> `xtemplate-caddy` currently does not support configuring the dot context in
> the Caddyfile format. To configure the dot context you must use Caddy's json
> configuration.

## Build

### `xcaddy` CLI

To build xtemplate_caddy locally, install [`xcaddy`](xcaddy), then build from
the directory root. Examples:

```shell
# build with CGO in order to use the sqlite3 db driver
CGO_ENABLED=1 xcaddy build --with github.com/infogulch/xtemplate-caddy

# build enable the sqlite_json build tag to get json funcs
GOFLAGS='-tags="sqlite_json"' CGO_ENABLED=1 xcaddy build --with github.com/infogulch/xtemplate-caddy
```

[xcaddy]: https://github.com/caddyserver/xcaddy

<details>

```shell
TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format="%cd-%h"
```

</details>

### A Go module

Create a go module `go mod init <modname>` with a `main.go` like this:

```go
package main

import (
    caddycmd "github.com/caddyserver/caddy/v2/cmd"

    _ "github.com/infogulch/xtemplate-caddy"

    // Add xtemplate dot providers:
    _ "github.com/infogulch/xtemplate/providers"
    _ "github.com/infogulch/xtemplate/providers/nats"

    // Add other caddy modules:
    // _ "github.com/greenpau/caddy-security"
)

func main() {
	caddycmd.Main()
}
```

Compile it with `go build -o caddy .`, then run with `./caddy run --config Caddyfile`

## Package history

This package has moved several times. Here are some previous names it has been known as:

* `github.com/infogulch/caddy-xtemplate` - Initial implementation to prove out the idea.
* `github.com/infogulch/xtemplate/caddy` - Refactored xtemplate to be usable from the cli and as a Go library, split Caddy integration into a separate module in the same repo.
* `github.com/infogulch/xtemplate-caddy` **Current package** - Caddy integration moved to its own repo, and refactored config organization. This should be the final rename 🤞.
