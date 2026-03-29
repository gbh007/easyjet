package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
)

type Generator struct {
	uid    string
	vlExpr string
}

func New(uid, vlExpr string) *Generator {
	plugins.RegisterDefaultPlugins()
	return &Generator{
		uid:    uid,
		vlExpr: vlExpr,
	}
}
