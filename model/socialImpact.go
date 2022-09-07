package model

type BuildFormula struct {
	FormulaId       string
	FormulaType     string
	FieldCount      string
	Query           string
	FormulaAsString string
	TenantId        string
	Metric          FormulaItem
	Activity        Activity
	Expert          Expert
	Formula         []FormulaItem
	Active          bool
	CreatedAt       string
	UpdatedAt       string
}
type FormulaItem struct {
	Id           string
	Name         string
	Description  string
	value        float32
	VariableType string
	DataType     string
	URL          string
	Unit         string
}

type Expert struct {
	ExpertID string
	ExpertPK string
}

