package app

type Config struct {
  Name        string
	Deployment  string
	Credentials CredentialsConfig
	Server      ServerConfig
  Broker      BrokerConfig
}

type CredentialsConfig struct {
	Key    string
	Secret string
}

type ServerConfig struct {
	Addr string
	Port int
}

type BrokerConfig struct {
	Addr string
	Port int
  Topics TopicsConfig
}

type TopicsConfig struct {
	Request string
	Response string
}
