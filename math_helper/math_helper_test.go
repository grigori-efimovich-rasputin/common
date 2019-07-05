package mathHelper


import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestClamp(t *testing.T) {
	assert := assert.New(t)

	ret := ClampInt(0, 100, 50)
	assert.Equal(ret, 50, "")

	ret = ClampInt(0, 100, -1)
	assert.Equal(ret, 0, "")

	ret = ClampInt(0, 100, 200)
	assert.Equal(ret, 100, "")
}


