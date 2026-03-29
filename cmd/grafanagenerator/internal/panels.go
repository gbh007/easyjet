package internal

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
)

func (g *Generator) withPanels(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithRow(dashboard.NewRowBuilder("Task Execution"))
	builder.WithPanel(g.runDurationTS())
	builder.WithPanel(g.runRPSTS())
	if g.vlExpr != "" {
		builder.WithRow(dashboard.NewRowBuilder("Logs"))
		builder.WithPanel(g.logs())
	}
	builder.WithRow(dashboard.NewRowBuilder("Resource usage"))
	builder.WithPanel(g.cpuUsageTS())
	builder.WithPanel(g.memUsageTS())
	builder.WithPanel(g.threadsTS())
	builder.WithPanel(g.netTS())
	builder.WithPanel(g.fdTS())
	builder.WithPanel(g.schedTS())
	return builder
}

func simpleLegend() *common.VizLegendOptionsBuilder {
	return common.
		NewVizLegendOptionsBuilder().
		DisplayMode(common.LegendDisplayModeTable).
		Placement(common.LegendPlacementBottom).
		Calcs([]string{"mean", "lastNotNull"}).
		SortBy("Mean").
		SortDesc(true).
		ShowLegend(true)
}

func (g *Generator) logs() *logs.PanelBuilder {
	return logs.
		NewPanelBuilder().
		Title("Logs").
		Targets([]cog.Builder[variants.Dataquery]{
			loki. // Примечание по сигнатуре частично совпадает с Victoria Logs, т.ч. используем его.
				NewDataqueryBuilder().
				Expr(g.vlExpr),
		}).
		SortOrder(common.LogsSortOrderDescending).
		EnableLogDetails(true).
		Height(12).Span(24).
		Datasource(logsVictoriaLogsDatasource)
}
