package lob

// Ask represent a resting Ask in the lob.
type Ask struct {
	Size  int
	Price float64
}

// Asks represent the entire Ask side of the lob.
type Asks struct {
	head  int
	ticks []Ask
}

// NewAsks creates a new instance of the Asks side
func NewAsks() *Asks {

	return &Asks{
		head:  0,
		ticks: make([]Ask, 0),
	}
}

// Peek returns the sum of asks at a certain price.
func (asks *Asks) Peek(price float64) int {

	if len(asks.ticks) == 0 {
		return 0
	}

	return 1
}

// SetAsk sets an resting ask at certain price and size
// if there's already an ask at that price we incremt the size.
func (asks *Asks) SetAsk(ask Ask) {

	// if it's the first order in the ask side set head to 0.
	if len(asks.ticks) == 0 {
		asks.ticks = append(asks.ticks, ask)
		asks.head = 0
		return
	}

	i := asks.Find(ask.Price)

	item := asks.ticks[i]
	itemhead := asks.ticks[asks.head]
	// found a matching price level for this ask
	if item.Price == ask.Price {
		// increment ask size of the level
		item.Size += ask.Size
		if asks.head == i && ask.Size == 0 {
			asks.head = asks.scan(i + 1)
		} else if ask.Price < asks.ticks[asks.head].Price {
			asks.head = i
		}
		return
	}

	// its a new price level
	before := ask.Price < item.Price
	ins := 0
	if before {
		ins = i
	} else {
		ins = i + 1
	}
	asks.insert(ins, ask)

	if before && ask.Price < itemhead.Price {
		asks.head = asks.scan(ins)
	}
}

// Find does a binary search on the ask side by price.
func (asks *Asks) Find(price float64) int {

	askSide := asks.ticks
	n := len(askSide)

	// Edge cases
	if price <= askSide[0].Price {
		return 0
	} else if price >= askSide[n-1].Price {
		return n - 1
	}

	i := 0
	j := n
	mid := 0

	for i < j {
		mid = (i + j) >> 1

		if price == askSide[mid].Price {
			return mid
		} else if price < askSide[mid].Price {
			if mid > 0 && price > askSide[mid-1].Price {
				return asks.getClosestMatch(mid-1, mid, price)
			}
			j = mid
		} else {
			if mid < n-1 && price < askSide[mid+1].Price {
				return asks.getClosestMatch(mid, mid+1, price)
			}
			i = mid + 1
		}
	}
	return 0
}

// getClosestMatch fetches the closest match of price between two values a and b.
func (asks *Asks) getClosestMatch(a int, b int, price float64) int {
	if price-asks.ticks[a].Price >= asks.ticks[b].Price-price {
		return b
	}
	return a
}

// scan step through existing levels from index `from` => upwards.
func (asks *Asks) scan(from int) int {

	for t := from; t < len(asks.ticks); t++ {
		if asks.ticks[t].Size > 0 {
			return t
		}
	}

	return from
}

// insert an ask at a given index
func (asks *Asks) insert(index int, ask Ask) {
	if len(asks.ticks) == index { // nil or empty slice or after last element
		asks.ticks = append(asks.ticks, ask)
	}
	asks.ticks = append(asks.ticks[:index+1], asks.ticks[index:]...) // index < len(a)
	asks.ticks[index] = ask
}

// Snapshot returns a snapshot of the lob for a given number of levels.
func (asks *Asks) Snapshot(levels int) []Ask {

	lob := make([]Ask, 0)

	if len(asks.ticks) == 0 {
		return lob
	}

	for t := asks.head; t < len(asks.ticks); t++ {
		if asks.ticks[t].Size > 0 {
			lob = append(lob, asks.ticks[t])
		} else if len(lob) == levels {
			break
		}

	}
	return lob
}
