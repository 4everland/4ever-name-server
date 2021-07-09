package fns

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"net/url"
)

const PluginName = "fns"

var log = clog.NewWithPlugin(PluginName)

func init() { plugin.Register(PluginName, setup) }

func setup(c *caddy.Controller) error {
	fns, err := fnsParse(c)
	if err != nil {
		return plugin.Error(fns.Name(), err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		fns.Next = next
		return fns
	})
	return nil
}

func fnsParse(c *caddy.Controller) (*FNS, error) {
	var (
		fns *FNS
		err error
		i   int
	)
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++

		fns, err = parseStanza(c)
		if err != nil {
			return fns, err
		}
	}

	return fns, nil
}

func parseStanza(c *caddy.Controller) (*FNS, error) {
	var (
		fns = new(FNS)
		err error
	)

	for c.NextBlock() {
		switch c.Val() {
		case "api-url":
			var u *url.URL
			args := c.RemainingArgs()
			if len(args) == 0 {
				return nil, c.ArgErr()
			}
			u, err = url.Parse(args[0])
			if err != nil {
				return fns, c.Err("invalid api-url")
			}

			if u.Scheme != "http" && u.Scheme != "https" {
				return fns, c.Err("api-url scheme not support")
			}

			fns.RRsRepo = NewRRsSDK(u.String())
		default:
			return fns, c.Errf("unknown property '%s'", c.Val())
		}
	}

	if fns.RRsRepo == nil {
		return fns, c.Err("missing api-url")
	}

	return fns, nil
}
