package entity

type Column struct {
	Name          string // guaranteed uniqueness
	DataType      DataType
	Length        *int
	NotNull       bool
	DefaultValue  *string
	AutoIncrement bool // Indicates whether the column is automatically incremented
	Comment       *string
}

type Index struct {
	Name    string
	Type    IndexType
	Columns []string
}

type TableOption struct {
	Collate       *string
	AutoIncrement *int64 // Set the auto-increment initial value of the table
	Comment       *string
}

type Table struct {
	Name      string // guaranteed uniqueness
	Columns   []*Column
	Indexes   []*Index
	Options   *TableOption
	CreatedAt int64
	UpdatedAt int64
}

type ResultSet struct {
	Columns      []string
	Rows         []map[string]interface{}
	AffectedRows int64
}
