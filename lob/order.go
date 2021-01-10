package lob

// Order represents a resting set of bids or asks in the lob.
type Order interface {
	Type() SIDE
	OrderPrice() float64
	OrderSize() int
}
