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

func (g *Generator) cpuUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(rate(process_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)
	query3 := fmt.Sprintf(
		`sum(rate(go_cpu_classes_user_cpu_seconds_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("CPU usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process %").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query3).
				LegendFormat("{{instance}} go code %").
				Datasource(metricDatasource),
		}).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.PercentUnit,
			},
		}).
		OverrideByQuery("C", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.PercentUnit,
			},
		}).
		Unit(units.Seconds).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
