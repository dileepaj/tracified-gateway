package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/rabbitmq"
	"github.com/gorilla/mux"
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
	// Try to acquire a distributed lock for the item

	err = rabbitmq.LockRequest(res)
	if err != nil {
		logrus.Error("Failed to aquire lock ", err)
	}
	return
}

func RetrieveQueueData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetQueueData(vars["ImageBase64"], vars["blockchain"], vars["version"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.PendingNFTS)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "URL Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}
