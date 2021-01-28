package infra

type SearchEngine interface {
	InitializeIndex(name string)
	CreateMany(v interface{}) error
	UpdateMany(v interface{}) error
	DeleteMany(oid []string) error
}

type SearchEngineOperation string

const (
	SearchEngineOperationCreate SearchEngineOperation = "create"
	SearchEngineOperationUpdate SearchEngineOperation = "update"
	SearchEngineOperationDelete SearchEngineOperation = "delete"
)
