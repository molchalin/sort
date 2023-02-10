package internal

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort_basic(t *testing.T) {
	var tests = []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"basic", "a\nc\nd\nb", "a\nb\nc\nd"},
		{"basic 2", "aaab\naaaa", "aaaa\naaab"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := Sort(buf, strings.NewReader(tt.in))
			require.NoError(t, err)
			got, _ := io.ReadAll(buf)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("mock error")
}

func TestSort_write_err(t *testing.T) {
	err := Sort(new(errWriter), strings.NewReader("12312"))

	assert.Error(t, err)
}

func TestSort_read_err(t *testing.T) {
	err := Sort(nil, iotest.ErrReader(errors.New("mock error")))
	assert.Error(t, err)
}

func createFile(t *testing.T, content string) string {
	file, err := os.CreateTemp("", "*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(file.Name())
	})
	t.Cleanup(func() {
		_ = file.Close()
	})
	_, err = file.Write([]byte(content))
	require.NoError(t, err)
	return file.Name()
}

func TestFilesSource_basic(t *testing.T) {
	var tests = []struct {
		name  string
		want  string
		files []string
	}{
		{"empty", "", nil},
		{
			"basic",
			"a\nb\nc\nd\ne\nf",
			[]string{
				createFile(t, "a\nb\nc"),
				createFile(t, "d\ne\nf"),
			},
		},
	}
	for _, tt := range tests {
		t.Parallel()
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r, err := FilesSource(tt.files)
			require.NoError(t, err)

			got, err := io.ReadAll(r)
			require.NoError(t, err)

			assert.Equal(t, tt.want, string(got))
		})
	}
}
