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

func (g *Generator) threadsTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Gorutines & threads").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} gorutines").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(process_num_threads{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} process threads").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_threads{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} go threads").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_sched_goroutines_not_in_go_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} gorutines syscall or CGO").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_sched_goroutines_runnable_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} runnable gorutines").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_sched_goroutines_running_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} running gorutines").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(go_sched_goroutines_waiting_goroutines{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} waiting gorutines").
				Datasource(metricDatasource),
		}).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.Number,
			},
		}).
		OverrideByQuery("C", []dashboard.DynamicConfigValue{
			{
				Id:    "unit",
				Value: units.Number,
			},
		}).
		Unit(units.Short).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
