package expertformula

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	equationbuilding "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula/equationBuilding"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
StellarExpertFormulBuilder
des- This method build stellar trasactions for expert formula

	 steps
		* map the formulaId and retrive the mapped id
		* build memo for the trasacions
		* map the experId and retive mapped id
		* build formula identity manageData opration
		* build auther identity manage data opration
		* loop through the formulaArray to see build the field definitions and build relevenat manage data oprations
		* get the execution template from fcld and build relevenat manage data oprations
		* load stellar account,build and sing the XDR
		* put XDR to stellar blockchain
*/
func StellarExpertFormulBuilder(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int,variableCount int) {
	w.Header().Set("Content-Type", "application/json")
	formulaArray := formulaJSON.MetricExpertFormula.Formula // formula array sent by the backend                               // formula array sent by the backend
	var manageDataOpArray []txnbuild.ManageData             // manageDataOpArray all manage data append to to this array
	expertIDMap := model.ExpertIDMap{}
	var formStat string
	// var startTransactionTime, endTransactionTime time.Time
	var expertMapID uint64
	var memo0, memo1 string
	object := dao.Connection{}
	formulaStatusDetails, errWhenGettingStatus := object.GetFormulaStatus(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingStatus != nil {
		logrus.Error("An error occured when getting formula status ", errWhenGettingStatus)
	}
	if formulaStatusDetails == nil {
		formStat = ""
	}
	if formulaStatusDetails != nil {
		formulaMapDet := formulaStatusDetails.(model.FormulaStore)
		formStat = formulaMapDet.Status
		logrus.Info("Status recorded : ", formStat)
	}

	if formStat == "QUEUE" {
		// ask user to try again
		logrus.Info("Requested formula is in the queue, please try again")
		commons.JSONErrorReturn(w, r, formStat, 400, "Requested formual is in the queue, please try again")
		return
	} else if formStat == "SUCCESS" {
		logrus.Info("Formula is already recorded in the blockchain and the gateway DB")
		// response indicating that formula is already recorded
		commons.JSONErrorReturn(w, r, formStat, 400, "Formula is already recorded in the blockchain and the gateway DB")
		return
	} else if formStat == "FAILED" || formStat == "" {
		logrus.Info("Requested formula id status is failed or a new binding request")
		// save expert formula in the database
		expertFormulaBuilder := model.FormulaStore{
			FormulaID:           formulaJSON.MetricExpertFormula.ID,
			MetricExpertFormula: formulaJSON.MetricExpertFormula,
			User:                formulaJSON.User,
			VariableCount:       fieldCount,
			Timestamp:           time.Now().String(),
			Status:              "FAILED",
		}
		// checked whether given formulaID already in the database or not
		formulaMap, err := object.GetExpertFormulaCount(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			logrus.Info(err)
		}
		// if formulA already in Database, not allowed to  build expert formula to that ID
		if formulaMap.(int64) != 0 {
			commons.JSONErrorReturn(w, r, "Formula Id is in gateway datastore", http.StatusBadRequest, "Duplicate formula IDs not allowed ")
			return
		}
		// if not,  retrived the current latest sequence number for formulaID
		dataFormulaID, err := object.GetNextSequenceValue("FORMULAID")
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "GetNextSequenceValue for formula Id was failed ")
			return
		}
		expertFormula := ExpertFormula{}
		// checked whether given ExpertID already in the database or not
		expertMapdata, err := object.GetExpertMapID(formulaJSON.User.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Unable to connect to gateway datastore ")
		}
		// if not,  retrived the current latest sequence number for expertID , map the expertID with incrementing interger
		if expertMapdata == nil {
			data, err := object.GetNextSequenceValue("EXPERTID")
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Mapping expert ID failed ")
				return
			}
			expertIDMap = model.ExpertIDMap{
				ExpertID:  formulaJSON.User.ID,
				ExpertPK:  formulaJSON.User.Publickey,
				MapID:     data.SequenceValue,
				FormulaID: formulaJSON.MetricExpertFormula.ID,
			}
			err1 := object.InsertExpertIDMap(expertIDMap)
			if err1 != nil {
				commons.JSONErrorReturn(w, r, err1.Error(), http.StatusInternalServerError, "Insert to ExpertIDMap was failed ")
				return
			}
			expertMapID = data.SequenceValue
		} else {
			expertMap := expertMapdata.(model.ExpertIDMap)
			expertMapID = expertMap.MapID
		}
		// formula identity operation
		formulaIdentityBuilder, errInFormulaIdentity := expertFormula.BuildFormulaIdentity(expertMapID, formulaJSON.MetricExpertFormula.Name)
		if errInFormulaIdentity != nil {
			commons.JSONErrorReturn(w, r, errInFormulaIdentity.Error(), http.StatusInternalServerError, "An error occured when building formula identity ")
			return
		}
		// append to the manage data array
		manageDataOpArray = append(manageDataOpArray, formulaIdentityBuilder)
		// author details opreation
		authorDetailsBuilder, errInAuthorBuilder := expertFormula.BuildPublicManageData(formulaJSON.User.Publickey)
		if errInAuthorBuilder != nil {
			commons.JSONErrorReturn(w, r, errInAuthorBuilder.Error(), http.StatusInternalServerError, "An error occured when building author identity ")
			return
		}
		// append to the manage data array
		manageDataOpArray = append(manageDataOpArray, authorDetailsBuilder)
		// loop through the formulaArray to see build the field definitions
		for i := 0; i < len(formulaArray); i++ {
			if formulaArray[i].Type == "VARIABLE" {
				// excute the variable builder
				variableBuilder, _, err := expertFormula.BuildVariableDefinitionManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Variable was failed ")
					return
				}
				// append to the manage data array
				manageDataOpArray = append(manageDataOpArray, variableBuilder)
			} else if formulaArray[i].Type == "REFERREDCONSTANT" {
				// execute the referred constant builder
				referredConstant, _, err := expertFormula.BuildReferredConstantManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Referred Constant was failed ")
					return
				}
				// url builder
				urlBuilder, err := expertFormula.BuildReference(formulaArray[i].MetricReference.Reference)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Referred Constant Reference was failed ")
					return
				}
				// append to the manage data array
				manageDataOpArray = append(manageDataOpArray, referredConstant)
				manageDataOpArray = append(manageDataOpArray, urlBuilder)
			} else if formulaArray[i].Type == "SEMANTICCONSTANT" {
				// execute the semantic constant builder
				sematicConstant, _, err := expertFormula.BuildSemanticConstantManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Semantic Constant was failed ")
					return
				}
				// value builder
				valueBuilder, err := expertFormula.BuildSemanticValue(formulaArray[i].Value)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Semantic Constant Value was failed ")
					return
				}
				// append to the manage data array
				manageDataOpArray = append(manageDataOpArray, sematicConstant)
				manageDataOpArray = append(manageDataOpArray, valueBuilder)
			}
		}

		/* logic section of the expert formula builder

		   * BuildExecutionTemplateByQuery() method will return the execution template that returns from the FCL
		   * if the lst_commands in the returned execution template is not empty
		   	-> Type 1 execution template(Start variable followed by a list of commands) - returns an array of manage data operations
		     else
		   	-> Type 2 execution template(Entity) - returns a single manage data operation
		*/
		executionTemplate, errInGettingExecutionTemplate := BuildExecutionTemplateByQuery(formulaJSON.MetricExpertFormula.FormulaAsQuery)
		if errInGettingExecutionTemplate != nil {
			commons.JSONErrorReturn(w, r, errInGettingExecutionTemplate.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
			return
		}
		expertFormulaBuilder.ExecutionTemplate = executionTemplate
		if executionTemplate.Lst_Commands != nil {
			manageDataOp, errTemplate1Builder := equationbuilding.Type1TemplateBuilder(formulaJSON.MetricExpertFormula.ID, executionTemplate)
			if errTemplate1Builder != nil {
				commons.JSONErrorReturn(w, r, errTemplate1Builder.Error(), http.StatusInternalServerError, "Error in building execution template type 1 failed ")
				return
			}
			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, manageDataOp...)
		} else {
			template1Builder, errInTemplate1Builder := equationbuilding.Type2TemplateBuilder(formulaJSON.MetricExpertFormula.ID, executionTemplate)
			if errInTemplate1Builder != nil {
				commons.JSONErrorReturn(w, r, errInTemplate1Builder.Error(), http.StatusInternalServerError, "Error in building execution template type 2 failed ")
				return
			}
			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, template1Builder)
		}

		// split the manage data array into two parts
		manageData2dArray := commons.ChunkSlice(manageDataOpArray, 25)
		for i, manadataOperationArray := range manageData2dArray {
			if i == 0 {

				// build memo0 send the transaction
				memo0, _, err = expertFormula.BuildMemo(0, uint32(fieldCount), dataFormulaID.SequenceValue)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Hex conversion issue in building memo ")
					return
				}
				if len(memo0) != 28 {
					commons.JSONErrorReturn(w, r, "Memo length error(expertFormulaBuilder) ", http.StatusInternalServerError, memo0)
					return
				}
			}
			memo := memo0
			if i != 0 {
				// here for insted of no of values we pass the current index of the manadataOperationArray array
				memo1, _, err = expertFormula.BuildMemo(1, uint32(fieldCount), dataFormulaID.SequenceValue)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Hex conversion issue in building memo")
					return
				}
				if len(memo1) != 28 {
					commons.JSONErrorReturn(w, r, "Memo length error(expertFormulaBuilder) ", http.StatusInternalServerError, memo1)
					return
				}
				memo = memo1
			}
			expertFormulaBuilder.FormulaMapID = dataFormulaID.SequenceValue
			expertFormulaBuilder.NoOfManageDataInTxn = len(manadataOperationArray)
			expertFormulaBuilder.TotalNoOfManageData = len(manageDataOpArray)

			timeNow := time.Now().UTC()
			entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
			id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
			logrus.Info("TXN UUID : ", id)
			expertFormulaBuilder.TxnUUID = id.String()

			buildMetricBind := model.SendToQueue{
				ExpertFormula: expertFormulaBuilder,
				Type:          "EXPERTFORMULA",
				User:          formulaJSON.User,
				Status:        "QUEUE",
				Operations:    manadataOperationArray,
				Memo:          []byte(memo),
			}
			err := services.SendToQueue(buildMetricBind)
			if err != nil {
				expertFormulaBuilder.ErrorMessage = err.Error()
				_, errResult := object.InsertExpertFormula(expertFormulaBuilder)
				if errResult != nil {
					logrus.Error("Error while inserting the Expoert formula into DB: ", errResult)
				}
				logrus.Error("Error when submitting managedata to queue (METRIC BINDING) ", err)
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Error{Code: http.StatusInternalServerError, Message: "Error when submitting managedata to queue (METRIC BINDING) " + err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
			logrus.Info("Expert formula request sent to queue ")
			_, errResult := object.InsertExpertFormula(expertFormulaBuilder)
			if errResult != nil {
				logrus.Error("Error while inserting the metric expoert formula into DB: ", errResult)
			}
			formulaIDMap := model.FormulaIDMap{
				FormulaID:     formulaJSON.MetricExpertFormula.ID,
				MapID:         dataFormulaID.SequenceValue,
				VariableCount: variableCount,
				FieldCount:    fieldCount,
			}
			// map the formulaID with incrementing Integer put those object to blockchain
			err1 := object.InsertFormulaIDMap(formulaIDMap)
			if err1 != nil {
				logrus.Error("Inserting formula to the export formula map was failed " + err1.Error())
			}
		}
		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusOK,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Expert formula request sent to queue",
		}
		json.NewEncoder(w).Encode(response)
		return

	} else {
		logrus.Info("Fomrula bind status is invalid : ", formStat)
		commons.JSONErrorReturn(w, r, formStat, 504, "Formula bind status is invalid, status : ")
		return
	}
}
