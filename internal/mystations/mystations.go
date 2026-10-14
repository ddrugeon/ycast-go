// Package mystations serves the stations of a stations.yml file, parsed with
// yaml.v3. The file is checked on every request (cheap stat call) and
// reloaded when it changed, so edits apply without restarting the server.
package mystations

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hgomez/ycast-go/internal/station"
	"gopkg.in/yaml.v3"
)

// Category is a named group of stations from stations.yml.
type Category struct {
	Name     string
	Stations []*station.Station
}

// MyStations is the hot-reloading source backed by a stations.yml file.
type MyStations struct {
	path string
	log  *slog.Logger

	mu      sync.Mutex
	mtime   time.Time
	size    int64
	cats    []Category
	loaded  bool
	lastErr string
}

// New returns the source for the given file. A nil *MyStations (empty path)
// behaves like an empty source.
func New(path string, log *slog.Logger) *MyStations {
	return &MyStations{path: path, log: log}
}

// Categories returns the current categories, reloading the file if needed.
func (m *MyStations) Categories() []Category {
	if m == nil || m.path == "" {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	fi, err := os.Stat(m.path)
	if err != nil {
		if m.lastErr != err.Error() {
			m.log.Error("station configuration not readable", "file", m.path, "err", err)
			m.lastErr = err.Error()
		}
		m.cats, m.loaded = nil, false
		return nil
	}
	if m.loaded && fi.ModTime().Equal(m.mtime) && fi.Size() == m.size {
		return m.cats
	}
	data, err := os.ReadFile(m.path)
	if err == nil {
		var cats []Category
		cats, err = ParseYAML(data)
		if err == nil {
			m.cats, m.mtime, m.size, m.loaded, m.lastErr = cats, fi.ModTime(), fi.Size(), true, ""
			m.log.Info("station configuration loaded", "file", m.path, "categories", len(cats))
			return m.cats
		}
	}
	// Keep serving the last good version when the file is half-edited.
	m.mtime, m.size = fi.ModTime(), fi.Size()
	if m.lastErr != err.Error() {
		m.log.Error("station configuration format error", "file", m.path, "err", err)
		m.lastErr = err.Error()
	}
	return m.cats
}

// Category returns the named category, or nil.
func (m *MyStations) Category(name string) *Category {
	cats := m.Categories()
	for i := range cats {
		if cats[i].Name == name {
			return &cats[i]
		}
	}
	return nil
}

// ByID finds a station by its full ID.
func (m *MyStations) ByID(id string) *station.Station {
	for _, c := range m.Categories() {
		for _, s := range c.Stations {
			if s.ID == id {
				return s
			}
		}
	}
	return nil
}

// Search returns the stations whose name contains the query
// (case-insensitive).
func (m *MyStations) Search(query string) []*station.Station {
	query = strings.ToLower(query)
	var out []*station.Station
	for _, c := range m.Categories() {
		for _, s := range c.Stations {
			if strings.Contains(strings.ToLower(s.Name), query) {
				out = append(out, s)
			}
		}
	}
	return out
}

// ParseYAML parses stations.yml with yaml.v3, walking the document as an
// ordered yaml.Node tree so that file order is preserved (maps would sort
// keys). Top-level keys are categories, their values map station names to
// URLs. Duplicate categories or station names keep their original position
// with the last definition winning, like conventional YAML loaders.
func ParseYAML(data []byte) ([]Category, error) {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if doc.Kind == 0 { // empty file
		return nil, nil
	}
	root := doc.Content[0]
	if isNull(root) {
		return nil, nil
	}
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("top level must be a mapping of categories")
	}

	var cats []Category
	catIndex := map[string]int{} // category name -> position in cats
	for i := 0; i+1 < len(root.Content); i += 2 {
		key, val := root.Content[i], root.Content[i+1]
		name := key.Value
		var mapping *yaml.Node
		if val.Kind == yaml.MappingNode {
			mapping = val
		} else if !isNull(val) { // "Category:" with no stations yet is fine
			return nil, fmt.Errorf("category %q must map station names to URLs", name)
		}
		c := Category{Name: name}
		stationIndex := map[string]int{}
		if mapping != nil {
			for j := 0; j+1 < len(mapping.Content); j += 2 {
				stKey, stVal := mapping.Content[j], mapping.Content[j+1]
				if isNull(stVal) {
					return nil, fmt.Errorf("station %q has no URL (nested categories are not supported)", stKey.Value)
				}
				if stVal.Kind != yaml.ScalarNode {
					return nil, fmt.Errorf("station %q must map to a URL", stKey.Value)
				}
				if k, dup := stationIndex[stKey.Value]; dup { // last one wins, keeps position
					c.Stations[k] = newStation(name, stKey.Value, stVal.Value)
					continue
				}
				stationIndex[stKey.Value] = len(c.Stations)
				c.Stations = append(c.Stations, newStation(name, stKey.Value, stVal.Value))
			}
		}
		if k, dup := catIndex[name]; dup { // last one wins, keeps position
			cats[k].Stations = c.Stations
			continue
		}
		catIndex[name] = len(cats)
		cats = append(cats, c)
	}
	return cats, nil
}

func isNull(n *yaml.Node) bool {
	return n.Kind == yaml.ScalarNode && (n.Tag == "!!null" || n.Value == "")
}

func newStation(category, name, url string) *station.Station {
	return &station.Station{
		ID:   station.MakeID(station.MyPrefix, station.MyChecksum(name, url)),
		Name: name,
		URL:  url,
		Tags: []string{category},
	}
}
