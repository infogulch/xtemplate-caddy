package xtemplate_caddy

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/infogulch/watch"
	"github.com/infogulch/xtemplate"
	"go.uber.org/zap/exp/zapslog"
	"golang.org/x/exp/slices"
)

func init() {
	caddy.RegisterModule(XTemplateModule{})
}

// CaddyModule returns the Caddy module information.
func (XTemplateModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.xtemplate",
		New: func() caddy.Module { return new(XTemplateModule) },
	}
}

type XTemplateModule struct {
	xtemplate.Config

	WatchTemplatePath bool `json:"watch_template_path"`
	WatchContextPath  bool `json:"watch_context_path"`

	FuncsModules []string `json:"funcs_modules,omitempty"`

	handler xtemplate.CancelHandler
	halt    chan<- struct{}
}

// Validate ensures t has a valid configuration. Implements caddy.Validator.
func (m *XTemplateModule) Validate() error {
	if m.Database.Driver != "" && slices.Index(sql.Drivers(), m.Database.Driver) == -1 {
		return fmt.Errorf("database driver '%s' does not exist", m.Database.Driver)
	}
	return nil
}

// Provision provisions t. Implements caddy.Provisioner.
func (m *XTemplateModule) Provision(ctx caddy.Context) error {
	// Wrap zap logger into a slog logger for xtemplate
	log := slog.New(zapslog.NewHandler(ctx.Logger().Core(), nil)).WithGroup("xtemplate-caddy")

	m.Logger = log

	var err error
	m.handler, err = xtemplate.Build(&m.Config)
	if err != nil {
		return err
	}

	var watchDirs []string
	if m.WatchTemplatePath {
		watchDirs = append(watchDirs, m.Template.Path)
	}
	if m.WatchContextPath {
		watchDirs = append(watchDirs, m.Context.Path)
	}

	if len(watchDirs) > 0 {
		changed, halt, err := watch.WatchDirs(watchDirs, 200*time.Millisecond)
		if err != nil {
			return err
		}
		m.halt = halt
		watch.React(changed, halt, func() (halt bool) {
			temphandler, err := xtemplate.Build(&m.Config)
			if err != nil {
				log.Info("failed to reload xtemplate", "error", err)
				return
			}
			temphandler, m.handler = m.handler, temphandler
			log.Info("reloaded templates after file changed")
			temphandler.Cancel()
			return
		})
	}
	return nil
}

func (m *XTemplateModule) ServeHTTP(w http.ResponseWriter, r *http.Request, _ caddyhttp.Handler) error {
	m.handler.ServeHTTP(w, r)
	return nil
}

// Cleanup discards resources held by t. Implements caddy.CleanerUpper.
func (m *XTemplateModule) Cleanup() error {
	if m.halt != nil {
		m.halt <- struct{}{}
		close(m.halt)
		m.halt = nil
	}
	if m.handler != nil {
		m.handler.Cancel()
		m.handler = nil
	}
	return nil
}

// Interface guards
var (
	_ caddy.Validator             = (*XTemplateModule)(nil)
	_ caddy.Provisioner           = (*XTemplateModule)(nil)
	_ caddyhttp.MiddlewareHandler = (*XTemplateModule)(nil)
	_ caddy.CleanerUpper          = (*XTemplateModule)(nil)
)
