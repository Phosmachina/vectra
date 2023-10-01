package main

import (
	"Vectra/src/model/i18n"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
)

var (
	dic map[string]map[string]string
)

func main() {

	dic = make(map[string]map[string]string)

	langs := []string{"fr", "en"}

	for _, lang := range langs {
		path := filepath.Join("i18n", lang)
		_ = loadData(path, lang, "")
	}

	i18n.GenerateCode(dic["en"])
}

func loadData(path string, lang string, prefix string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		key := entry.Name()
		if entry.IsDir() {
			err := loadData(filepath.Join(path, key), lang, prefix+key+".")
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(key, ".csv") {
			if _, ok := dic[lang]; !ok {
				dic[lang] = make(map[string]string)
			}
			fullKey := prefix + strings.TrimSuffix(key, ".csv")

			data, err := os.ReadFile(filepath.Join(path, key))
			if err != nil {
				return err
			}

			reader := csv.NewReader(strings.NewReader(string(data)))
			records, err := reader.ReadAll()
			if err != nil {
				return err
			}

			for _, record := range records {
				if len(record) < 2 {
					continue
				}
				dic[lang][fullKey+"."+record[0]] = strings.ReplaceAll(record[1], "\\n", "\n")
			}
		}
	}

	return nil
}
