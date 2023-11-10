package storage

import (
	. "github.com/Phosmachina/FluentKV/reldb"
	. "github.com/Phosmachina/FluentKV/reldb/impl"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	FallbackDirPath = "fallback"
	DataDirPath     = "data"
	BinDirPath      = filepath.Join(DataDirPath, "bin")
	DbDirPath       = filepath.Join(DataDirPath, "db")
	ConfigDirPath   = filepath.Join(DataDirPath, "config")
	I18nDirPath     = filepath.Join(DataDirPath, "i18n")
	ConfigFilePath  = filepath.Join(ConfigDirPath, "configuration.yml")
	TokenFilePath   = filepath.Join(ConfigDirPath, "init-token")

	lock                 = &sync.Mutex{}
	instance             *Storage
	CookieNameForSession = "session-id"
	CookieNameForCSRF    = "csrf-token"
	SessionKeyForUserId  = "user-id"
)

type Storage struct {
	DB     *IRelationalDB
	Config *configuration
}

func newStorage() *Storage {
	s := &Storage{}
	bdb, err := NewBadgerDB(DbDirPath)
	if err != nil {
		log.Fatal("Failed to open Badger database.")
	}

	GobRegistration()

	s.DB = &bdb
	s.Config = &configuration{
		Roles:       map[string]int{},
		AccessRules: []AccessRule{},
	}

	return s
}

func GetStorage() *Storage {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = newStorage()
			instance.LoadConfiguration()
		}
	}

	return instance
}

func (s Storage) LoadConfiguration() {
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, s.Config); err != nil {
		log.Fatal(err)
	}
	s.Config.CurrentLang = s.Config.DefaultLang
}

type SessionItem struct {
	LastViewed time.Time
	UA         string
}

type AccessRule struct {
	Target    string
	Component string
	Role      string
}
