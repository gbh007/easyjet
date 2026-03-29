package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) fdTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("File descriptors").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(process_open_fds{instance=~"$%s"}) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}}").
				Datasource(metricDatasource),
		}).
		Unit(units.Short).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
