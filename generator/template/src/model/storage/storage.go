package storage

import (
	. "github.com/phosmachina/FluentKV/reldb"
	. "github.com/phosmachina/FluentKV/reldb/impl"
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

// region DB Object

type Role struct {
	DBObject
	Name  string
	Level int
}

func NewRole(name string, level int) Role {
	role := Role{Name: name, Level: level}
	role.IObject = role
	return role
}

func (r Role) ToString() string  { return ToString(r) }
func (r Role) TableName() string { return NameOfStruct[Role]() }

type User struct {
	DBObject
	IsActivated bool
	Password    []byte

	Firstname string
	Lastname  string
	Email     string

	Sessions map[string]SessionItem
}

func NewUser(isActivated bool, password []byte, email string) User {
	user := User{IsActivated: isActivated, Password: password, Email: email}
	user.IObject = user
	return user
}

func (u User) ToString() string  { return ToString(u) }
func (u User) TableName() string { return NameOfStruct[User]() }

// endregion
