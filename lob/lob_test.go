package lob

import (
	"math"
	"math/rand"
	"testing"
)

func TestLimitOrderBook(t *testing.T) {

	randomBook := func(side SIDE, levels int, min, max int) *LimitOrderBook {
		lob := NewLimitOrderBook()
		p := 0.
		for t := 0; t < levels; t++ {
			if side == ASK {
				for i := 0; i < 100; i++ {
					p = math.Max(float64(min), float64(rand.Intn(10000)))
					p = math.Min(float64(max), p)

					if lob.Peek(p, ASK) == nil {
						break
					}
				}

				lob.Ask(p, rand.Intn(25))
			} else {
				for i := 0; i < 100; i++ {
					p = math.Max(float64(min), float64(rand.Intn(10000)))
					p = math.Min(float64(max), p)

					if lob.Peek(p, BID) == nil {
						break
					}
				}

				lob.Bid(p, rand.Intn(25))
			}
		}

		return lob
	}

	t.Run("bid correct number of levels", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)
		b := randomBook(BID, levels, 0, 1000)

		s := b.Snapshot(levels * 2)

		if len(s.bids.ticks) == levels {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.bids.ticks), levels)
		}

		s2 := b.Snapshot(levels >> 1)

		if len(s2.bids.ticks) == levels>>1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.bids.ticks), levels)
		}

		s3 := b.Snapshot(1)

		if len(s3.bids.ticks) == 1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.bids.ticks), levels)
		}
	})
}
