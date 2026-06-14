package colormap

import (
	"fmt"
	"sort"
)

type Algorithm interface {
	Name() string
	Partition(px *Pixels, opts Options) []Region
}

var registry = map[string]Algorithm{}

func Register(a Algorithm) {
	if a == nil {
		panic("colormap: Register(nil)")
	}

	name := a.Name()
	if _, dup := registry[name]; dup {
		panic("colormap: algorithm already registered: " + name)
	}

	registry[name] = a
}

func Algorithms() []string {
	names := make([]string, 0, len(registry))

	for n := range registry {
		names = append(names, n)
	}

	sort.Strings(names)

	return names
}

func getAlgorithm(name string) (Algorithm, error) {
	if a, ok := registry[name]; ok {
		return a, nil
	}

	return nil, fmt.Errorf("colormap: unknown algorithm %q (available: %v)", name, Algorithms())
}
