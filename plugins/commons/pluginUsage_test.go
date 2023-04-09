package commons

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallPath(t *testing.T) {
	p := &Plugin{
		Source:  "github.com/owner/repo",
		Version: "1.0.0",
		Name:    "test",
	}

	expectedPath := strings.Join([]string{"github.com/owner/repo", "1.0.0", "yatas-test"}, string(os.PathSeparator))
	assert.Equal(t, expectedPath, p.InstallPath())
}

func TestTagName(t *testing.T) {
	p := &Plugin{
		Version: "latest",
	}
	assert.Equal(t, "latest", p.TagName())

	p.Version = "1.0.0"
	assert.Equal(t, "v1.0.0", p.TagName())
}

func TestAssetName(t *testing.T) {
	p := &Plugin{
		Name: "test",
	}

	expectedName := strings.Join([]string{"yatas-test", runtime.GOOS, runtime.GOARCH}, "_") + ".zip"
	assert.Equal(t, expectedName, p.AssetName())
}

func TestValidate(t *testing.T) {
	p := &Plugin{
		Name:    "test",
		Version: "1.0.0",
		Source:  "github.com/owner/repo",
		Type:    "checks",
	}

	err := p.Validate()
	assert.Nil(t, err)

	p.Type = "invalid_type"
	err = p.Validate()
	assert.NotNil(t, err)
}
