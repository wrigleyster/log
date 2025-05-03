package config

import (
	"github.com/wrigleyster/gorm"
	"github.com/wrigleyster/gorm/sqlite"
	"sync"
)

type Configuration struct {
	db gorm.DataSource
}

var config *Configuration
var mutex = &sync.Mutex{}

func instance() *Configuration {
	if config == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if config == nil {
			config = &Configuration{
				db: sqlite.New(GetDbName()),
			}
		}
	}
	return config
}

func GetDbName() string {
	return "sqlite.db"
}

func GetDb() gorm.DataSource {
	return instance().db
}
