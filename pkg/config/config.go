package config

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func Load(c interface{}, container string, fileNames ...string) error {
	env := os.Getenv("ENV")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetConfigType("yaml")

	for _, fileName := range fileNames {
		viper.SetConfigName(fileName)
	}

	viper.AddConfigPath("./" + container + "/")

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			panic(err)
		}
	}
	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name := fileName2 + "-" + env2
			viper.SetConfigName(name)
			viper.MergeInConfig()
		}
	}
	er3 := BindEnv(c)
	if er3 != nil {
		return er3
	}
	er4 := viper.Unmarshal(c)
	return er4
}

func BindEnv(conf interface{}, parts ...string) error {
	ifv := reflect.Indirect(reflect.ValueOf(conf))
	ift := reflect.TypeOf(ifv)
	num := min(ift.NumField(), ifv.NumField())
	for i := 0; i < num; i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			return BindEnv(v.Interface(), append(parts, tv)...)
		default:
			return viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
	return nil
}

func min(n1 int, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
