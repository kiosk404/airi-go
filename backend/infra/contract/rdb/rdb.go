package rdb

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/rdb/entity"
	"gorm.io/gorm"
)

type Option func(option OptionService)

type OptionService interface {
	WithMaster() Option            // WithMaster 强制读主库
	WithTransaction(tx RDB) Option // WithTransaction 使用一个已有的事务
	Debug() Option                 // Debug 启用调试模式
	WithDeleted() Option           // WithDeleted 返回软删的数据
	WithSelectForUpdate() Option   // WithSelectForUpdate 开启当前读
}

//go:generate mockgen -destination=mocks/db.go -package=mocks . Provider
type Provider interface {
	// NewSession 创建一个新的数据库会话
	NewSession(ctx context.Context, opts ...Option) RDB
	// Transaction 执行一个事务
	Transaction(ctx context.Context, fc func(tx RDB) error, opts ...Option) error
}

//go:generate mockgen -destination  ../../../internal/mock/infra/contract/rdb/rdb_mock.go  --package rdb  -source rdb.go
type RDB interface {
	CreateTable(ctx context.Context, req *CreateTableRequest) (*CreateTableResponse, error)
	AlterTable(ctx context.Context, req *AlterTableRequest) (*AlterTableResponse, error)
	DropTable(ctx context.Context, req *DropTableRequest) (*DropTableResponse, error)
	GetTable(ctx context.Context, req *GetTableRequest) (*GetTableResponse, error)

	InsertData(ctx context.Context, req *InsertDataRequest) (*InsertDataResponse, error)
	UpdateData(ctx context.Context, req *UpdateDataRequest) (*UpdateDataResponse, error)
	DeleteData(ctx context.Context, req *DeleteDataRequest) (*DeleteDataResponse, error)
	SelectData(ctx context.Context, req *SelectDataRequest) (*SelectDataResponse, error)
	UpsertData(ctx context.Context, req *UpsertDataRequest) (*UpsertDataResponse, error)

	ExecuteSQL(ctx context.Context, req *ExecuteSQLRequest) (*ExecuteSQLResponse, error)
	Transaction(ctx context.Context, fc func(tx RDB) error) error
	DB() *gorm.DB
}

// CreateTableRequest Create table request
type CreateTableRequest struct {
	Table *entity.Table
}

// CreateTableResponse Create table response
type CreateTableResponse struct {
	Table *entity.Table
}

// AlterTableOperation Modify table operation
type AlterTableOperation struct {
	Action       entity.AlterTableAction
	Column       *entity.Column
	OldName      *string
	Index        *entity.Index
	IndexName    *string
	NewTableName *string
}

// AlterTableRequest Modify table request
type AlterTableRequest struct {
	TableName  string
	Operations []*AlterTableOperation
}

// AlterTableResponse Modify table response
type AlterTableResponse struct {
	Table *entity.Table
}

// DropTableRequest drop table request
type DropTableRequest struct {
	TableName string
	IfExists  bool
}

// DropTableResponse Delete table response
type DropTableResponse struct {
	Success bool
}

// GetTableRequest Get table information request
type GetTableRequest struct {
	TableName string
}

// GetTableResponse Get table information response
type GetTableResponse struct {
	Table *entity.Table
}

// InsertDataRequest insert data request
type InsertDataRequest struct {
	TableName string
	Data      []map[string]interface{}
}

// InsertDataResponse
type InsertDataResponse struct {
	AffectedRows int64
}

// Condition defines query conditions
type Condition struct {
	Field    string
	Operator entity.Operator
	Value    interface{}
}

// ComplexCondition
type ComplexCondition struct {
	Conditions       []*Condition
	NestedConditions []*ComplexCondition // Conditions mutual exclusion example: WHERE (age > = 18 AND status = 'active') OR (age > = 21 AND status = 'pending')
	Operator         entity.LogicalOperator
}

// UpdateDataRequest
type UpdateDataRequest struct {
	TableName string
	Data      map[string]interface{}
	Where     *ComplexCondition
	Limit     *int
}

// UpdateDataResponse
type UpdateDataResponse struct {
	AffectedRows int64
}

// DeleteDataRequest Delete data request
type DeleteDataRequest struct {
	TableName string
	Where     *ComplexCondition
	Limit     *int
}

// DeleteDataResponse
type DeleteDataResponse struct {
	AffectedRows int64
}

type OrderBy struct {
	Field     string               // sort field
	Direction entity.SortDirection // sort direction
}

// SelectDataRequest query data request
type SelectDataRequest struct {
	TableName string
	Fields    []string // The field to query, if empty, query all fields
	Where     *ComplexCondition
	OrderBy   []*OrderBy // sort condition
	Limit     *int       // Limit the number of rows returned
	Offset    *int       // Offset
}

// SelectDataResponse
type SelectDataResponse struct {
	ResultSet *entity.ResultSet
	Total     int64 // Total number of eligible records (excluding paging)
}

type UpsertDataRequest struct {
	TableName string
	Data      []map[string]interface{} // Data to be updated or inserted
	Keys      []string                 // The column name used to identify a unique record, if empty, the primary key is used by default
}

type UpsertDataResponse struct {
	AffectedRows  int64 // Number of rows affected
	InsertedRows  int64 // Number of newly inserted rows
	UpdatedRows   int64 // updated rows
	UnchangedRows int64 // Constant number of rows (no rows inserted or updated)
}

// ExecuteSQLRequest Execute SQL request
type ExecuteSQLRequest struct {
	SQL    string
	Params []interface{} // For parameterized queries

	// SQLType indicates the type of SQL: parameterized or raw SQL. It takes effect if OperateType is 0.
	SQLType entity.SQLType
}

// ExecuteSQLResponse
type ExecuteSQLResponse struct {
	ResultSet *entity.ResultSet
}
