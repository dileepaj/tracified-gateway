package businesslogic

import (
	"errors"
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/dao"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/stellar/go/xdr"
)

type AbstractXDR struct {
	TxnBody []model.TransactionCollectionBody
}

func (AP *AbstractXDR) GenesisAndDataXDRToCron() (error, int, string) {
	var task []model.Task // array to decide whether the actions are done
	connect := dao.Connection{}
	customLogger := utilities.NewCustomLogger()
	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		customLogger.LogWriter("TxnBody.XDR "+TxnBody.XDR, 2)
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

		_, err2, _ := connect.InsertTempOrphan(AP.TxnBody[i])
		if err2 != nil {
			customLogger.LogWriter("Error @ InsertSpecialToTempOrphan @SubmitSpecial "+err2.Error(), 3)
			task = append(task, model.Task{
				Identifier: AP.TxnBody[i].Identifier,
				Status:     false,
			})
		}

		if AP.TxnBody[i].TxnType == "0" {
			insertIdentifierMap(TxnBody.Identifier, AP.TxnBody[i].Identifier)
		}
	}
	unSubmittedTasks := commons.CheckTaskStatus(task)
	if len(unSubmittedTasks) == 0 {
		return nil, http.StatusOK, "Success"
	} else {
		tasksString := commons.TasksToString(unSubmittedTasks)
		return errors.New("Unsubmitted XDRs to TempOrphan: " + tasksString), http.StatusInternalServerError, ""
	}
}
