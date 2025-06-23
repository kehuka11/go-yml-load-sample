package config

// yamlで記載された設定を読み込み、Goで使用できる形で保持する
type TestConfig struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Driver   string `yaml:"driver"`
	} `yaml:"database"`

	Server struct {
		Port    int `yaml:"port"`
		Timeout struct {
			Api int `yaml:"api"`
			Db  int `yaml:"db"`
		} `yaml:"timeout"`
	} `yaml:"server"`

	Test struct {
		Check  bool `yaml:"check"`
		Normal bool `yaml:"normal"`
	} `yaml:"test"`
}
