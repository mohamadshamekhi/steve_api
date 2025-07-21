package mellipayamak

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	From     string `koanf:"from"`
}

type MelliPayamakAdapter struct {
	cfg Config
}

func (m MelliPayamakAdapter) SendPatternSMS(patternID uint, params []string, to string) error {
	panic("implement me")
}
