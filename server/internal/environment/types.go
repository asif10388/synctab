package environment

const (
	ApiServiceDir             = ""
	ApiServiceConfDir         = ApiServiceDir + "/conf"
	ApiServiceSqlDir          = "/apiserver/model/sql/synctabdb"
	ApiServiceOverrideEnvFile = ApiServiceConfDir + "/apiservice-override.env"
	ApiServiceEnvDefaultsFile = ApiServiceConfDir + "/apiservice-defaults.env"
)

type Environment struct {
	data  map[string]string
	etype string
}
