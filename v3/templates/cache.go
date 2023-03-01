package templates

import (
	"html/template"
	"sync"
)

type tCache struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

func newCache() *tCache {
	return &tCache{templates: make(map[string]*template.Template)}
}

func (tc *tCache) Get(key string) (*template.Template, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	t, ok := tc.templates[key]
	return t, ok
}

func (tc *tCache) Set(key string, value *template.Template) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.templates[key] = value
}
