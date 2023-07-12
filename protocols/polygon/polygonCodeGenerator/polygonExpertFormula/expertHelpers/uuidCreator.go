package experthelpers

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func GenerateTransactionUUID() string {
	timeNow := time.Now().UTC()
	entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
	return id.String()
}
