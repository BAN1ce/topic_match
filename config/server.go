package config

type Server struct {
	Port       int
	BrokerPort int
}

func (e Server) GetPort() int {
	return e.Port
}

func (e Server) GetBrokerPort() int {
	return e.BrokerPort
}
