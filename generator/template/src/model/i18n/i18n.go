package i18n

import (
	. "Vectra/src/model/service"
	"Vectra/src/model/storage"
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type I18n struct {
	dic map[string]map[string]string
	mu  sync.RWMutex
}

var instance *I18n
var once sync.Once

func GetInstance() *I18n {
	once.Do(func() {
		instance = &I18n{}
		instance.dic = make(map[string]map[string]string)
	})
	return instance
}

func (i *I18n) SetUp(langs ...string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, lang := range langs {
		path := filepath.Join(storage.I18nDirPath, lang)
		err := i.loadData(path, lang, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *I18n) loadData(path string, lang string, prefix string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		key := entry.Name()
		if entry.IsDir() {
			err := i.loadData(filepath.Join(path, key), lang, prefix+key+".")
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(key, ".ini") {
			if _, ok := i.dic[lang]; !ok {
				i.dic[lang] = make(map[string]string)
			}
			fullKey := prefix + strings.TrimSuffix(key, ".ini")

			data, err := os.ReadFile(filepath.Join(path, key))
			if err != nil {
				return err
			}

			cfg, _ := ini.LoadSources(ini.LoadOptions{}, data)
			for _, k := range cfg.Section("").Keys() {
				i.dic[lang][fullKey+"."+k.Name()] = k.Value()
			}
		}
	}

	return nil
}

// Get is a method of the I18n type that allows for dynamic string localization.
// It retrieves the translation string associated with the provided key.
// The method also supports pluralization of the translation based on the first argument in args,
// when it's an integer and different from 1 (which stands for singular form).
//
// Parameters:
//
// - key: The key associated with the translation string in the dictionary.
//
// - args: The optional arguments which can be used for string formatting and pluralization.
//
// Returns:
//
// - The formatted translation string associated with the key if it exists, otherwise "Key not found".
func (i *I18n) Get(key string, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if len(args) > 0 {
		// If there are parameters and the first one is an integer,
		// check its value for pluralization.
		if count, ok := args[0].(int); ok {
			if count != 1 {
				key = fmt.Sprintf("%s_plural", key)
			} else {
				key = fmt.Sprintf("%s_singular", key)
			}
		}
	}

	if val, ok := i.dic[GetApiV1().GetStore().Config.CurrentLang][key]; ok {
		sprintf := fmt.Sprintf(val, args...)
		return sprintf
	} else {
		return "Key not found"
	}
}
