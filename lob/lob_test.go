package lob_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/notaquant/ob/lob"
)

func ExampleLimitOrderBook() {
	b := lob.NewLimitOrderBook()
	// Add asks
	// price, size
	b.Ask(1000, 1)
	b.Ask(1010, 2)
	b.Ask(999, 3)
	b.Ask(1200, 4)

	// add some bids
	b.Bid(998, 4)
	b.Bid(100, 3)
	b.Bid(997, 2)
	b.Bid(900, 1)

	bid := b.Peek(997, lob.BID)

	fmt.Println("Bid Size at ", bid.OrderPrice(), "is :", bid.OrderSize())

}

func TestLimitOrderBook(t *testing.T) {

	randomBook := func(side lob.SIDE, levels int, min, max int) *lob.LimitOrderBook {
		ob := lob.NewLimitOrderBook()
		p := 0.
		for t := 0; t < levels; t++ {
			if side == lob.ASK {
				for i := 0; i < 100; i++ {
					p = math.Max(float64(min), float64(rand.Intn(10000)))
					p = math.Min(float64(max), p)

					if ob.Peek(p, lob.ASK) == nil {
						break
					}
				}

				ob.Ask(p, rand.Intn(25)+1)
			} else {
				for i := 0; i < 100; i++ {
					p = math.Max(float64(min), float64(rand.Intn(10000)))
					p = math.Min(float64(max), p)

					if ob.Peek(p, lob.BID) == nil {
						break
					}
				}

				ob.Bid(p, rand.Intn(25)+1)
			}
		}

		return ob
	}

	t.Run("TestBidCorrectNumberOfLevels", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)
		b := randomBook(lob.BID, levels, 0, 1000)

		s := b.Snapshot(levels * 2)

		if len(s.Bids.Ticks()) == levels {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}

		s2 := b.Snapshot(levels >> 1)

		if len(s2.Bids.Ticks()) == levels>>1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}

		s3 := b.Snapshot(1)

		if len(s3.Bids.Ticks()) == 1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}
	})
	t.Run("TestAskCorrectNumberOfLevels", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)
		b := randomBook(lob.ASK, levels, 0, 1000)

		s := b.Snapshot(levels * 2)

		if len(s.Bids.Ticks()) == levels {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}

		s2 := b.Snapshot(levels >> 1)

		if len(s2.Asks.Ticks()) == levels>>1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}

		s3 := b.Snapshot(1)

		if len(s3.Asks.Ticks()) == 1 {
			t.Fatalf("failed to verify correct number of ticks got %d ticks expected %d ticks", len(s.Bids.Ticks()), levels)
		}
	})
	t.Run("TestBidDescSequential", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.BID, levels, 0, 1000)

		s := b.Snapshot(levels)

		for i := 0; i < len(s.Bids.Ticks()); i++ {
			l1 := s.Bids.Ticks()[i]
			l2 := s.Bids.Ticks()[i+1]

			if l1.Price < l2.Price {
				t.Fatal("failed to correctly parse sequentially descending bids")
			}

		}
	})
	t.Run("TestAskAscSequential", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.ASK, levels, 0, 1000)

		s := b.Snapshot(levels)

		for i := 0; i < len(s.Asks.Ticks())-1; i++ {
			l1 := s.Asks.Ticks()[i]
			l2 := s.Asks.Ticks()[i+1]

			if l1.Price > l2.Price {
				t.Fatal("failed to correctly parse sequentially ascending asls")
			}

		}
	})
	t.Run("TestBidNonZeroVol", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.BID, levels, 0, 1000)

		s := b.Snapshot(levels)

		if len(s.Bids.Ticks()) < 0 {
			t.Fatal("failed to assert tick length")
		}

		for _, bid := range s.Bids.Ticks() {
			if bid.Size == 0 {
				t.Fatal("failed to bid with non zero volume")
			}
		}
	})
	t.Run("TestAskNonZeroVol", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.ASK, levels, 0, 1000)

		s := b.Snapshot(levels)

		if len(s.Asks.Ticks()) < 0 {
			t.Fatal("failed to assert tick length")
		}

		for _, ask := range s.Asks.Ticks() {
			if ask.Size == 0 {
				t.Fatal("failed to bid with non zero volume")
			}
		}
	})
	t.Run("TestBidCheckQuote", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.BID, levels, 100, 10000)

		b.Bid(10001, 1)

		s := b.Snapshot(1)

		if s.Bids.Ticks()[0].Price != 10001 || s.Bids.Ticks()[0].Size != 1 {
			t.Fatal("failed to correctly bid at new price")
		}

	})
	t.Run("TestAskCheckQuote", func(t *testing.T) {

		levels := 1000 + rand.Intn(1000)

		b := randomBook(lob.BID, levels, 100, 10000)

		b.Bid(99, 1)

		s := b.Snapshot(1)

		if s.Bids.Ticks()[0].Price != 99 || s.Bids.Ticks()[0].Size != 1 {
			t.Fatal("failed to correctly ask at new price")
		}

	})
}
