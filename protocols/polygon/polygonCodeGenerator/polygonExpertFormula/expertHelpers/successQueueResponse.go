package experthelpers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
)

//Handle Queue, Success and Invalid response
func SuccessOrQueueResponse(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, deployStatus int) {
	logger := utilities.NewCustomLogger()

	if deployStatus == 118 { // SUCCESS
		logger.LogWriter("Contract for formula "+formulaJSON.MetricExpertFormula.Name+" has been added to the blockchain", constants.INFO)
		w.WriteHeader(http.StatusBadRequest)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusBadRequest,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Requested formula is in the blockchain :  Status : " + strconv.Itoa(deployStatus),
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if deployStatus == 116 { // QUEUE
		logger.LogWriter("Requested formula is in the queue, please try again", constants.INFO)
		w.WriteHeader(http.StatusBadRequest)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusBadRequest,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Requested formula is in the queue :  Status : " + strconv.Itoa(deployStatus),
		}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		logger.LogWriter("Invalid formula status "+strconv.Itoa(deployStatus), constants.INFO)
		commons.JSONErrorReturn(w, r, strconv.Itoa(deployStatus), http.StatusInternalServerError, "Invalid formula status : ")
		return
	}
}
