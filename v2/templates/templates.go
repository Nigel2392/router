package templates

import (
	"errors"
	"html/template"
	"io/fs"
	"strings"
)

var (
	// Default base template suffixes
	BASE_TEMPLATE_SUFFIXES = []string{".tmpl", ".html"}
	// Default directory to look in for base templates
	BASE_TEMPLATE_DIRS = []string{"templates/base"}
	// Default directory to look in for templates
	TEMPLATE_DIRS = []string{"templates"}
	// Functions to add to templates
	DEFAULT_FUNCS = make(template.FuncMap)
	// Template file system
	TEMPLATEFS fs.FS
)

var USE_TEMPLATE_CACHE = true

func GetTemplate(templateName string) (*template.Template, string, error) {
	// Check if template is cached
	var t *template.Template
	var ok bool
	if t, ok = templateCache.Get(templateName); !ok || !USE_TEMPLATE_CACHE {
		// If not, cache it
		var base_template_dirs = BASE_TEMPLATE_DIRS
		var directories = TEMPLATE_DIRS
		var extensions = BASE_TEMPLATE_SUFFIXES

		// Search fs for all base templates, in every base directory
		var base_templates = make([]string, 0)
		for _, base_template_dir := range base_template_dirs {
			// Read all files in base template directory
			files, err := fs.ReadDir(TEMPLATEFS, base_template_dir)
			if err != nil {
				return nil, "", errors.New("Error reading base template directory: " + base_template_dir + " (" + err.Error() + ")")
			}
			// Add all files to base templates
			for _, file := range files {
				var name = file.Name()
				// Check if file is a template
				for _, extension := range extensions {
					if name[len(name)-len(extension):] == extension {
						base_templates = append(base_templates, base_template_dir+"/"+file.Name())
					}
				}
			}
		}
		var template_name string
		// Search fs for all templates, in every directory
		if len(directories) > 0 {
			for _, directory := range directories {
				// Check if file exists
				var dirName = NicePath(false, directory, templateName)
				if _, err := fs.Stat(TEMPLATEFS, dirName); err == nil {
					template_name = dirName
					break
				}
			}
		} else {
			template_name = NicePath(false, templateName)
		}
		var err error
		var t = template.New(template_name)
		t.Funcs(DEFAULT_FUNCS)
		t, err = t.ParseFS(TEMPLATEFS, append(base_templates, template_name)...)
		if err != nil {
			return nil, "", err
		}
		templateCache.Set(templateName, t)

		// Render template
		return t, FilenameFromPath(template_name), nil
	}
	var name = FilenameFromPath(templateName)
	if t == nil {
		var err = errors.New("template not found")
		return nil, "", err
	}
	return t, name, nil
}

func NicePath(forceSuffixSlash bool, p ...string) string {
	var b strings.Builder
	for i, s := range p {
		s = strings.Replace(s, "\\", "/", -1)
		if s == "/" {
			b.WriteString(s)
			continue
		}
		if i != 0 {
			s = strings.TrimPrefix(s, "/")
		}
		if i == len(p)-1 && forceSuffixSlash && !strings.HasSuffix(s, "/") || i != len(p)-1 && !strings.HasSuffix(s, "/") {
			s += "/"
		}
		b.WriteString(s)
	}
	return b.String()
}

func NameFromPath(p string) string {
	var name = FilenameFromPath(p)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[0]
	}
	return name
}

func FilenameFromPath(p string) string {
	p = strings.Replace(p, "\\", "/", -1)
	name := strings.Split(p, "/")[len(strings.Split(p, "/"))-1]
	return name
}
