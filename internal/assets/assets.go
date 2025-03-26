package assets

import (
	"net/url"
	"path"
)

var assetsPath = ""

func init() {

}

func Setup(assets string) {
	assetsPath = assets
}

func FullPath(p string) string {
	if p == "" {
		return ""
	}

	// 如果 p 是完整的 URL，則直接返回
	if u, err := url.ParseRequestURI(p); err == nil && u.Scheme != "" {
		return p
	}

	return path.Join(assetsPath, p)
}
