package constants

// routing
const (
	RootPrefix  = "/mydata"
	DocPath     = "documents"
	DocIdPath   = "docid"
	StartTsPath = "startTimestamp"
	EndTsPath   = "endTimestamp"
	TsPath      = "timestamp"
	KeyPath     = "key"
	ValuePath   = "value"
)

//
const (
	DefaultMaxRecordsReturn = 500
	DefaultCacheTimeMins    = 5
)

// DB
const (
	DbType         = "postgres"
	DbDSNBase      = DbType + "://"
	DbHostEnv      = "DB_HOST"
	DbUserEnv      = "DB_USER"
	DbPassEnv      = "DB_PASSWORD"
	DbSchemaEnv    = "DB_SCHEMA"
	DbPortEnv      = "DB_PORT"
	DbInsecureConn = "sslmode=disable"
)
