package experthelpers

import (
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
)

func BuildExpertObject(formulaID string, formulaName string, metricExpertFormula model.ExpertFormula, fieldCount int, verify model.Verify) model.EthereumExpertFormula {

	polygonFormulaObject := model.EthereumExpertFormula{
		FormulaID:           formulaID,
		FormulaName:         formulaName,
		MetricExpertFormula: metricExpertFormula,
		VariableCount:       int32(fieldCount),
		ContractAddress:     "",
		Timestamp:           time.Now().String(),
		TransactionHash:     "",
		TransactionCost:     "",
		TransactionUUID:     "",
		GOstring:            "",
		TransactionSender:   commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
		Verify:              verify,
		ErrorMessage:        "",
		ActualStatus:        101, // SMART_CONTRACT_GENERATION_STARTED
	}

	return polygonFormulaObject
}
