package model

type BuildFormula struct {
	FormulaID       string
	FormulaType     string
	Name            string
	FieldCount      int32
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

type FormulaItemRequest struct {
	ID                string
	Value             float64
	Name              string
	Description       string
	Key               string
	Type              string
	MeasurementUnit   string
	MetricReferenceId string
	Precision         int64
	MetricReference   MetricReference
}

type Expert struct {
	ExpertID string
	ExpertPK string
}

type FormulaIDMap struct {
	FormulaID string
	MapID     int64
}

type ExpertIDMap struct {
	ExpertID  string
	ExpertPK  string
	MapID     int64
	FormulaID string
}

type Counters struct {
	ID            string
	SequenceValue int64
}

type ValueIDMap struct {
	ValueId   string
	ValueType string
	ValueName string
	MapID     int64
}

type UnitIDMap struct {
	Unit  string
	MapID int64
}

type FormulaBuildingRequest struct {
	ID              string
	Name            string
	Metric          FormulaItem
	Formula         []FormulaItemRequest
	FormulaAsString string
	FormulaAsQuery  string
	Activity        Activity
	Expert          Expert
	Active          bool
	CreatedAt       string
	UpdatedAt       string
}

type MetricReference struct {
	Name            string
	MeasurementUnit string
	Url             string
}
