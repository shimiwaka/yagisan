package schema

type Settings struct {
	Username     string `yaml:"db_username"`
	Pass         string `yaml:"db_pass"`
	Host         string `yaml:"db_host"`
	Port         int    `yaml:"db_port"`
	Name         string `yaml:"db_name"`
	MailHost     string `yaml:"mail_host"`
	MailAddress  string `yaml:"mail_address"`
	MailPassword string `yaml:"mail_password"`
	ServiceHost  string `yaml:"service_host"`
}