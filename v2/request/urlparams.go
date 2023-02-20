package request

import "strconv"

type URLParams map[string]string

func (u URLParams) GetDefault(key, def string) string {
	if v, ok := u[key]; ok {
		return v
	}
	return def
}

func (u URLParams) Get(key string) string {
	return u.GetDefault(key, "")
}

func (u URLParams) Set(key, value string) {
	u[key] = value
}

func (u URLParams) Has(key string) bool {
	_, ok := u[key]
	return ok
}

func (u URLParams) Delete(key string) {
	delete(u, key)
}

func (u URLParams) GetInt(key string) int {
	var zero = "0"
	var item = u.GetDefault(key, zero)
	i, err := strconv.Atoi(item)
	if err != nil {
		return 0
	}
	return i
}
