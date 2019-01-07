package builder

import (
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/dileepaj/tracified-gateway/dao"
	"fmt"
	"strings"
	"github.com/stellar/go/xdr"
	// "github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/model"
)

// type InsertData struct{}


func XDRSubmitter(TDP []model.TransactionCollectionBody) bool {
	object := dao.Connection{}
	var copy model.TransactionCollectionBody
	for i := 0; i < len(TDP); i++ {
		TDP[i].Status = "Pending"
		var txe xdr.Transaction
		err:= xdr.SafeUnmarshalBase64(TDP[i].XDR, &txe)
		if err != nil {
			fmt.Println(err)
		}

		TDP[i].PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		TDP[i].TxnType = TxnType
		TDP[i].Status = "pending"

		copy=TDP[i]
		err1 := object.InsertTransaction(TDP[i])
		if err1 != nil {
			TDP[i].Status = "failed"
		}

	}
	for i := 0; i < len(TDP); i++ {
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TDP[i].XDR}

		response := display.SubmitXDR()
		if response.Error.Code == 503 {
			TDP[i].Status = "pending"
		} else {
			TDP[i].TxnHash = response.TXNID

			upd := model.TransactionCollectionBody{TxnHash: response.TXNID, Status: "done"}
			err2 := object.UpdateTransaction(copy, upd)
			if err2 != nil {
				TDP[i].Status = "pending"
			} else {
				TDP[i].Status = "done"
			}
		}
	}

	return true
}
