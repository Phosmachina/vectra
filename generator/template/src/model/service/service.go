package service

import (
	. "Vectra/src/model"
	. "Vectra/src/model/storage"
	. "github.com/Phosmachina/FluentKV/reldb"
	"io"
	"log"
	"os"
	"sync"
)

var (
	lock = &sync.Mutex{}
)

type service struct {
	store            *Storage
	accessManager    *AccessManager
	firstLaunchToken string
}

func newService() *service {
	checkVolumeFiles()
	a := &service{}
	a.store = GetStorage()
	a.setupAccessManager()
	a.checkAppToken()
	return a
}

func (s *service) GetStore() *Storage {
	return s.store
}

func (s *service) GetAccessManager() *AccessManager {
	return s.accessManager
}

func (s *service) IsFirstLaunch() bool {
	// Check to see if a user has the "admin" role.
	return len(VisitWrp[Role, User](s.accessManager.DefaultRoles["admin"])) == 0
}

// checkAppToken if is the first launch, the initialization token will be generated and write in file.
func (s *service) checkAppToken() {
	if !s.IsFirstLaunch() {
		return
	}

	s.firstLaunchToken = GenerateToken(45)

	err := os.WriteFile(TokenFilePath, []byte(s.firstLaunchToken+"\n"), 0755)
	if err != nil {
		log.Fatal("Failed to write the app token in file.")
	}
}

// checkVolumeFiles checks the presence of all necessary files in folder that could be
// potentially mounted as Docker volume.
// If a file is missing, it is copied from the fallback folder.
func checkVolumeFiles() {
	_, err := os.Stat(ConfigFilePath)
	if err == nil {
		return
	}

	src, err := os.Open(FallbackDirPath + "configuration.yml")
	if err != nil {
		log.Fatal("Fail to find the fallback config file.")
	}
	defer src.Close()
	dest, err := os.OpenFile(ConfigFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("Fail to create the config file.")
		return
	}
	defer dest.Close()
	_, _ = io.Copy(dest, src)
}

// region Access management

// TODO use a robust RBAC like Casbin

// setupAccessManager create and fill value for AccessManager of this service.
//
// Firstly, roles are read in the configuration file. DefaultRoles are filled with it:
// it takes value saved in db or create a new one. After that, DefaultRoles it used to fill RulesRoutes and RulesTables.
func (s *service) setupAccessManager() {
	config := s.store.Config
	db := *s.store.DB

	s.accessManager = NewAccessManager()

	for name, level := range config.Roles {
		roleInDB := FindFirst(db, func(id string, role *Role) bool {
			return role.Name == name
		})
		if roleInDB == nil {
			role := NewRole()
			role.Name = name
			role.Level = level
			roleWrp := Insert[Role](db, role)
			s.accessManager.DefaultRoles[name] = roleWrp
		} else {
			s.accessManager.DefaultRoles[name] = roleInDB
		}
	}

	for _, rule := range config.AccessRules {
		switch rule.Target {
		case "route":
			s.accessManager.RulesRoutes[rule.Component] = s.accessManager.DefaultRoles[rule.Role]
		case "table":
			s.accessManager.RulesTables[rule.Component] = s.accessManager.DefaultRoles[rule.Role]
		}
	}
}

type AccessManager struct {
	DefaultRoles map[string]*ObjWrapper[Role]

	RulesRoutes map[string]*ObjWrapper[Role]
	RulesTables map[string]*ObjWrapper[Role]
}

func NewAccessManager() *AccessManager {
	manager := AccessManager{
		DefaultRoles: map[string]*ObjWrapper[Role]{},
		RulesRoutes:  map[string]*ObjWrapper[Role]{},
		RulesTables:  map[string]*ObjWrapper[Role]{},
	}
	return &manager
}

func (m AccessManager) CheckAccessForRoute(userId string, page string) bool {

	db := *GetApiV1().GetStore().DB

	var userRoleWrp *ObjWrapper[Role] = nil
	if userId == "" {
		userRoleWrp = GetApiV1().GetAccessManager().DefaultRoles["none"]
	} else {
		userRoleWrp = AllFromLink[User, Role](db, userId)[0]
	}

	return m.RulesRoutes[page] != nil && userRoleWrp.Value.Level >= m.RulesRoutes[page].Value.Level
}

func CheckAccessForTable[T IObject](userId string, idOfT string) bool {

	db := *GetApiV1().GetStore().DB
	m := GetApiV1().GetAccessManager()

	ttn := TableName[T]()
	userRoleWrp := AllFromLink[User, Role](db, userId)[0]

	if m.RulesTables[ttn] == nil {
		return false
	}

	if userRoleWrp.Value.Level > m.RulesTables[ttn].Value.Level {
		return true
	} else if userRoleWrp.Value.Level == m.RulesTables[ttn].Value.Level {
		return len(Visit[User, T](db, idOfT)) > 0
	}

	return false
}

// endregion
