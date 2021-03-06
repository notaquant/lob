# Limit Order Book Data Structure in Go

This package implements a limit order book management system, it can rebuild, simulate, calculate live anaylitcs
on given trade data (Historical data in CSV format or Live WebSocket feeds).

* Thanks to [@AziDynamics](https://twitter.com/azidynamics) for the initial javascript [implementation](https://github.com/azidyn/lob) of a LOB which this implementation is a port.

## Usage

```golang
package main

import (
    "github.com/notaquant/ob/lob"
    "fmt"
)

func main() {
    csvStream := lob.NewStream("BITMEX-21-DEC-2020.csv")
    lob := ob.New().Build(csvStream)
    err := lob.Simulate()
    if err != nil {
        panic(err)
    }
    fmt.Println("Best Bid", lob.BestBid())
    fmt.Println("Best Ask", lob.BestAsk())
}
```

## Analytics

The following data can be extracted :

- Bid and Ask Volume
- Volume Delta
- Distribution of Volume relative to Price
- Open Interest
- Funding Rate (Interest Rate on Perpetual Swaps)
