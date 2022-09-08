package model

type BuildFormula struct {
	FormulaID       string
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
	ID           string
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

type FormulaIDMap struct {
	FormulaID string
	MapID     int64
	ByteID    []byte
	BitString string
}

type Counters struct {
	ID            string
	SequenceValue int64
}