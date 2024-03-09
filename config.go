package labraboard

type Config struct {
	ConnectionString string `yaml:"connectionString" env:"CONNECTION_STRING" env-description:"Connection string to the database" env-required:"true"`
	HttpPort         int    `yaml:"httpPort" env:"HTTP_PORT" env-description:"HTTP port to serve the application" env-default:"8080"`
}
