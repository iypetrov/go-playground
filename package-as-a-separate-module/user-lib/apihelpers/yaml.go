package apihelpers

import (
	go_yaml "github.com/goccy/go-yaml"
	"github.com/iypetrov/user-lib/apis"
)

func ToYaml(u apis.User) string {
	bytes, err := go_yaml.Marshal(u)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
