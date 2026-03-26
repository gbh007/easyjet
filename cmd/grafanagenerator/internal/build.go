package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

// Build builds the Grafana dashboard.
func (g *Generator) Build() (dashboard.Dashboard, error) {
	builder := dashboard.
		NewDashboardBuilder("EasyJet").
		Uid(g.uid).
		Timezone("Asia/Krasnoyarsk").
		Time("now-6h", "now").
		WeekStart("monday").
		Refresh("1m").
		Tooltip(dashboard.DashboardCursorSyncCrosshair)

	g.WithPanels(builder)
	g.WithVariables(builder)
	g.WithTagsAndAnnotations(builder)

	d, err := builder.Build()
	if err != nil {
		return dashboard.Dashboard{}, fmt.Errorf("build dashboard: %w", err)
	}

	return d, nil
}
