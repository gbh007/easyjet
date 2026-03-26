package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) WithPanels(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithRow(dashboard.NewRowBuilder("Resource usage"))
	builder.WithPanel(g.cpuUsageTS())
	builder.WithPanel(g.cpuUsageTSK8s())
	builder.WithPanel(g.memUsageTS())
	builder.WithPanel(g.memUsageGoTS())
	builder.WithRow(dashboard.NewRowBuilder("Task Execution"))
	builder.WithPanel(g.runDurationTS())
	return builder
}

func (g *Generator) runDurationTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`histogram_quantile(0.95, 
			sum(
				rate(easyjet_core_run_duration_bucket{project_id=~"$%s", instance=~"$%s"}[$__rate_interval])
			) by (le, project_id, result)
		)`,
		projectIDVariableName,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("P95 Run Duration").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{project_id}} - {{result}}").
				Datasource(metricDatasource),
		}).
		Unit(units.Seconds).
		Datasource(metricDatasource)
}

func (g *Generator) cpuUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(rate(process_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("CPU usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}}").
				Datasource(metricDatasource),
		}).
		Unit(units.Seconds).
		Datasource(metricDatasource)
}

func (g *Generator) cpuUsageTSK8s() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(rate(process_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("CPU usage (k8s format)").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}}").
				Datasource(metricDatasource),
		}).
		Unit(units.PercentUnit).
		Datasource(metricDatasource)
}

func (g *Generator) memUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(process_resident_memory_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("Memory usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}}").
				Datasource(metricDatasource),
		}).
		Unit(units.BytesIEC).
		Datasource(metricDatasource)
}

func (g *Generator) memUsageGoTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(go_memstats_sys_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("Memory usage (Go)").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}}").
				Datasource(metricDatasource),
		}).
		Unit(units.BytesIEC).
		Datasource(metricDatasource)
}
