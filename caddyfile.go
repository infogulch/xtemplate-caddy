package xtemplate_caddy

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/infogulch/xtemplate"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("xtemplate", parseCaddyfile)
}

// parseCaddyfile sets up the handler from Caddyfile tokens.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	t := &XTemplateModule{
		// Inherit xtemplate defaults
		Config: *xtemplate.New(),

		// Set defaults on custom fields
		WatchTemplatePath: true,
	}

	for h.Next() {
		for h.NextBlock(0) {
			switch h.Val() {
			case "template":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					switch h.Val() {
					case "path":
						if !h.AllArgs(&t.Template.Path) {
							return nil, h.ArgErr()
						}
					case "template_extension":
						if !h.AllArgs(&t.Template.TemplateExtension) {
							return nil, h.ArgErr()
						}
					case "delimiters":
						if !h.AllArgs(&t.Template.Delimiters.Left, &t.Template.Delimiters.Right) {
							return nil, h.ArgErr()
						}
					default:
						return nil, h.Errf("unknown config option")
					}
				}
			case "context":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					switch h.Val() {
					case "path":
						if !h.AllArgs(&t.Context.Path) {
							return nil, h.ArgErr()
						}
					default:
						return nil, h.Errf("unknown config option")
					}
				}
			case "database":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					switch h.Val() {
					case "driver":
						if !h.AllArgs(&t.Database.Driver) {
							return nil, h.ArgErr()
						}
					case "connstr":
						if !h.AllArgs(&t.Database.Connstr) {
							return nil, h.ArgErr()
						}
					default:
						return nil, h.Errf("unknown config option")
					}
				}
			case "config":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					var key, val string
					key = h.Val()
					if _, ok := t.UserConfig[key]; ok {
						return nil, h.Errf("config key '%s' repeated", key)
					}
					if !h.Args(&val) {
						return nil, h.ArgErr()
					}
					t.UserConfig[key] = val
				}
			case "funcs_modules":
				t.FuncsModules = h.RemainingArgs()
			case "watch_template_path":
				var b string
				if !h.AllArgs(&b) {
					return nil, h.ArgErr()
				}
				switch b {
				case "true":
					t.WatchTemplatePath = true
				case "false":
					t.WatchTemplatePath = false
				default:
					return nil, h.Errf("arg must be bool, got: %s", b)
				}
			case "watch_context_path":
				var b string
				if !h.AllArgs(&b) {
					return nil, h.ArgErr()
				}
				switch b {
				case "true":
					t.WatchContextPath = true
				case "false":
					t.WatchContextPath = false
				default:
					return nil, h.Errf("arg must be bool, got: %s", b)
				}
			default:
				return nil, h.Errf("unknown config option")
			}
		}
	}
	return t, nil
}
