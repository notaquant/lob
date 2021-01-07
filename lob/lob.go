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
	asks  *Asks
	bids  *Bids
	depth int
}

// NewLimitOrderBook creates a new instance of a lob data structure.
func NewLimitOrderBook() *LimitOrderBook {
	return &LimitOrderBook{
		asks:  NewAsks(),
		bids:  NewBids(),
		depth: 0,
	}
}

// Bid at a price level with a given size.
func (lob *LimitOrderBook) Bid(price float64, size int) {

	lob.bids.SetBid(Bid{
		Price: price,
		Size:  size,
	})
}

// Ask at a price level with a given size.
func (lob *LimitOrderBook) Ask(price float64, size int) {
	lob.asks.SetAsk(Ask{
		Price: price,
		Size:  size,
	})
}

// Peek the cumulative sum of orders at a given price.
func (lob *LimitOrderBook) Peek(price float64, side SIDE) int {
	switch side {
	case ASK:
		return lob.asks.Peek(price)
	case BID:
		return lob.bids.Peek(price)
	}

	return 0
}

// Snapshot returns a snapshot of resting orders at a given level.
func (lob *LimitOrderBook) Snapshot(levels int) *LimitOrderBook {
	bids := lob.bids.Snapshot(levels)
	asks := lob.asks.Snapshot(levels)

	return &LimitOrderBook{
		bids: &Bids{head: lob.bids.head, ticks: bids},
		asks: &Asks{head: lob.asks.head, ticks: asks},
	}
}
