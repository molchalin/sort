package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func Sort(dst io.Writer, src io.Reader) error {
	var lines []string
	scn := bufio.NewScanner(src)
	for scn.Scan() {
		lines = append(lines, scn.Text())
	}
	if err := scn.Err(); err != nil {
		return fmt.Errorf("read error: %v", err)
	}

	sort.Strings(lines)
	for i, line := range lines {
		var err error
		if i < len(lines)-1 {
			_, err = fmt.Fprintln(dst, line)
		} else {
			_, err = fmt.Fprint(dst, line)
		}
		if err != nil {
			return fmt.Errorf("write error: %v", err)
		}
	}
	return nil
}

func FilesSource(fileNames []string) (io.Reader, error) {
	if len(fileNames) == 0 {
		return bytes.NewReader(nil), nil
	}
	files := make([]io.Reader, 0, 2*len(fileNames)-1)
	for i, name := range fileNames {
		file, err := os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("can't open file: %v", err)
		}
		files = append(files, file)
		if i < len(fileNames)-1 {
			files = append(files, strings.NewReader("\n"))
		}
	}
	return io.MultiReader(files...), nil
}
