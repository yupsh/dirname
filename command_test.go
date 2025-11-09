package command_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/dirname"
)

// ==============================================================================
// Test Basic Functionality
// ==============================================================================

func TestDirname_SimplePath(t *testing.T) {
	result := run.Quick(command.Dirname("/usr/local/bin/script.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/usr/local/bin"})
}

func TestDirname_NoPath(t *testing.T) {
	result := run.Quick(command.Dirname("script.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

func TestDirname_Root(t *testing.T) {
	result := run.Quick(command.Dirname("/"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/"})
}

func TestDirname_CurrentDir(t *testing.T) {
	result := run.Quick(command.Dirname("."))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

func TestDirname_ParentDir(t *testing.T) {
	result := run.Quick(command.Dirname(".."))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

// ==============================================================================
// Test Trailing Slashes
// ==============================================================================

func TestDirname_TrailingSlash(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/dir/"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/path/to"})
}

func TestDirname_MultipleTrailingSlashes(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/dir///"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/path/to"})
}

// ==============================================================================
// Test Multiple Paths
// ==============================================================================

func TestDirname_MultiplePaths(t *testing.T) {
	result := run.Quick(command.Dirname(
		"/usr/bin/script.sh",
		"/home/user/doc.txt",
		"/tmp/temp",
	))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"/usr/bin",
		"/home/user",
		"/tmp",
	})
}

// ==============================================================================
// Test Different Path Styles
// ==============================================================================

func TestDirname_AbsolutePath(t *testing.T) {
	result := run.Quick(command.Dirname("/usr/local/bin/myapp"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/usr/local/bin"})
}

func TestDirname_RelativePath(t *testing.T) {
	result := run.Quick(command.Dirname("relative/path/file.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"relative/path"})
}

func TestDirname_SingleLevel(t *testing.T) {
	result := run.Quick(command.Dirname("filename"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

// ==============================================================================
// Test Special Characters
// ==============================================================================

func TestDirname_Spaces(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/file with spaces.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/path/to"})
}

func TestDirname_SpecialChars(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/file-name_v2.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/path/to"})
}

// ==============================================================================
// Test Zero Terminator Flag
// ==============================================================================

func TestDirname_ZeroTerminator(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/file.txt", command.Zero))

	assertion.NoError(t, result.Err)
	// Should end with null byte instead of newline
	assertion.True(t, strings.HasSuffix(result.Stdout[0], "\x00"), "ends with null")
	assertion.True(t, strings.HasPrefix(result.Stdout[0], "/path/to"), "starts with dirname")
}

func TestDirname_ZeroTerminatorMultiple(t *testing.T) {
	result := run.Quick(command.Dirname("/a/file1.txt", "/b/file2.txt", command.Zero))

	assertion.NoError(t, result.Err)
	// Each output should be null-terminated
	fullOutput := strings.Join(result.Stdout, "")
	parts := strings.Split(fullOutput, "\x00")
	assertion.True(t, len(parts) >= 2, "at least 2 parts")
}

// ==============================================================================
// Test No Arguments
// ==============================================================================

func TestDirname_NoArguments(t *testing.T) {
	result := run.Quick(command.Dirname())

	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestDirname_OutputError(t *testing.T) {
	result := run.Command(command.Dirname("/path/to/file")).
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

func TestDirname_OutputError_ZeroFlag(t *testing.T) {
	result := run.Command(command.Dirname("/path/to/file", command.Zero)).
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestDirname_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"simple file", "file.txt", "."},
		{"with path", "/usr/local/bin/app", "/usr/local/bin"},
		{"trailing slash", "/path/to/dir/", "/path/to"},
		{"root", "/", "/"},
		{"current dir", ".", "."},
		{"parent dir", "..", "."},
		{"relative", "sub/dir/file", "sub/dir"},
		{"two levels", "/a/b", "/a"},
		{"three levels", "/a/b/c", "/a/b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Quick(command.Dirname(tt.path))

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, []string{tt.expected})
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestDirname_ScriptPath(t *testing.T) {
	result := run.Quick(command.Dirname("/usr/local/bin/deploy.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/usr/local/bin"})
}

func TestDirname_LogFile(t *testing.T) {
	result := run.Quick(command.Dirname("/var/log/app.log"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/var/log"})
}

func TestDirname_ConfigFile(t *testing.T) {
	result := run.Quick(command.Dirname("/etc/nginx/nginx.conf"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/etc/nginx"})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestDirname_EmptyPath(t *testing.T) {
	result := run.Quick(command.Dirname(""))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

func TestDirname_OnlySlashes(t *testing.T) {
	result := run.Quick(command.Dirname("///"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/"})
}

func TestDirname_HiddenFile(t *testing.T) {
	result := run.Quick(command.Dirname("/home/user/.bashrc"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/home/user"})
}

func TestDirname_DeepPath(t *testing.T) {
	result := run.Quick(command.Dirname("/a/b/c/d/e/f/g/file.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/a/b/c/d/e/f/g"})
}

// ==============================================================================
// Test Unicode Paths
// ==============================================================================

func TestDirname_Unicode(t *testing.T) {
	result := run.Quick(command.Dirname("/path/to/ファイル.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/path/to"})
}

func TestDirname_UnicodeDir(t *testing.T) {
	result := run.Quick(command.Dirname("/路径/文件.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/路径"})
}

// ==============================================================================
// Test With Basename Comparison
// ==============================================================================

func TestDirname_ComplementsBasename(t *testing.T) {
	// dirname gets directory part, basename gets file part
	path := "/usr/local/bin/script.sh"

	result := run.Quick(command.Dirname(path))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/usr/local/bin"})
}

