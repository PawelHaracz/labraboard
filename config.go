package labraboard

type Config struct {
	ConnectionString string `yaml:"connectionString" env:"CONNECTION_STRING" env-description:"Connection string to the database" env-required:"true" json:"connectionString,omitempty"`
	HttpPort         int    `yaml:"httpPort" env:"HTTP_PORT" env-description:"HTTP port to serve the application" env-default:"8080" json:"httpPort,omitempty"`
	RedisHost        string `yaml:"redisHost" env:"REDIS_HOST" env-description:"Redis host" env-default:"localhost" json:"redisHost,omitempty"`
	RedisPort        int    `yaml:"redisPort" env:"REDIS_PORT" env-description:"Redis port" env-default:"6379" json:"redisPort,omitempty"`
	RedisPassword    string `yaml:"redisPassword" env:"REDIS_PASSWORD" env-description:"Redis password" env-default:"eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81" json:"redisPassword,omitempty"`
	RedisDB          int    `yaml:"redisDB" env:"REDIS_DB" env-description:"Redis database" env-default:"0" json:"redisDB,omitempty"`
	LogLevel         int8   `yaml:"logLevel" env:"LOG_LEVEL" env-description:"Redis database" env-default:"1" json:"logLevel,omitempty"`
	UsePrettyLogs    bool   `yaml:"usePrettyLogs" env:"USE_PRETTY_LOGS" env-description:"use pretty logs instead of json. Logs are pushed to stdout" env-default:"false" json:"usePrettyLogs,omitempty"`
	ServiceDiscovery string `yaml:"serviceDiscovery" env:"SERVICE_DISCOVERY" env-default:"http://localhost" json:"serviceDiscovery,omitempty"`
	FrontendPath     string `yaml:"frontendPath" env:"FRONTEND_PATH" env-default:"/app/client" json:"frontendPath,omitempty"`
}
