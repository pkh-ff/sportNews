package conf

import (
	"github.com/spf13/viper"
	"sportNews/conf/database"
	"strings"
)

type Conf struct {
	Project string          `mapstructure:"project"`
	App     App             `mapstructure:"app"`
	DB      database.Config `mapstructure:"db"`
	Assets  string          `mapstructure:"assets"`
}

type App struct {
	Addr         string `mapstructure:"addr"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
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
