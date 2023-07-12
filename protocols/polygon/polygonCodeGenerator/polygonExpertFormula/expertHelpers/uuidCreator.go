package experthelpers

import (
	"math/rand"
	"time"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/oklog/ulid"
)

//Generate transaction UUID
func GenerateTransactionUUID() string {
	timeNow := time.Now().UTC()
	logger := utilities.NewCustomLogger()
	entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
	logger.LogWriter("TXN UUID : "+id.String(), constants.INFO)
	return id.String()
}
