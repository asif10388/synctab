package environment

const (
	ApiServiceDir             = ""
	ApiServiceConfDir         = ApiServiceDir + "/conf"
	ApiServiceSqlDir          = "/apiserver/model/sql/synctabdb"
	ApiServiceDevEnvFile      = ApiServiceConfDir + "/apiservice-dev.env"
	ApiServiceProdEnvFile     = ApiServiceConfDir + "/apiservice-prod.env"
	ApiServiceEnvDefaultsFile = ApiServiceConfDir + "/apiservice-defaults.env"
)

type Environment struct {
	data  map[string]string
	etype string
}
