package ticker

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestNewTimeTicker(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("Instance can be created", func(t *testing.T) {
		var rate = time.Duration(fakerInstance.IntBetween(0, 999_999_999))

		ticker := NewTimeTicker(rate)

		require.IsType(t, &TimeTicker{}, ticker)
		require.Equal(t, rate, ticker.rate)
	})
}

func TestTimeTicker_Init(t *testing.T) {
	t.Run("Ticker can be initialized", func(t *testing.T) {
		var rate = time.Second

		ticker := TimeTicker{
			rate: rate,
		}

		times := ticker.Init()

		require.IsType(t, make(<-chan time.Time), times)
		require.GreaterOrEqual(t, -time.Now().Sub(<-times), rate)
	})
}

func TestTimeTicker_Rate(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("The current rate is returned", func(t *testing.T) {
		var rate = time.Duration(fakerInstance.IntBetween(0, 999_999_999))

		ticker := TimeTicker{
			rate: rate,
		}

		require.Equal(t, rate, ticker.rate)
		require.Equal(t, rate, ticker.Rate())
	})
}
