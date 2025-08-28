package tickers

import "time"

func NewTickerChanWithInitial(interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	oldTicks := ticker.C
	newTicks := make(chan time.Time, 1)

	go func() {
		newTicks <- time.Now()

		for tick := range oldTicks {
			newTicks <- tick
		}
	}()

	ticker.C = newTicks

	return ticker
}
