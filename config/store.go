package config

type Store struct {
	MessageExpired int
	Default        string
	Redis          StoreRedis
}

type StoreRedis struct {
	// redis config
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
