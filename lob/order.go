package lob

// Order represents a resting set of bids or asks in the lob.
type Order interface {
	Type() SIDE
}
