package bucharest_test

import (
	"os"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/stretchr/testify/assert"
)

var envFile = []byte(`FOO=BAR
FOZ="BAZ"
NUMZ=8
BOOLZ=true`)

func TestLoadEnvFile(t *testing.T) {
	path := "/tmp/.env"
	temp, err := os.Create(path)
	assert.NoError(t, err)
	assert.NotEmpty(t, temp)

	_, err = temp.Write(envFile)
	assert.NoError(t, err)

	err = temp.Close()
	assert.NoError(t, err)

	env, err := NewENV(path)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	all := env.All()
	assert.NotEmpty(t, all)

	foo := env.String("FOO")
	assert.Equal(t, "BAR", foo)

	foz := env.String("FOZ")
	assert.Equal(t, "BAZ", foz)

	numz := env.Int("NUMZ")
	assert.Equal(t, 8, numz)

	boolz := env.Bool("BOOLZ")
	assert.Equal(t, true, boolz)

	v := env.Viper()
	assert.NotEmpty(t, v)

	err = os.Remove(path)
	assert.NoError(t, err)
}
