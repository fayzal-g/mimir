// SPDX-License-Identifier: AGPL-3.0-only

package validation

import (
	"strings"

	"github.com/grafana/mimir/pkg/mimirpb"
)

// For each write request, the group is obtained from the first non-empty group label from the list of timeseries
func GroupLabel(o *Overrides, userID string, timeseries []mimirpb.PreallocTimeseries) string {
	if len(timeseries) == 0 {
		return ""
	}

	groupLabel := o.SeparateMetricsGroupLabel(userID)
	if groupLabel == "" {
		// If not set, label value will be "" and dropped by Prometheus
		return groupLabel
	}

	for _, label := range timeseries[0].Labels {
		if label.Name == groupLabel {
			// label.Value string is cloned as underlying PreallocTimeseries contains
			// unsafe strings that should not be retained
			return strings.Clone(label.Value)
		}
	}

	return ""
}
