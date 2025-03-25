package conf

import (
	"github.com/stretchr/testify/assert"
	"sportNews/conf/database"
	"testing"
)

func TestNew(t *testing.T) {
	conf, err := New("app.conf.example.yaml")

	c := &Conf{
		Project: "kingler",
		App: App{
			Addr:         "9011",
			Mode:         "debug",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
		},
		DB:     newDBConf("kingler"),
		Assets: "http://assets.localhost",
	}

	assert.Nil(t, err)
	assert.Equal(t, c, conf)
}

func newDBConf(dbName string) database.Config {
	return database.Config{
		Master: database.ConfigNode{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "1234",
		},
		Slave: database.ConfigNode{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "1234",
		},
		Database:        dbName,
		Timezone:        "UTC",
		DialTimeout:     "10s",
		ReadTimeout:     "30s",
		WriteTimeout:    "60s",
		ConnMaxLifetime: 0,
		MaxIdleConns:    2,
		MaxOpenConns:    0,
	}
}
