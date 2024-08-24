package environment

const (
	ApiServiceDir             = "/etc/apiservice"
	ApiServiceConfDir         = ApiServiceDir + "/conf"
	ApiServiceDevEnvFile      = ApiServiceConfDir + "/apiservice-dev.env"
	ApiServiceProdEnvFile     = ApiServiceConfDir + "/apiservice-prod.env"
	ApiServiceEnvDefaultsFile = ApiServiceConfDir + "/apiservice-defaults.env"
)

type Environment struct {
	data  map[string]string
	etype string
}
