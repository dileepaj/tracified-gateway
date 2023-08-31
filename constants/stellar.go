package constants

import "github.com/stellar/go/txnbuild"

var (
	MinBaseFee         int64 = 10000
	TransactionTimeOut       = txnbuild.NewInfiniteTimeout()
)
