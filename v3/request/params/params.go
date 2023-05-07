package params

import "strconv"

type URLParams map[string]string

func (u URLParams) Get(key string, def ...string) string {
	if v, ok := u[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (u URLParams) Has(key string) bool {
	_, ok := u[key]
	return ok
}

func (u URLParams) GetInt(key string, def ...int) int {
	var item = u.Get(key)
	i, err := strconv.Atoi(item)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return i
}
