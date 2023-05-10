package businesslogic

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/apiDemo/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/stellar/go/xdr"
)

type AbstractXDR struct {
	TxnBody []model.TransactionCollectionBody
}

func (AP *AbstractXDR) GenesisAndDataXDRToCron() (error, int, string) {
	var task []model.Task                                         // array to decide whether the actions are done
	task = append(task, model.Task{Identifier: "", Status: true}) // starting with a true for bipass
	connect := dao.Connection{}
	customLogger := utilities.NewCustomLogger()
	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		customLogger.LogWriter("Index ", 2)
		customLogger.LogWriter(i, 2)
		customLogger.LogWriter("TxnBody.XDR ", 2)
		customLogger.LogWriter(TxnBody.XDR, 2)

		// decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			customLogger.LogWriter("Error @ SafeUnmarshalBase64 @SubmitSpecial "+err.Error(), 3)
		}
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		stellarRetriever.MapXDROperations(&AP.TxnBody[i], txe.Operations)
		AP.TxnBody[i].Status = "pending"
		customLogger.LogWriter(AP.TxnBody[i].Identifier, 2)

		err2, err2Code := connect.InsertTempOrphan(AP.TxnBody[i])
		if err2 != nil {
			customLogger.LogWriter("Error @ InsertSpecialToTempOrphan @SubmitSpecial "+err2.Error(), 3)
			task = append(task, model.Task{Identifier: AP.TxnBody[i].Identifier, Status: false})
			// w.WriteHeader(http.StatusBadRequest)
			// response := apiModel.SubmitXDRSuccess{
			// 	Status: "Index[" + strconv.Itoa(i) + "] TXN: Scheduling Failed!",
			// }
			// json.NewEncoder(w).Encode(response)
			return errors.New("Index[" + strconv.Itoa(i) + "] TXN: Scheduling Failed!" + err2.Error()), err2Code, ""
		}

		if AP.TxnBody[i].TxnType == "0" {
			insertIdentifierMap(TxnBody.Identifier, AP.TxnBody[i].Identifier)
		}
	}
	unSubmittedTasks := checkTaskStatus(task)
	if len(unSubmittedTasks) == 0 {
		// w.WriteHeader(http.StatusOK)
		// result := apiModel.SubmitXDRSuccess{
		// 	Status: "Success",
		// }
		// json.NewEncoder(w).Encode(result)
		return nil, http.StatusOK, "Success"
	} else {
		tasksString := tasksToString(unSubmittedTasks)
		return errors.New("Unsubmitted Tasks: " + tasksString), http.StatusInternalServerError, ""
	}
}

func checkTaskStatus(tasks []model.Task) []model.Task {
	var failedTasks []model.Task
	for _, task := range tasks {
		if !task.Status {
			failedTasks = append(failedTasks, task)
		}
	}
	return failedTasks
}

func tasksToString(tasks []model.Task) string {
	var tasksString []string
	for _, task := range tasks {
		tasksString = append(tasksString, fmt.Sprintf("%s %v", task.Identifier, task.Status))
	}
	return fmt.Sprintf("[%s]", strings.Join(tasksString, ", "))
}
