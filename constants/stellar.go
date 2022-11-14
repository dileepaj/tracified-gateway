package constants

import "github.com/stellar/go/txnbuild"

var (
	MinBaseFee         int64 = 100000
	TransactionTimeOut       = txnbuild.NewInfiniteTimeout()
)
