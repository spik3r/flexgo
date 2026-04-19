package flexgo

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

// Each entry is a path under ./example/ (use forward slashes). The
// golden file name is the path with slashes replaced by underscores,
// so example/basics/align → testdata/basics_align.golden.
var exampleGoldens = []string{
	"basics/align",
	"basics/basic",
	"basics/centered",
	"basics/justify",
	"basics/spacing",
	"builder/alignself",
	"builder/basic",
	"dynamic",
	"layouts/headerbodyfooter",
	"layouts/modal",
	"margins/centered_layout",
	"margins/hautocenter",
	"margins/vautocenter",
}

func TestExampleGolden(t *testing.T) {
	for _, path := range exampleGoldens {
		t.Run(path, func(t *testing.T) {
			out, err := runExample(path)
			if err != nil {
				t.Fatalf("failed to run example %s: %v", path, err)
			}

			goldenPath := filepath.Join("testdata", goldenName(path)+".golden")
			if *update {
				if err := os.WriteFile(goldenPath, []byte(out), 0o644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				return
			}

			expected, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			if out != string(expected) {
				t.Fatalf("example %s output mismatch", path)
			}
		})
	}
}

func goldenName(examplePath string) string {
	return strings.ReplaceAll(examplePath, "/", "_")
}

func runExample(path string) (string, error) {
	cmd := exec.Command("go", "run", "./example/"+path)
	cmd.Env = append(os.Environ(), "FLEXGO_GOLDEN=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
