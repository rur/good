package service

type Routemap struct {
	Block  string
	Path   string
	Routes map[string]Routemap
}

type Sitemap map[string]Routemap

func (sm *Sitemap) GetPath(path string) (string, bool) {
	if path == "" {
		return "", false
	}
	parts := strings.Split(path, " > ")
	cursor := sm
	for i := 0; i < len(parts); i++ {
		cursor, ok = cursor.Routes[parts[i]]
		if !ok {
			break
		}
		if i == len(parts)-1 && cursor.Path != "" {
			return cursor.Path, true
		}
	}
	return "", false
}
