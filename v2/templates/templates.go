package templates

import (
	"errors"
	"html/template"
	"io/fs"
	"strings"
)

// Template manager
// Used for easily fetching templates.
type Manager struct {
	// Use template cache?
	USE_TEMPLATE_CACHE bool
	// Template cache to use, if enabled
	tCache *tCache
	// Default base template suffixes
	BASE_TEMPLATE_SUFFIXES []string
	// Default directory to look in for base templates
	BASE_TEMPLATE_DIRS []string
	TEMPLATE_DIRS      []string
	// Functions to add to templates
	DEFAULT_FUNCS template.FuncMap
	// Template file system
	TEMPLATEFS fs.FS
}

func (tm *Manager) cache() *tCache {
	if tm.tCache == nil {
		tm.tCache = newCache()
	}
	return tm.tCache
}

// Initialize the template manager
func (tm *Manager) Init() {
	tm.tCache = newCache()
}

// Get a template
func (tm *Manager) Get(templateName string, funcs template.FuncMap) (*template.Template, string, error) {
	// Check if template is cached
	var t *template.Template
	var ok bool
	if t, ok = tm.cache().Get(templateName); !ok || !tm.USE_TEMPLATE_CACHE {
		// If not, cache it
		var base_template_dirs = tm.BASE_TEMPLATE_DIRS
		var directories = tm.TEMPLATE_DIRS
		var extensions = tm.BASE_TEMPLATE_SUFFIXES

		// Search fs for all base templates, in every base directory
		var base_templates = make([]string, 0)
		for _, base_template_dir := range base_template_dirs {
			// Read all files in base template directory
			files, err := fs.ReadDir(tm.TEMPLATEFS, base_template_dir)
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
				var _, err = fs.Stat(tm.TEMPLATEFS, dirName)
				if err == nil {
					template_name = dirName
					break
				}
			}
		}
		if template_name == "" {
			template_name = NicePath(false, templateName)
		}
		var err error
		var t = template.New(template_name)
		t.Funcs(tm.DEFAULT_FUNCS)
		if funcs != nil {
			t.Funcs(funcs)
		}
		t, err = t.ParseFS(tm.TEMPLATEFS, append(base_templates, template_name)...)
		if err != nil {
			return nil, "", errors.New("Error parsing template: " + template_name + " (" + err.Error() + ")")
		}
		tm.cache().Set(templateName, t)

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

// Render a template from a string
func (tm *Manager) GetFromString(templateString string, templateName string) (*template.Template, error) {
	var t = template.New(templateName)
	t.Funcs(tm.DEFAULT_FUNCS)
	var err error
	t, err = t.Parse(templateString)
	if err != nil {
		return nil, errors.New("Error parsing template: " + templateName + " (" + err.Error() + ")")
	}
	return t, nil
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

func FilenameFromPath(p string) string {
	p = strings.Replace(p, "\\", "/", -1)
	name := strings.Split(p, "/")[len(strings.Split(p, "/"))-1]
	return name
}
