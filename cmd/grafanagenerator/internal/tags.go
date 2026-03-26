package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g *Generator) WithTagsAndAnnotations(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	tags := []string{"easyjet"}
	builder.Tags(tags)

	builder.Link(
		dashboard.
			NewDashboardLinkBuilder("GitHub Repository").
			Url("https://github.com/gbh007/easyjet").
			Type(dashboard.DashboardLinkTypeLink).
			TargetBlank(true),
	)

	builder.Link(
		dashboard.
			NewDashboardLinkBuilder("EasyJet Dashboards").
			Tags(tags).
			Type(dashboard.DashboardLinkTypeDashboards).
			KeepTime(true).
			AsDropdown(true).
			TargetBlank(true),
	)

	return builder
}
