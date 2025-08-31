package entity

type DataType string

const (
	TypeInt       DataType = "INT"
	TypeVarchar   DataType = "VARCHAR"
	TypeText      DataType = "TEXT"
	TypeBoolean   DataType = "BOOLEAN"
	TypeJson      DataType = "JSON"
	TypeTimestamp DataType = "TIMESTAMP"
	TypeFloat     DataType = "FLOAT"
	TypeBigInt    DataType = "BIGINT"
	TypeDouble    DataType = "DOUBLE"
)

type IndexType string

const (
	PrimaryKey IndexType = "PRIMARY KEY"
	UniqueKey  IndexType = "UNIQUE KEY"
	NormalKey  IndexType = "KEY"
)

// AlterTableAction defines the type of action to modify a table
type AlterTableAction string

const (
	AddColumn    AlterTableAction = "ADD COLUMN"
	DropColumn   AlterTableAction = "DROP COLUMN"
	ModifyColumn AlterTableAction = "MODIFY COLUMN"
	RenameColumn AlterTableAction = "RENAME COLUMN"
	AddIndex     AlterTableAction = "ADD INDEX"
)

type LogicalOperator string

const (
	AND LogicalOperator = "AND"
	OR  LogicalOperator = "OR"
)

type Operator string

const (
	OperatorEqual        Operator = "="
	OperatorNotEqual     Operator = "!="
	OperatorGreater      Operator = ">"
	OperatorGreaterEqual Operator = ">="
	OperatorLess         Operator = "<"
	OperatorLessEqual    Operator = "<="

	OperatorLike    Operator = "LIKE"
	OperatorNotLike Operator = "NOT LIKE"

	OperatorIn    Operator = "IN"
	OperatorNotIn Operator = "NOT IN"

	OperatorIsNull    Operator = "IS NULL"
	OperatorIsNotNull Operator = "IS NOT NULL"
)

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"  // ascending order
	SortDirectionDesc SortDirection = "DESC" // descending order
)

type SQLType int32

const (
	SQLType_Parameterized SQLType = 0
	SQLType_Raw           SQLType = 1 // Complete/raw SQL
)
