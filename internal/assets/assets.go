package assets

import "fmt"

var assetsPath = ""

func init() {

}

func Setup(assets string) {
	assetsPath = assets
}

func FullPath(path string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", assetsPath, path)
}
