package lob

// Bid represent a resting Bid in the lob.
type Bid struct {
	Size  int
	Price float64
}

// Type represents a generic interface implementation for a Bid or an Ask.
func (bid *Bid) Type() SIDE {
	return BID
}

// OrderSize represents a generic interface implement for a Bid or an Ask.
func (bid *Bid) OrderSize() int {
	return bid.Size
}

// OrderPrice represents a generic interface implement for a Bid or an Ask.
func (bid *Bid) OrderPrice() float64 {
	return bid.Price
}

// Bids represent the entire Bid side of the lob.
type Bids struct {
	head  int
	ticks []Bid
}

// NewBids creates a new instance of the Bids side
func NewBids() *Bids {

	return &Bids{
		head:  0,
		ticks: make([]Bid, 0),
	}
}

// Head returns the head item of the bids list.
func (bids *Bids) Head() int {
	return bids.head
}

// Ticks returns the ticks populating the bid side.
func (bids *Bids) Ticks() []Bid {
	return bids.ticks
}

// Peek returns the sum of Bids at a certain price.
func (bids *Bids) Peek(price float64) *Bid {

	if len(bids.ticks) == 0 {
		return nil
	}
	i := bids.Find(price)
	item := bids.ticks[i]

	if item.Price == price {
		return &item
	}
	return nil
}

// SetBid sets an resting Bid at certain price and size
// if there's already an Bid at that price we incremt the size.
func (bids *Bids) SetBid(Bid Bid) {

	// if it's the first order in the Bid side set head to 0.
	if len(bids.ticks) == 0 {
		bids.ticks = append(bids.ticks, Bid)
		bids.head = 0
		return
	}

	i := bids.Find(Bid.Price)

	item := bids.ticks[i]
	itemhead := bids.ticks[bids.head]
	// found a matching price level for this Bid
	if item.Price == Bid.Price {
		// increment Bid size of the level
		item.Size += Bid.Size
		if bids.head == i && Bid.Size == 0 {
			bids.head = bids.scan(i + 1)
		} else if Bid.Price < bids.ticks[bids.head].Price {
			bids.head = i
		}
		return
	}

	// its a new price level
	before := Bid.Price < item.Price
	ins := 0
	if before {
		ins = i
	} else {
		ins = i + 1
	}
	bids.insert(ins, Bid)

	if before && Bid.Price < itemhead.Price {
		bids.head = bids.scan(ins)
	}
}

// Find does a binary search on the Bid side by price.
func (bids *Bids) Find(price float64) int {

	BidSide := bids.ticks
	n := len(BidSide)

	// Edge cases
	if price <= BidSide[0].Price {
		return 0
	} else if price >= BidSide[n-1].Price {
		return n - 1
	}

	i := 0
	j := n
	mid := 0

	for i < j {
		mid = (i + j) >> 1

		if price == BidSide[mid].Price {
			return mid
		} else if price < BidSide[mid].Price {
			if mid > 0 && price > BidSide[mid-1].Price {
				return bids.getClosestMatch(mid-1, mid, price)
			}
			j = mid
		} else {
			if mid < n-1 && price < BidSide[mid+1].Price {
				return bids.getClosestMatch(mid, mid+1, price)
			}
			i = mid + 1
		}
	}
	return 0
}

// getClosestMatch fetches the closest match of price between two values a and b.
func (bids *Bids) getClosestMatch(a int, b int, price float64) int {
	if price-bids.ticks[a].Price >= bids.ticks[b].Price-price {
		return b
	}
	return a
}

// scan step through existing levels from index `from` => upwards.
func (bids *Bids) scan(from int) int {

	for t := from; t < len(bids.ticks); t++ {
		if bids.ticks[t].Size > 0 {
			return t
		}
	}

	return from
}

// insert an Bid at a given index
func (bids *Bids) insert(index int, Bid Bid) {
	if len(bids.ticks) == index { // nil or empty slice or after last element
		bids.ticks = append(bids.ticks, Bid)
	}
	bids.ticks = append(bids.ticks[:index+1], bids.ticks[index:]...) // index < len(a)
	bids.ticks[index] = Bid
}

// Snapshot returns a snapshot of the lob for a given number of levels.
func (bids *Bids) Snapshot(levels int) []Bid {

	lob := make([]Bid, 0)

	if len(bids.ticks) == 0 {
		return lob
	}

	for t := bids.head; t < len(bids.ticks); t++ {
		if bids.ticks[t].Size > 0 {
			lob = append(lob, bids.ticks[t])
		} else if len(lob) == levels {
			break
		}

	}
	return lob
}
