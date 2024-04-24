package util_test

import (
	"github.com/haapjari/repository-search-api/internal/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchLibrary_Success(t *testing.T) {
	url := "github.com/charmbracelet/charm@v0.12.6"

	// Call the function

	p, err := util.FetchLibrary(url)
	assert.NoError(t, err)

	t.Log(p)

	l, err := util.CalcLOC(p, "Go")
	assert.NoError(t, err)

	t.Log(p)
	t.Log(l)
}
