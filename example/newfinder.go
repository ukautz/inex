package main

import (
	"fmt"
	"github.com/ukautz/inex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var help = fmt.Sprintf(`
Usage: %s <parameters>

parameters:
  (-p|--path) <path>
    path to directory to find recursively in
  (-e|--exclude) <regex>                            [multiple allowed]
    regular expression matching file paths which are not to be displayed
  (-i|--include) <regex>                            [multiple allowed]
    regular expression matching file paths which are to be displayed
  (-c|--contains) <string>                          [multiple allowed]
    this string must be contained in paths which are to be displayed
  (-C|--not-contains) <string>                      [multiple allowed]
    this string must NOT be contained in paths which are to be displayed
  (-t|--type) (f|file|d|directory)
    whether path to be displayed must be file or directory
`, os.Args[0])

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

// looks like this, when run from inex folder:
//	$ go run example/main.go -i ".*\.go" -i "(_test|inex)" -e "\/"
//	inex.go
//	inex_test.go
//	matcher_test.go
func main() {
	l := len(os.Args)
	if l <= 1 || l%2 != 1 {
		die(help)
	}
	filter := inex.NewRoot()
	path := "."
	for i := 1; i < l; i += 2 {
		switch os.Args[i] {
		case "-e", "--exclude":
			filter = filter.Exclude(&inex.RegexpMatcher{regexp.MustCompile(os.Args[i+1])})
		case "-i", "--include":
			filter = filter.Include(&inex.RegexpMatcher{regexp.MustCompile(os.Args[i+1])})
		case "-c", "--contains":
			partial := os.Args[i+1]
			filter = filter.Include(inex.FuncMatcher(func(str string) bool {
				return strings.Contains(str, partial)
			}))
		case "-C", "--not-contains":
			partial := os.Args[i+1]
			filter = filter.Include(inex.FuncMatcher(func(str string) bool {
				return !strings.Contains(str, partial)
			}))
		case "-t", "--type":
			typ := os.Args[i+1]
			switch typ {
			case "f", "file", "d", "directory":
			default:
				die("invalid --type provided. must be \"f\", \"file\", \"d\" or \"directory\"")
			}
			filter = filter.Include(inex.FuncMatcher(func(str string) bool {
				stat, err := os.Stat(str)
				if err != nil {
					die("failed to stat \"%s\": %s\n", str, err)
				}
				if typ == "f" || typ == "file" {
					return !stat.IsDir()
				}
				return stat.IsDir()
			}))
		case "-p", "--path":
			if path != "." {
				die("%s can only be provided once\n\n%s", os.Args[i], help)
			}
			path = os.Args[i+1]
		default:
			die("unsupported operation \"%s\"\n\n%s", os.Args[i], help)
		}
	}
	if filter.IsRoot() {
		die("no filters provided\n\n%s", help)
	}
	filter = filter.Root()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if filter.Match(path) {
			fmt.Println(path)
		}
		return nil
	})
}
