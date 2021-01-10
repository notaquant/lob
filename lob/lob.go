package lob

// SIDE represents a book side either bid or ask
type SIDE int

var (
	// ASK represents the ask side
	ASK SIDE = 0
	// BID represents the bid side
	BID SIDE
)

// LimitOrderBook a limit order book.
type LimitOrderBook struct {
	Asks  *Asks
	Bids  *Bids
	depth int
}

// NewLimitOrderBook creates a new instance of a lob data structure.
func NewLimitOrderBook() *LimitOrderBook {
	return &LimitOrderBook{
		Asks:  NewAsks(),
		Bids:  NewBids(),
		depth: 0,
	}
}

// Bid at a price level with a given size.
func (lob *LimitOrderBook) Bid(price float64, size int) {

	lob.Bids.SetBid(Bid{
		Price: price,
		Size:  size,
	})
}

// Ask at a price level with a given size.
func (lob *LimitOrderBook) Ask(price float64, size int) {
	lob.Asks.SetAsk(Ask{
		Price: price,
		Size:  size,
	})
}

// Peek the cumulative sum of orders at a given price.
func (lob *LimitOrderBook) Peek(price float64, side SIDE) Order {
	switch side {
	case ASK:
		return lob.Asks.Peek(price)
	case BID:
		return lob.Bids.Peek(price)
	}

	return nil
}

// Snapshot returns a snapshot of resting orders at a given level.
func (lob *LimitOrderBook) Snapshot(levels int) *LimitOrderBook {
	bids := lob.Bids.Snapshot(levels)
	asks := lob.Asks.Snapshot(levels)

	return &LimitOrderBook{
		Asks:  &Asks{head: lob.Asks.head, ticks: asks},
		Bids:  &Bids{head: lob.Bids.head, ticks: bids},
		depth: 0,
	}
}
