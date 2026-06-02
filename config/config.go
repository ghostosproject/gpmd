package config

type Config struct {
	Version        string      `json:"version" env:"VERSION" envDefault:"2.1.1"`
	Ports          ConfigPorts `json:"ports" envPrefix:"PORT_"`
	PrivateKeyPath string      `json:"private" env:"PRIVATE_KEY_PATH" envDefault:"./private.key"`
	Peers          []string    `json:"peers"`
}

type ConfigPorts struct {
	P2P        string `json:"p2p" env:"P2P" envDefault:"8282"`
	Background string `json:"background" env:"BACKGROUND" envDefault:"1111"`
	GVM        string `json:"gvm" env:"GVM" envDefault:"5989"`
}
