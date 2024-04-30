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
			case "templates_path":
				if !h.AllArgs(&t.TemplatesDir) {
					return nil, h.ArgErr()
				}
			case "template_extension":
				if !h.AllArgs(&t.TemplateExtension) {
					return nil, h.ArgErr()
				}
			case "delimiters":
				if !h.AllArgs(&t.LDelim, &t.RDelim) {
					return nil, h.ArgErr()
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
