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

func (g *Generator) schedTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Schedule & GC").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.80, 
						sum(
							rate(go_sched_latencies_seconds_bucket{instance=~"$%s"}[$__rate_interval])
						) by (le, instance)
					)`,
					instanceVariableName,
				)).
				LegendFormat("Sched P80 {{instance}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.95, 
						sum(
							rate(go_sched_latencies_seconds_bucket{instance=~"$%s"}[$__rate_interval])
						) by (le, instance)
					)`,
					instanceVariableName,
				)).
				LegendFormat("Sched P95 {{instance}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(rate(go_sched_latencies_seconds_count{instance=~"$%s"}[$__rate_interval])) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("Sched RPS {{instance}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.80, 
						sum(
							rate(go_sched_pauses_total_gc_seconds_bucket{instance=~"$%s"}[$__rate_interval])
						) by (le, instance)
					)`,
					instanceVariableName,
				)).
				LegendFormat("GC P80 {{instance}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`histogram_quantile(0.95, 
						sum(
							rate(go_sched_pauses_total_gc_seconds_bucket{instance=~"$%s"}[$__rate_interval])
						) by (le, instance)
					)`,
					instanceVariableName,
				)).
				LegendFormat("GC P95 {{instance}}").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(rate(go_sched_pauses_total_gc_seconds_count{instance=~"$%s"}[$__rate_interval])) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("GC RPS {{instance}}").
				Datasource(metricDatasource),
		}).
		OverrideByQuery("C", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.RequestsPerSecond,
			},
		}).
		OverrideByQuery("F", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.RequestsPerSecond,
			},
		}).
		Unit(units.Seconds).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
