package assets

import (
	"net/url"
	"strings"
)

var assetsPath = ""
var assets2Path = ""

func init() {

}

func Setup(assets, assets2 string) {
	assetsPath = assets
	assets2Path = assets2
}

// FullAssetsPath
// assets 1完整連結
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

// FullAssets2Path
// assets 2完整連結
func FullAssets2Path(p string) string {
	if p == "" {
		return ""
	}
	// 如果 p 是完整的 URL，則直接返回
	if u, err := url.ParseRequestURI(p); err == nil && u.Scheme != "" {
		return p
	}

	return strings.TrimRight(assets2Path, "/") + "/" + strings.TrimLeft(p, "/")
}
