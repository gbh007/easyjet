package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
)

type Generator struct {
	uid string
}

func New(uid string) *Generator {
	plugins.RegisterDefaultPlugins()
	return &Generator{
		uid: uid,
	}
}
