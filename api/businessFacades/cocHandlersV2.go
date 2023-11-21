package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func COCTransferRequestInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var cocstatus model.COCState
	err := json.NewDecoder(r.Body).Decode(&cocstatus)
	if err != nil {
		logrus.Error("Error decoding data from payload : ", err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.CommonBadResponse{
			Message: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	if cocstatus.COCStatus != model.COC_TRANSFER_ENABLED {
		logrus.Error("Invalid status:  ", err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.CommonBadResponse{
			Message: "Invalid COC status",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	repo := dao.Connection{}
	saveErr := repo.InsertCOCStatus(cocstatus)
	if saveErr != nil {
		logrus.Error("failed to save COC transfer request", err)
		result := apiModel.CommonBadResponse{
			Message: "failed to save COC transfer request",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	result := apiModel.CommonSuccessMessage{
		Message: "Success"}
	json.NewEncoder(w).Encode(result)
	return
}
