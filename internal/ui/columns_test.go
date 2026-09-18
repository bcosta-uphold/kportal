package ui

import (
	"testing"

	"github.com/lukaszraczylo/kportal/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestResolveColumns_Default(t *testing.T) {
	cols := ResolveColumns(nil)
	assert.Equal(t, defaultColumns, cols)

	cols = ResolveColumns(&config.Config{})
	assert.Equal(t, defaultColumns, cols)
}

func TestResolveColumns_CustomOrderAndVisibility(t *testing.T) {
	cfg := &config.Config{
		TUI: &config.TUISpec{
			Columns: []config.TableColumn{
				{Name: "status"},
				{Name: "alias", Width: 40},
			},
		},
	}

	cols := ResolveColumns(cfg)

	assert.Len(t, cols, 2)
	assert.Equal(t, ColKeyStatus, cols[0].Key)
	assert.Equal(t, ColKeyAlias, cols[1].Key)
	assert.Equal(t, 40, cols[1].Width)
}

func TestResolveColumns_UsesDefaultWidthWhenUnset(t *testing.T) {
	cfg := &config.Config{
		TUI: &config.TUISpec{
			Columns: []config.TableColumn{
				{Name: "namespace"},
			},
		},
	}

	cols := ResolveColumns(cfg)

	assert.Len(t, cols, 1)
	assert.Equal(t, ColumnWidthNamespace, cols[0].Width)
}

func TestResolveColumns_UnknownNamesSkipped(t *testing.T) {
	cfg := &config.Config{
		TUI: &config.TUISpec{
			Columns: []config.TableColumn{
				{Name: "bogus"},
			},
		},
	}

	// All entries invalid -> fall back to defaults.
	assert.Equal(t, defaultColumns, ResolveColumns(cfg))
}

func TestColumnValue(t *testing.T) {
	fwd := &ForwardStatus{
		Context:    "prod",
		Namespace:  "default",
		Alias:      "api",
		Type:       "service",
		Resource:   "api",
		RemotePort: 8080,
		LocalPort:  8081,
		Status:     "Active",
	}

	assert.Equal(t, "prod", columnValue(ResolvedColumn{Key: ColKeyContext, Width: 10}, fwd))
	assert.Equal(t, "8080", columnValue(ResolvedColumn{Key: ColKeyRemote}, fwd))
	assert.Equal(t, "8081", columnValue(ResolvedColumn{Key: ColKeyLocal}, fwd))
	assert.Equal(t, "Active", columnValue(ResolvedColumn{Key: ColKeyStatus}, fwd))
}

func TestColumnIndex(t *testing.T) {
	cols := []ResolvedColumn{{Key: ColKeyContext}, {Key: ColKeyStatus}}
	assert.Equal(t, 1, columnIndex(cols, ColKeyStatus))
	assert.Equal(t, -1, columnIndex(cols, ColKeyLocal))
}
