package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	missing := map[string]string{}

	filepath.Walk("_site", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".html" {
			return nil
		}

		vpath, _ := filepath.Rel("_site", path)
		if !strings.HasPrefix(vpath, "v") {
			return nil
		}

		i := strings.Index(vpath, "/")
		if i == -1 {
			return nil
		}
		p := vpath[i:]
		_, err = os.Stat(filepath.Join("_site", p))
		if os.IsNotExist(err) {
			missing[p] = strings.TrimSuffix(vpath, ".html")
			// fmt.Println("missing:", path)
		}
		return nil
	})

	paths := []string{}
	for path, _ := range missing {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%s,https://docs.projectcalico.org/%s\n", path, missing[path])
	}

	return nil
}
