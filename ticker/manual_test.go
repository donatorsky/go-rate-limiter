package ticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewManualTicker(t *testing.T) {
	t.Run("Instance can be created", func(t *testing.T) {
		require.IsType(t, &ManualTicker{}, NewManualTicker())
	})
}

func TestManualTicker_Init(t *testing.T) {
	t.Run("Ticker can be initialized", func(t *testing.T) {
		ticker := ManualTicker{
			ticks: make(chan time.Time),
		}

		require.IsType(t, make(<-chan time.Time), ticker.Init())
	})
}

func TestManualTicker_Tick(t *testing.T) {
	t.Run("Tick is generated", func(t *testing.T) {
		ticker := ManualTicker{
			ticks: make(chan time.Time, 1),
		}

		require.Len(t, ticker.ticks, 0)
		ticker.Tick()
		require.Len(t, ticker.ticks, 1)
	})
}
