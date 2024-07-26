package models

type ApplicationConfig struct {
	QiscusConfig struct {
		BaseUrl string `yaml:"baseurl"`
		AppCode string `yaml:"appcode"`
		Secret  string `yaml:"secret"`
		Sender  string `yaml:"sender"`
	} `yaml:"qiscus"`
	Evoucher struct {
		Url string `yaml:"url"`
	} `yaml:"evoucher"`
}
