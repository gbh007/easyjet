package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g *Generator) withTagsAndAnnotations(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
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

	builder.Annotation(
		dashboard.
			NewAnnotationQueryBuilder().
			Enable(true).
			Expr(`sum(changes(easyjet_start_timestamp{instance=~"$instance"}[$__interval])) by (instance)`).
			IconColor("super-light-blue").
			Placement(dashboard.AnnotationQueryPlacementInControlsMenu).
			Name("app started (metrics)").
			Datasource(metricDatasource),
	)

	return builder
}
