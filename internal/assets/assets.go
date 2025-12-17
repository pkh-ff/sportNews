package assets

import (
	"net/url"
	"strings"
)

var assetsPath = ""

func init() {

}

func Setup(assets string) {
	assetsPath = assets
}

// FullAssetsPath
// assets 完整連結
func FullAssetsPath(p string) string {
	if p == "" {
		return ""
	}

	// 如果 p 是完整的 URL，則直接返回
	if u, err := url.ParseRequestURI(p); err == nil && u.Scheme != "" {
		return p
	}

	return strings.TrimRight(assetsPath, "/") + "/" + strings.TrimLeft(p, "/")
}
