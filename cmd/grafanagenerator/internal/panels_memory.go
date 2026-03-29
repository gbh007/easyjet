package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) memUsageTS() *timeseries.PanelBuilder {
	query := fmt.Sprintf(
		`sum(process_resident_memory_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)
	query2 := fmt.Sprintf(
		`sum(go_memstats_sys_bytes{instance=~"$%s"}) by (instance)`,
		instanceVariableName,
	)

	return timeseries.NewPanelBuilder().
		Title("Memory usage").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{instance}} process resident").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(query2).
				LegendFormat("{{instance}} go sys").
				Datasource(metricDatasource),
		}).
		Unit(units.BytesIEC).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
