/*
Package gotext implements GNU gettext utilities.

For quick/simple translations you can use the package level functions directly.

    import (
	    "fmt"
	    "github.com/tanyinloo/gotext"
    )

    func main() {
        // Configure package
        gotext.Configure("/path/to/locales/root/dir", "en_UK", "domain-name")

        // Translate text from default domain
        fmt.Println(gotext.Get("My text on 'domain-name' domain"))

        // Translate text from a different domain without reconfigure
        fmt.Println(gotext.GetD("domain2", "Another text on a different domain"))
    }

*/
package gotext

import (
	"encoding/gob"
	"sync"
)

// Global environment variables
type config struct {
	sync.RWMutex
	// Default domain to look at when no domain is specified. Used by package level functions.
	domain string

	// Language set.
	language string

	// Path to library directory where all locale directories and Translation files are.
	library string

	// Storage for package level methods
	storage *Locale
}

// var globalConfig *config
var configMap map[string]*config
var defaultLang string

func init() {
	// Init default configuration
	// globalConfig = &config{
	// 	domain:   "default",
	// 	language: "en_US",
	// 	library:  "/usr/local/share/locale",
	// 	storage:  nil,
	// }
	configMap = make(map[string]*config)

	// Register Translator types for gob encoding
	gob.Register(TranslatorEncoding{})
}

// loadStorage creates a new Locale object at package level based on the Global variables settings.
// It's called automatically when trying to use Get or GetD methods.
func loadStorage(lang string) {
	config, ok := configMap[lang]
	if !ok {
		return
	}
	config.Lock()
	defer config.Unlock()
	storage := config.storage
	if storage == nil {
		storage = NewLocale(config.library, config.language)
	}

	if _, ok := storage.Domains[config.domain]; !ok {
		storage.AddDomain(config.domain)
	}
	storage.SetDomain(config.domain)
	config.storage = storage
	configMap[lang] = config
}

// GetDomain is the domain getter for the package configuration
func GetDomain(lang string) string {
	var dom string

	config := configMap[SimplifiedLocale(lang)]
	config.RLock()
	defer config.RUnlock()
	if config.storage != nil {
		dom = config.storage.GetDomain()
	}
	if dom == "" {
		dom = config.domain
	}

	return dom
}

// SetDomain sets the name for the domain to be used at package level.
// It reloads the corresponding Translation file.
// func SetDomain(dom string) {
// 	globalConfig.Lock()
// 	globalConfig.domain = dom
// 	if globalConfig.storage != nil {
// 		globalConfig.storage.SetDomain(dom)
// 	}
// 	globalConfig.Unlock()

// 	loadStorage(true)
// }

// GetLanguage is the language getter for the package configuration
func GetLanguage() string {
	return defaultLang
}

// SetLanguage sets the language code to be used at package level.
// It reloads the corresponding Translation file.
func SetLanguage(l string) {
	defaultLang = SimplifiedLocale(l)
}

// GetLibrary is the library getter for the package configuration
func GetLibrary(lang string) string {
	return configMap[lang].library
}

// SetLibrary sets the root path for the loale directories and files to be used at package level.
// It reloads the corresponding Translation file.
// func SetLibrary(lib string) {
// 	globalConfig.Lock()
// 	globalConfig.library = lib
// 	globalConfig.Unlock()

// 	loadStorage(true)
// }

// Configure sets all configuration variables to be used at package level and reloads the corresponding Translation file.
// It receives the library path, language code and domain name.
// This function is recommended to be used when changing more than one setting,
// as using each setter will introduce a I/O overhead because the Translation file will be loaded after each set.
func AddConfig(lib, lang, dom string) {
	simplifyLang := SimplifiedLocale(lang)
	c := config{
		domain:   dom,
		library:  lib,
		language: simplifyLang,
	}
	configMap[simplifyLang] = &c
	defaultLang = lang
	loadStorage(simplifyLang)
}

// Get uses the default domain globally set to return the corresponding Translation of a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func Get(str string, vars ...interface{}) string {
	return GetD(GetDomain(defaultLang), str, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string in the default domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetN(str, plural string, n int, vars ...interface{}) string {
	return GetND(GetDomain(defaultLang), str, plural, n, vars...)
}

// GetD returns the corresponding Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetD(dom, str string, vars ...interface{}) string {
	config, ok := configMap[defaultLang]
	if !ok {
		return Printf(str, vars...)
	}
	config.RLock()
	defer config.RUnlock()
	// Return Translation

	if _, ok := config.storage.Domains[dom]; !ok {
		config.storage.AddDomain(dom)
	}

	tr := config.storage.GetD(dom, str, vars...)

	return tr
}

// GetND retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetND(dom, str, plural string, n int, vars ...interface{}) string {
	config, ok := configMap[defaultLang]
	if !ok {
		return Printf(str, vars...)
	}
	config.RLock()
	defer config.RUnlock()
	if _, ok := config.storage.Domains[dom]; !ok {
		config.storage.AddDomain(dom)
	}

	tr := config.storage.GetND(dom, str, plural, n, vars...)

	return tr
}

// GetC uses the default domain globally set to return the corresponding Translation of the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetC(str, ctx string, vars ...interface{}) string {
	return GetDC(GetDomain(defaultLang), str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context in the default domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	return GetNDC(GetDomain(defaultLang), str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetDC(dom, str, ctx string, vars ...interface{}) string {
	config, ok := configMap[defaultLang]
	if !ok {
		return Printf(str, vars...)
	}
	config.RLock()
	defer config.RUnlock()
	// Return Translation

	tr := config.storage.GetDC(dom, str, ctx, vars...)

	return tr
}

// GetNDC retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNDC(dom, str, plural string, n int, ctx string, vars ...interface{}) string {
	config, ok := configMap[defaultLang]
	if !ok {
		return Printf(str, vars...)
	}
	config.RLock()
	defer config.RUnlock()
	// Return Translation
	tr := config.storage.GetNDC(dom, str, plural, n, ctx, vars...)

	return tr
}
