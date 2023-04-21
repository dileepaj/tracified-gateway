package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//var redisClient *redis.Client

func BuyHandlerLock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var res model.PendingNFTS
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&res)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("This is the start", res)
	// Try to acquire a distributed lock for the item

	err = services.LockRequest(res)
	if err != nil {
		logrus.Error("Failed to aquire lock ", err)
	}

}
