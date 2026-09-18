package ui

import (
	"fmt"
	"strings"

	"github.com/lukaszraczylo/kportal/internal/config"
)

// ColumnKey identifies a forwards-table column.
type ColumnKey string

// Recognized column keys, matching the names accepted in tui.columns.
const (
	ColKeyContext   ColumnKey = "context"
	ColKeyNamespace ColumnKey = "namespace"
	ColKeyAlias     ColumnKey = "alias"
	ColKeyType      ColumnKey = "type"
	ColKeyResource  ColumnKey = "resource"
	ColKeyRemote    ColumnKey = "remote"
	ColKeyLocal     ColumnKey = "local"
	ColKeyStatus    ColumnKey = "status"
)

// ResolvedColumn describes a single column ready for rendering: its header
// text, truncation width, and identity.
type ResolvedColumn struct {
	Key    ColumnKey
	Header string
	Width  int
}

// defaultColumns is the built-in column set and order, used when the config
// does not specify tui.columns.
var defaultColumns = []ResolvedColumn{
	{Key: ColKeyContext, Header: "CONTEXT", Width: ColumnWidthContext},
	{Key: ColKeyNamespace, Header: "NAMESPACE", Width: ColumnWidthNamespace},
	{Key: ColKeyAlias, Header: "ALIAS", Width: ColumnWidthAlias},
	{Key: ColKeyType, Header: "TYPE", Width: ColumnWidthType},
	{Key: ColKeyResource, Header: "RESOURCE", Width: ColumnWidthResource},
	{Key: ColKeyRemote, Header: "REMOTE", Width: 0},
	{Key: ColKeyLocal, Header: "LOCAL", Width: 0},
	{Key: ColKeyStatus, Header: "STATUS", Width: 0},
}

// columnHeaders maps every known column key to its display header.
var columnHeaders = map[ColumnKey]string{
	ColKeyContext:   "CONTEXT",
	ColKeyNamespace: "NAMESPACE",
	ColKeyAlias:     "ALIAS",
	ColKeyType:      "TYPE",
	ColKeyResource:  "RESOURCE",
	ColKeyRemote:    "REMOTE",
	ColKeyLocal:     "LOCAL",
	ColKeyStatus:    "STATUS",
}

// ResolveColumns builds the ordered list of columns to render, based on the
// tui.columns configuration. A nil config, or one with no tui.columns entries,
// falls back to the built-in default column set and order.
func ResolveColumns(cfg *config.Config) []ResolvedColumn {
	var configured []config.TableColumn
	if cfg != nil {
		configured = cfg.GetTableColumns()
	}

	if len(configured) == 0 {
		return defaultColumns
	}

	defaultWidth := make(map[ColumnKey]int, len(defaultColumns))
	for _, c := range defaultColumns {
		defaultWidth[c.Key] = c.Width
	}

	resolved := make([]ResolvedColumn, 0, len(configured))
	for _, c := range configured {
		key := ColumnKey(strings.ToLower(strings.TrimSpace(c.Name)))
		header, ok := columnHeaders[key]
		if !ok {
			// Unknown column names are rejected by the validator before this
			// point is reached; skip defensively rather than render garbage.
			continue
		}

		width := c.Width
		if width <= 0 {
			width = defaultWidth[key]
		}

		resolved = append(resolved, ResolvedColumn{Key: key, Header: header, Width: width})
	}

	if len(resolved) == 0 {
		return defaultColumns
	}

	return resolved
}

// columnValue returns the display value for the given column of a forward,
// truncated to the column's width when applicable.
func columnValue(col ResolvedColumn, fwd *ForwardStatus) string {
	switch col.Key {
	case ColKeyContext:
		return truncate(fwd.Context, col.Width)
	case ColKeyNamespace:
		return truncate(fwd.Namespace, col.Width)
	case ColKeyAlias:
		return truncate(fwd.Alias, col.Width)
	case ColKeyType:
		return truncate(fwd.Type, col.Width)
	case ColKeyResource:
		return truncate(fwd.Resource, col.Width)
	case ColKeyRemote:
		return fmt.Sprintf("%d", fwd.RemotePort)
	case ColKeyLocal:
		return fmt.Sprintf("%d", fwd.LocalPort)
	case ColKeyStatus:
		return fwd.Status
	default:
		return ""
	}
}

// columnIndex returns the position of key within cols, or -1 if absent.
func columnIndex(cols []ResolvedColumn, key ColumnKey) int {
	for i, c := range cols {
		if c.Key == key {
			return i
		}
	}
	return -1
}
