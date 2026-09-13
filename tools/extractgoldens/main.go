// Command extractgoldens extracts the worked JSON request/event examples from the vendored Core-API
// entity docs (spec/core-api/doc/entities/*.md) into individual golden fixture files under
// pkg/entities/testdata, for the entity types goucrt currently implements.
//
// The Core-API AsyncAPI spec only defines `features` and `device_class` as structured, machine-checkable
// enums (see internal/spec); cmd_id, attribute names and state values exist only as prose and worked
// examples in these docs. This tool is how that prose becomes something pkg/entities' tests can check
// goucrt's behavior against mechanically, instead of by eye.
//
// Run it from the repository root after updating spec/core-api (see spec/core-api/README.md):
//
//	go run ./tools/extractgoldens
//
// It fully regenerates pkg/entities/testdata/<entity>/{commands,events}/ for each entity listed in
// implementedEntities below - review the diff afterwards, since a spec update can rename or remove an
// example goucrt's tests reference by filename.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// implementedEntities maps a vendored doc's basename (without the entity_/.md wrapper) to the entity
// output directory name under pkg/entities/testdata. Only entity types goucrt actually implements are
// listed; add a type here once pkg/entities gains a constructor for it.
var implementedEntities = map[string]string{
	"button":       "button",
	"switch":       "switch",
	"light":        "light",
	"cover":        "cover",
	"media_player": "media_player",
	"climate":      "climate",
	"sensor":       "sensor",
	"remote":       "remote",
}

var (
	headingRE    = regexp.MustCompile(`^(#{2,4})\s+(.*)$`)
	fenceOpenRE  = regexp.MustCompile("^```json\\s*$")
	fenceCloseRE = regexp.MustCompile("^```\\s*$")
	slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
)

type section int

const (
	sectionNone section = iota
	sectionCommands
	sectionEvents
)

func main() {
	specDir := "internal/spec/core-api/doc/entities"
	outRoot := "pkg/entities/testdata"

	entries, err := os.ReadDir(specDir)
	if err != nil {
		fatalf("read %s: %v", specDir, err)
	}

	total := 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasPrefix(name, "entity_") || !strings.HasSuffix(name, ".md") {
			continue
		}
		stem := strings.TrimSuffix(strings.TrimPrefix(name, "entity_"), ".md")
		outDir, ok := implementedEntities[stem]
		if !ok {
			continue // not yet implemented by goucrt; nothing to check it against.
		}

		n, err := extractFile(filepath.Join(specDir, name), filepath.Join(outRoot, outDir))
		if err != nil {
			fatalf("%s: %v", name, err)
		}
		fmt.Printf("%-14s %2d example(s)\n", stem, n)
		total += n
	}
	fmt.Printf("%d example(s) extracted\n", total)
}

func extractFile(path, outDir string) (count int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Clear previously generated output for this entity so removed/renamed examples don't linger.
	if err := os.RemoveAll(outDir); err != nil {
		return 0, err
	}

	sec := sectionNone
	heading := ""
	headingCount := map[string]int{}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()

		if m := headingRE.FindStringSubmatch(line); m != nil {
			level, text := len(m[1]), strings.TrimSpace(m[2])
			switch {
			case level == 3:
				switch {
				case strings.EqualFold(text, "Command examples"):
					sec = sectionCommands
				case strings.EqualFold(text, "Event examples"):
					sec = sectionEvents
				default:
					sec = sectionNone
				}
			case level == 4 && sec != sectionNone:
				heading = text
			}
			continue
		}

		if sec == sectionNone || !fenceOpenRE.MatchString(line) {
			continue
		}

		var buf strings.Builder
		for scanner.Scan() {
			body := scanner.Text()
			if fenceCloseRE.MatchString(body) {
				break
			}
			buf.WriteString(body)
			buf.WriteByte('\n')
		}

		var js interface{}
		if err := json.Unmarshal([]byte(buf.String()), &js); err != nil {
			return count, fmt.Errorf("invalid JSON in %q block under %q: %w", heading, path, err)
		}
		pretty, err := json.MarshalIndent(js, "", "  ")
		if err != nil {
			return count, err
		}

		sub := "events"
		if sec == sectionCommands {
			sub = "commands"
		}
		dir := filepath.Join(outDir, sub)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return count, err
		}

		slug := slugify(heading)
		headingCount[sub+"/"+slug]++
		if n := headingCount[sub+"/"+slug]; n > 1 {
			slug = fmt.Sprintf("%s_%d", slug, n)
		}

		outPath := filepath.Join(dir, slug+".json")
		pretty = append(pretty, '\n')
		if err := os.WriteFile(outPath, pretty, 0o644); err != nil {
			return count, err
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugNonAlnum.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
