package conf

import (
	"github.com/stretchr/testify/assert"
	"sportNews/conf/database"
	"testing"
)

// 測試 api-server APP config是否符合格式
func TestAppConfNew(t *testing.T) {
	conf, err := New("app.conf.example.yaml")

	c := &Conf{
		App: App{
			Name:         "api-server",
			Addr:         "9011",
			Debug:        true,
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
		},
		DB:     newDBConf("sport_news"),
		Assets: "http://assets.localhost",
	}

	assert.Nil(t, err)
	assert.Equal(t, c, conf)
}

// 測試 crawler-process APP config是否符合格式
func TestProcessConfNew(t *testing.T) {
	conf, err := New("process.conf.example.yaml")

	c := &Conf{
		App: App{
			Name:  "crawler-process",
			Debug: true,
			Process: Process{
				News:        7200,
				Ranking:     43200,
				SyncPicture: 1800,
			},
		},
		DB:     newDBConf("sport_news"),
		Assets: "",
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
