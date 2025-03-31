package conf

import (
	"github.com/spf13/viper"
	"sportNews/conf/database"
	"strings"
)

type Conf struct {
	App    App             `mapstructure:"app"`
	DB     database.Config `mapstructure:"db"`
	Assets string          `mapstructure:"assets"`
}

type App struct {
	Name         string  `mapstructure:"name"`
	Addr         string  `mapstructure:"addr"`
	Debug        bool    `mapstructure:"debug"`
	Process      Process `mapstructure:"process"`
	ReadTimeout  string  `mapstructure:"read_timeout"`
	WriteTimeout string  `mapstructure:"write_timeout"`
}

type Process struct {
	News    int `mapstructure:"news"`
	Ranking int `mapstructure:"ranking"`
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func New(file string) (conf *Conf, err error) {
	conf = &Conf{}
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
