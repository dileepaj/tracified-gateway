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
	ID                string `json:"ID" bson:"id" validate:"required"`
	Value             float64
	Name              string
	Description       string
	Key               string
	Type              string `json:"Type" bson:"type" validate:"required"`
	MeasurementUnit   string
	MetricReferenceId string
	Precision         int64
	MetricReference   MetricReference
}

type Expert struct {
	ExpertID string `json:"ExpertID" bson:"expertID" validate:"required"`
	ExpertPK string `json:"ExpertPK" bson:"expertPK" validate:"required"`
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
	ID              string               `json:"ID" bson:"id" validate:"required"`
	Name            string               `json:"Name" bson:"name" validate:"required"`
	Metric          MetricItem           `json:"Metric" bson:"metric" validate:"required"`
	Formula         []FormulaItemRequest `json:"Formula" bson:"formula" validate:"required"`
	FormulaAsString string               `json:"FormulaAsString" bson:"formulaAsString" validate:"required"`
	FormulaAsQuery  string               `json:"FormulaAsQuery" bson:"formulaAsQuery" validate:"required"`
	Activity        Activity             `json:"Activity" bson:"activity" validate:"required"`
	Expert          Expert               `json:"Expert" bson:"expert" validate:"required"`
	CreatedAt       string               `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt       string               `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
}

type MetricReference struct {
	Name            string
	MeasurementUnit string
	Url             string
}

type MetricItem struct {
	ID          string `json:"ID" bson:"id" validate:"required"`
	Name        string `json:"Name" bson:"name" validate:"required"`
	Description string `json:"Description" bson:"description" validate:"required"`
}
