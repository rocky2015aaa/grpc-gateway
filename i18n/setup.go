// Package i18n handles language translation
package i18n

import (
	//	"path"
	//	"runtime"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

var T i18n.TranslateFunc

func init() {
	path := ""
	pwd, _ := os.Getwd()
	if strings.Contains(pwd, "") {
		path = "../" + path
	}

	i18n.MustLoadTranslationFile(path)

	T, _ = i18n.Tfunc("")
}
