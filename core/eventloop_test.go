package core_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/stretchr/testify/assert"
)

func TestEventloopLoadStore(t *testing.T) {
	// file not found
	core.DataFile = "file-not-found.json"
	data := core.LoadState()

	assert.Equal(t, 0, len(*data.Devices))

	// permission denied
	core.DataFile = "/etc/shadow"
	data = core.LoadState()

	assert.Nil(t, data)

	// program refers to an no-existent device
	core.DataFile = "invalid-data1.json"
	data = core.LoadState()

	assert.Nil(t, data)

	// missing closing brace
	core.DataFile = "invalid-data2.json"
	data = core.LoadState()

	assert.Nil(t, data)

	// schedule refers to an no-existent program
	core.DataFile = "invalid-data3.json"
	data = core.LoadState()

	assert.Nil(t, data)

	// valid data
	core.DataFile = "data_test.json"
	data = core.LoadState()

	assert.NotNil(t, data)
	assert.Equal(t, 5, len(*data.Devices))

	core.DataFile = "data_test2.json"
	data.StoreState()

	str1, _ := ioutil.ReadFile("data_test.json")
	str2, _ := ioutil.ReadFile("data_test2.json")

	os.Remove("data_test2.json")
	assert.Equal(t, str1, str2)
}
