package models

type ApplicationConfig struct {
	Logging struct {
		Enabled bool   `yaml:"enabled"`
		Output  string `yaml:"output"`
	} `yaml:"logging"`
	QiscusConfig struct {
		BaseUrl string `yaml:"baseurl"`
		AppCode string `yaml:"appcode"`
		Secret  string `yaml:"secret"`
		Sender  string `yaml:"sender"`
	} `yaml:"qiscus"`
	Evoucher struct {
		Url string `yaml:"url"`
	} `yaml:"evoucher"`
	Kalista struct {
		Database struct {
			Server   string `yaml:"server"`
			Name     string `yaml:"name"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Port     int    `yaml:"port"`
		} `yaml:"database"`
	} `yaml:"kalista"`
	Binaries struct {
		Inkscape string `yaml:"inkscape"`
	} `yaml:"binaries"`
}
