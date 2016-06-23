[![GoDoc](https://godoc.org/github.com/leonelquinteros/gotext?status.svg)](https://godoc.org/github.com/leonelquinteros/gotext)

# Gotext

GNU gettext utilities for Go 


# Examples

## Using package for single language/domain settings

For quick/simple translations on a single file, you can use the package level functions directly.

```go
import "github.com/leonelquinteros/gotext"

func main() {
    // Configure package
    gotext.Configure("/path/to/locales/root/dir", "en_UK", "domain-name")
    
    // Translate text from default domain
    println(gotext.Get("My text on 'domain-name' domain"))
    
    // Translate text from a different domain without reconfigure
    println(gotext.GetD("domain2", "Another text on a different domain"))
}

```

## Using dynamic variables on translations

All translation strings support dynamic variables to be inserted without translate. 
Use the fmt.Printf syntax (from Go's "fmt" package) to specify how to print the non-translated variable inside the translation string. 

```go
import "github.com/leonelquinteros/gotext"

func main() {
    // Configure package
    gotext.Configure("/path/to/locales/root/dir", "en_UK", "domain-name")
    
    // Set variables
    name := "John"
    
    // Translate text with variables
    println(gotext.Get("Hi, my name is %s", name))
}

```


## Using Locale object

When having multiple languages/domains/libraries at the same time, you can create Locale objects for each variation 
so you can handle each settings on their own.

```go
import "github.com/leonelquinteros/gotext"

func main() {
    // Create Locale with library path and language code
    l := gotext.NewLocale("/path/to/locales/root/dir", "es_UY")
    
    // Load domain '/path/to/locales/root/dir/es_UY/default.po'
    l.AddDomain("default")
    
    // Translate text from default domain
    println(l.Get("Translate this"))
    
    // Load different domain
    l.AddDomain("translations")
    
    // Translate text from domain
    println(l.GetD("translations", "Translate this"))
}
```

This is also helpful for using inside templates (from the "text/template" package), where you can pass the Locale object to the template.
If you set the Locale object as "Loc" in the template, then the templace code would look like: 

```
{{ .Loc.Get "Translate this" }}
```


## Using the Po object to handle .po files and PO-formatted strings

For when you need to work with PO files and strings, 
you can directly use the Po object to parse it and access the translations in there in the same way.

```go
import "github.com/leonelquinteros/gotext"

func main() {
    // Set PO content
    str := `
msgid "Translate this"
msgstr "Translated text"

msgid "Another string"
msgstr ""

msgid "One with var: %s"
msgstr "This one sets the var: %s"
`
    
    // Create Po object
    po := new(Po)
    po.Parse(str)
    
    println(po.Get("Translate this"))
}
```


## Use plural forms of translations

PO format supports defining one or more plural forms for the same translation.

```go
import "github.com/leonelquinteros/gotext"

func main() {
    // Set PO content
    str := `
msgid "Translate this"
msgstr "Translated text"

msgid "Another string"
msgstr ""

msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular: %s"
msgstr[1] "This one is the plural: %s"
msgstr[2] "And this is the second plural form: %s"
`
    
    // Create Po object
    po := new(Po)
    po.Parse(str)
    
    println(po.GetN("One with var: %s", "Several with vars: %s", 2, v))
    // "And this is the second plural form: Variable"
}
```


# Contribute 

- Please, contribute.
- Use the package on your projects.
- Report issues on Github. 
- Send pull requests for bugfixes and improvements.
- Send proposals on Github issues.
