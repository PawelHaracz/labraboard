package labraboard

type Config struct {
	ConnectionString string `yaml:"connectionString" env:"CONNECTION_STRING" env-description:"Connection string to the database" env-required:"true"`
	HttpPort         int    `yaml:"httpPort" env:"HTTP_PORT" env-description:"HTTP port to serve the application" env-default:"8080"`
	RedisHost        string `yaml:"redisHost" env:"REDIS_HOST" env-description:"Redis host" env-default:"localhost"`
	RedisPort        int    `yaml:"redisPort" env:"REDIS_PORT" env-description:"Redis port" env-default:"63"`
	RedisPassword    string `yaml:"redisPassword" env:"REDIS_PASSWORD" env-description:"Redis password" env-default:""`
	RedisDB          int    `yaml:"redisDB" env:"REDIS_DB" env-description:"Redis database" env-default:"0"`
}
