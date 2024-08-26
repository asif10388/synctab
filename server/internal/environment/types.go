package environment

const (
	ApiServiceDir             = ""
	ApiServiceConfDir         = ApiServiceDir + "/conf"
	ApiServiceDevEnvFile      = ApiServiceConfDir + "/apiservice-dev.env"
	ApiServiceProdEnvFile     = ApiServiceConfDir + "/apiservice-prod.env"
	ApiServiceEnvDefaultsFile = ApiServiceConfDir + "/apiservice-defaults.env"
	ApiServiceSqlDir          = ApiServiceDir + "/sql"
	ApiServiceSqlPrimaryDir   = ApiServiceSqlDir + "/primary"
)

type Environment struct {
	data  map[string]string
	etype string
}
