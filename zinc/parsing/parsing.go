package parsing

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var verbose = flag.Bool("v", false, "verbose logging")
var protobufRootDir = flag.String("protobufs", "", "root directory of protobuf sources")

var (
	itemCountRegexp      = regexp.MustCompile(`^([0-9]+) items$`)
	protobufSourceRegexp = regexp.MustCompile(`^// source: (.*)$`)
)

func Parse(path string, emit func(string, string)) error {
	if *verbose {
		log.Println("reading " + path)
	}
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
		if strings.Contains(path, "/protoc/") && *protobufRootDir != "" {
			if protobufPath := resolveProtobufPath(path); protobufPath != "" {
				path = protobufPath
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
	var protobufPath string
	err := withReader(genPath, func(r *bufio.Reader) error {
		for {
			line, err := readLine(r)
			if err == io.EOF {
				// we reached the end of the file without finding "source: " line, give up
				break
			}
			if err != nil {
				return err
			}
			matches := protobufSourceRegexp.FindStringSubmatch(line)

			if len(matches) == 2 {
				protobufPath = path.Join(*protobufRootDir, matches[1])
				break
			}
		}
		return nil
	})
	if err != nil {
		return ""
	} else {
		return protobufPath
	}
}
