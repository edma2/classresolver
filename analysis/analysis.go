package analysis

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	itemCountRegexp      = regexp.MustCompile(`^([0-9]+) items$`)
	protobufSourceRegexp = regexp.MustCompile(`^// source: (.*)$`)

	PantsRoot string
)

func IsAnalysisFile(path string) bool {
	return strings.HasSuffix(path, ".analysis") && isRegularFile(path)
}

func isRegularFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

func ReadAnalysisFile(path string, emit func(string, string)) error {
	log.Println("reading" + path)
	return withReader(path, func(r *bufio.Reader) error {
		if err := readUntil(r, "class names:"); err != nil {
			return err
		}
		if err := readClassNames(r, emit); err != nil {
			return err
		}
		return nil
	})
}

func readUntil(r *bufio.Reader, s string) error {
	for {
		line, err := readLine(r)
		if strings.Contains(line, s) {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func withReader(path string, f func(*bufio.Reader) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	return f(r)
}

func readClassNames(r *bufio.Reader, emit func(string, string)) error {
	var itemCount int

	itemCountLine, err := readLine(r)
	if err != nil {
		return err
	}
	matches := itemCountRegexp.FindStringSubmatch(itemCountLine)

	if len(matches) != 2 {
		return errors.New("unexpected item count line: " + itemCountLine)
	} else {
		itemCountString := matches[1]
		itemCount, err = strconv.Atoi(itemCountString)
		if err != nil {
			return err
		}
	}

	for i := 0; i < itemCount; i++ {
		l, err := readLine(r)
		if err != nil {
			return err
		}
		path, class := parseMapping(l)

		if strings.Contains(path, "/.pants.d/gen/protoc/") {
			if protoPath := resolveProtobufPath(path); protoPath != "" {
				path = protoPath
			}
		}
		emit(class, path)
	}

	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	return strings.TrimSpace(line), err
}

// foo -> bar => (foo, bar)
func parseMapping(line string) (string, string) {
	split := strings.SplitN(line, " ", 3)
	return strings.TrimSpace(split[0]), strings.TrimSpace(split[2])
}

func resolveProtobufPath(genPath string) string {
	var protoPath string
	err := withReader(genPath, func(r *bufio.Reader) error {
		for {
			line, err := readLine(r)
			if err == io.EOF {
				// This shouldn't happen?
				break
			}
			if err != nil {
				return err
			}
			matches := protobufSourceRegexp.FindStringSubmatch(line)

			if len(matches) == 2 {
				protoPath = path.Join(PantsRoot, "science/src/protobuf", matches[1])
				break
			}
		}
		return nil
	})
	if err != nil {
		return ""
	} else {
		return protoPath
	}
}
