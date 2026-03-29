package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g *Generator) netTS() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title("Network").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`sum(rate(process_network_receive_bytes_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} recv").
				Datasource(metricDatasource),
			prometheus.NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					`-sum(rate(process_network_transmit_bytes_total{instance=~"$%s"}[$__rate_interval])) by (instance)`,
					instanceVariableName,
				)).
				LegendFormat("{{instance}} send").
				Datasource(metricDatasource),
		}).
		Unit(units.BytesPerSecondIEC).
		Legend(simpleLegend()).
		Datasource(metricDatasource)
}
