package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullAssetsPathEmpty(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	assetsPath = "https://cdn.example.com"
	assert.Equal(t, "", FullAssetsPath(""))
}

func TestFullAssetsPathHttpsUrl(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	p := "https://127.0.0.1"
	assert.Equal(t, p, FullAssetsPath(p))
}

func TestFullAssetsPathHttpUrl(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	p := "http://127.0.0.1"
	assert.Equal(t, p, FullAssetsPath(p))
}

func TestFullAssetsPathJoinBaseSlash(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	assetsPath = "https://cdn.example.com/"
	p := "a.png"

	assert.Equal(t, "https://cdn.example.com/a.png", FullAssetsPath(p))
}

func TestFullAssetsPathJoinPathSlash(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	assetsPath = "https://cdn.example.com"
	p := "/a.png"

	assert.Equal(t, "https://cdn.example.com/a.png", FullAssetsPath(p))
}

func TestFullAssetsPathHostLike(t *testing.T) {
	orig := assetsPath
	t.Cleanup(func() { assetsPath = orig })

	assetsPath = "https://cdn.example.com"
	p := "127.0.0.1/a.png"

	assert.Equal(t, "https://cdn.example.com/127.0.0.1/a.png", FullAssetsPath(p))
}
