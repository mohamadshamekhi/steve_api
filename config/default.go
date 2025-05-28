package config

func Default() Config {
	return Config{
		HTTPServer: HTTPServer{
			Port: 7777,
		},
	}
}
