package polygonexpertformula

import (
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	experthelpers "github.com/dileepaj/tracified-gateway/protocols/polygon/polygonCodeGenerator/polygonExpertFormula/expertHelpers"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func PolygonExpertFormulaContractGenerator(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	object := dao.Connection{}
	var deployStatus int
	// reqType := "POLYGONEXPERT"
	logger := utilities.NewCustomLogger()

	formulaDetails, errWhenGettingFormulaDetails := object.GetPolygonFormulaStatus(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if errWhenGettingFormulaDetails != nil {
		logger.LogWriter("An error occurred when getting formula status, ERROR : "+errWhenGettingFormulaDetails.Error(), constants.ERROR)
	}
	if formulaDetails == nil {
		deployStatus = 0
	}
	if formulaDetails != nil {
		deployStatus = formulaDetails.(model.EthereumExpertFormula).Status
		logger.LogWriter("Polygon formula contract deploy status : "+strconv.FormatInt(int64(deployStatus), 10), constants.INFO)
	}

	if deployStatus != 0 || deployStatus != 119 {
		//handle Queue, Success, invalid status
		experthelpers.SuccessOrQueueResponse(w, r, formulaJSON, deployStatus)
	} else {
		if deployStatus == 119 {
			logger.LogWriter("Requested formula is in the failed status, trying to redeploy", constants.INFO)
		} else {
			logger.LogWriter("New expert formula request, initiating new deployment", constants.INFO)
		}

		//create expert formula
		formulaObj := experthelpers.BuildExpertObject(formulaJSON.MetricExpertFormula.ID, formulaJSON.MetricExpertFormula.Name, formulaJSON.MetricExpertFormula, fieldCount, formulaJSON.Verify)

		if deployStatus == 0 {
			transactionUuid := experthelpers.GenerateTransactionUUID()
			formulaObj.TransactionUUID = transactionUuid
		}
	}

}
