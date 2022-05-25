package bucharest

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type ENV interface {
	All() Map
	Bool(key string) bool
	Int(key string) int
	String(key string) string
	Viper() *viper.Viper
}

type env struct{ preventDefaultPointer uuid.UUID }

func (e *env) All() Map {
	return viper.AllSettings()
}

func (e *env) Bool(key string) bool {
	return viper.GetBool(key)
}

func (e *env) Int(key string) int {
	return viper.GetInt(key)
}

func (e *env) String(key string) string {
	return viper.GetString(key)
}

func (e *env) Viper() *viper.Viper {
	return viper.GetViper()
}

func NewENV(filename string) (ENV, error) {
	viper.New()
	viper.SetConfigFile(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	return &env{preventDefaultPointer: uuid.New()}, viper.ReadInConfig()
}
