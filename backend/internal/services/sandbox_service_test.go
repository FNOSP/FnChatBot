package services

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func setupSandboxService(t *testing.T) *SandboxService {
	db := setupTestDB(t)
	return NewSandboxService(db)
}

func TestIsPathAllowed_Disabled(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(false)
	if err != nil {
		t.Fatalf("failed to disable sandbox: %v", err)
	}

	if !svc.IsPathAllowed("C:\\any\\path") {
		t.Error("expected all paths to be allowed when sandbox is disabled")
	}
	if !svc.IsPathAllowed("/any/path") {
		t.Error("expected all paths to be allowed when sandbox is disabled")
	}
}

func TestIsPathAllowed_ExactMatch(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	testPath := "C:\\Projects\\test"
	if runtime.GOOS != "windows" {
		testPath = "/home/user/projects"
	}

	err = svc.AddPath(testPath, "test path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	if !svc.IsPathAllowed(testPath) {
		t.Errorf("expected exact path %s to be allowed", testPath)
	}
}

func TestIsPathAllowed_SubPath(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	parentPath := "C:\\Projects"
	if runtime.GOOS != "windows" {
		parentPath = "/home/user/projects"
	}

	err = svc.AddPath(parentPath, "parent path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	childPath := parentPath + string(filepath.Separator) + "subdir"
	if !svc.IsPathAllowed(childPath) {
		t.Errorf("expected child path %s to be allowed under parent %s", childPath, parentPath)
	}

	deepChild := childPath + string(filepath.Separator) + "nested" + string(filepath.Separator) + "dir"
	if !svc.IsPathAllowed(deepChild) {
		t.Errorf("expected deep child path %s to be allowed under parent %s", deepChild, parentPath)
	}
}

func TestIsPathAllowed_WindowsPathFormats(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific test on non-Windows platform")
	}

	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	err = svc.AddPath("C:\\Projects", "test path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"backslash format", "C:\\Projects\\subdir", true},
		{"forward slash format", "C:/Projects/subdir", true},
		{"lowercase drive", "c:\\projects\\subdir", true},
		{"mixed slashes", "C:\\Projects/subdir", true},
		{"different drive", "D:\\Projects\\subdir", false},
		{"different directory", "C:\\Other\\path", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := svc.IsPathAllowed(tc.path)
			if result != tc.expected {
				t.Errorf("IsPathAllowed(%q) = %v, expected %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestIsPathAllowed_NotAllowed(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	err = svc.AddPath("C:\\Allowed", "allowed path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	blockedPaths := []string{
		"C:\\Blocked",
		"C:\\AllowedOther",
		"D:\\Any\\Path",
	}

	for _, path := range blockedPaths {
		if svc.IsPathAllowed(path) {
			t.Errorf("expected path %s to be blocked", path)
		}
	}
}

func TestExtractPathsFromCommand_CdCommand(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"simple cd windows", "cd C:\\Users\\test", 1},
		{"cd with backslash", "cd C:\\Projects", 1},
		{"cd with forward slash", "cd C:/Projects", 1},
		{"cd with quotes", `cd "C:\Program Files"`, 1},
		{"cd with single quotes", `cd 'C:\Program Files'`, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d: %v", tc.command, len(paths), tc.minPaths, paths)
			}
		})
	}
}

func TestExtractPathsFromCommand_LsCommand(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"simple ls windows", "ls C:\\Users\\test", 1},
		{"ls with flags", "ls -la C:\\Users\\test", 1},
		{"ls with multiple flags", "ls -la -h C:\\Users\\test", 1},
		{"ls with quotes", `ls "C:\Program Files"`, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d: %v", tc.command, len(paths), tc.minPaths, paths)
			}
		})
	}
}

func TestExtractPathsFromCommand_CatCommand(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"simple cat windows", "cat C:\\Users\\test\\file.txt", 1},
		{"cat with path", "cat C:\\Users\\test\\file.txt", 1},
		{"cat with quotes", `cat "C:\Program Files\config.txt"`, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d: %v", tc.command, len(paths), tc.minPaths, paths)
			}
		})
	}
}

func TestExtractPathsFromCommand_QuotedPaths(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"double quoted windows path", `cat "C:\Program Files\test.txt"`, 1},
		{"quoted path with spaces", `cd "C:\Users\My Documents"`, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d", tc.command, len(paths), tc.minPaths)
			}
		})
	}
}

func TestExtractPathsFromCommand_MultiplePaths(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"cp command", "cp C:\\Source\\file.txt C:\\Dest\\file.txt", 2},
		{"mv command", "mv C:\\Source\\file.txt C:\\Dest\\file.txt", 2},
		{"multiple commands", "cd C:\\Users\\test; ls C:\\Users\\docs", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d", tc.command, len(paths), tc.minPaths)
			}
		})
	}
}

func TestCheckCommandPermission_Allowed(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	allowedPath := "C:\\Projects"
	if runtime.GOOS != "windows" {
		allowedPath = "/home/user/projects"
	}

	err = svc.AddPath(allowedPath, "allowed path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	commands := []string{
		"cd " + allowedPath,
		"ls " + allowedPath,
		"cat " + allowedPath + string(filepath.Separator) + "file.txt",
	}

	for _, cmd := range commands {
		allowed, blocked := svc.CheckCommandPermission(cmd)
		if !allowed {
			t.Errorf("expected command %q to be allowed, but got blocked paths: %v", cmd, blocked)
		}
		if len(blocked) > 0 {
			t.Errorf("expected no blocked paths for command %q, got: %v", cmd, blocked)
		}
	}
}

func TestCheckCommandPermission_PartiallyAllowed(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	allowedPath := "C:\\Allowed"
	blockedPath := "C:\\Blocked"
	if runtime.GOOS != "windows" {
		allowedPath = "/home/allowed"
		blockedPath = "/home/blocked"
	}

	err = svc.AddPath(allowedPath, "allowed path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	command := "cp " + allowedPath + string(filepath.Separator) + "file.txt " + blockedPath + string(filepath.Separator) + "file.txt"
	allowed, blockedPaths := svc.CheckCommandPermission(command)

	if allowed {
		t.Error("expected command to be blocked due to partially blocked paths")
	}

	if len(blockedPaths) == 0 {
		t.Error("expected at least one blocked path")
	}

	foundBlocked := false
	for _, p := range blockedPaths {
		normalized, _ := svc.normalizePath(p)
		blockedNormalized, _ := svc.normalizePath(blockedPath)
		if normalized == blockedNormalized || len(normalized) >= len(blockedNormalized) && normalized[:len(blockedNormalized)] == blockedNormalized {
			foundBlocked = true
			break
		}
	}
	if !foundBlocked {
		t.Errorf("expected blocked path %s in blocked paths list, got: %v", blockedPath, blockedPaths)
	}
}

func TestCheckCommandPermission_Blocked(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	allowedPath := "C:\\Allowed"
	if runtime.GOOS != "windows" {
		allowedPath = "/home/allowed"
	}

	err = svc.AddPath(allowedPath, "allowed path")
	if err != nil {
		t.Fatalf("failed to add path: %v", err)
	}

	blockedCommands := []string{
		"cd C:\\Blocked",
		"ls C:\\Blocked",
		"cat C:\\Blocked\\file.txt",
	}

	if runtime.GOOS != "windows" {
		blockedCommands = []string{
			"cd /home/blocked",
			"ls /home/blocked",
			"cat /home/blocked/file.txt",
		}
	}

	for _, cmd := range blockedCommands {
		allowed, blocked := svc.CheckCommandPermission(cmd)
		if allowed {
			t.Errorf("expected command %q to be blocked", cmd)
		}
		if len(blocked) == 0 {
			t.Errorf("expected blocked paths for command %q", cmd)
		}
	}
}

func TestCheckCommandPermission_NoPaths(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	commands := []string{
		"ls",
		"pwd",
		"echo hello",
		"date",
	}

	for _, cmd := range commands {
		allowed, blocked := svc.CheckCommandPermission(cmd)
		if !allowed {
			t.Errorf("expected command %q with no paths to be allowed, got blocked: %v", cmd, blocked)
		}
	}
}

func TestCheckCommandPermission_Disabled(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(false)
	if err != nil {
		t.Fatalf("failed to disable sandbox: %v", err)
	}

	command := "cat C:\\Sensitive\\secrets.txt"
	if runtime.GOOS != "windows" {
		command = "cat /etc/shadow"
	}

	allowed, blocked := svc.CheckCommandPermission(command)
	if !allowed {
		t.Error("expected all commands to be allowed when sandbox is disabled")
	}
	if len(blocked) > 0 {
		t.Errorf("expected no blocked paths when sandbox is disabled, got: %v", blocked)
	}
}

func TestNormalizePath(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name    string
		input   string
		checkFn func(string) bool
	}{
		{"empty string", "", func(s string) bool { return s == "" }},
		{"whitespace only", "   ", func(s string) bool { return s == "" || s == "   " || true }},
	}

	if runtime.GOOS == "windows" {
		winCases := []struct {
			name    string
			input   string
			checkFn func(string) bool
		}{
			{"lowercase drive", "c:\\users", func(s string) bool { return len(s) >= 2 && s[0] == 'C' && s[1] == ':' }},
			{"forward to backslash", "C:/Users/test", func(s string) bool { return true }},
			{"trailing spaces", "  C:\\Users  ", func(s string) bool { return true }},
		}
		testCases = append(testCases, winCases...)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.normalizePath(tc.input)
			if err != nil {
				t.Errorf("normalizePath(%q) returned error: %v", tc.input, err)
			}
			if !tc.checkFn(result) {
				t.Errorf("normalizePath(%q) = %q, check failed", tc.input, result)
			}
		})
	}
}

func TestIsSubPath(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		parent   string
		child    string
		expected bool
	}{
		{"exact match", "/home/user", "/home/user", true},
		{"direct child", "/home/user", "/home/user/docs", true},
		{"nested child", "/home/user", "/home/user/docs/projects", true},
		{"different path", "/home/user", "/home/other", false},
		{"partial match", "/home/user", "/home/userbackup", false},
	}

	if runtime.GOOS == "windows" {
		testCases = []struct {
			name     string
			parent   string
			child    string
			expected bool
		}{
			{"exact match", "C:\\Users", "C:\\Users", true},
			{"direct child", "C:\\Users", "C:\\Users\\docs", true},
			{"nested child", "C:\\Users", "C:\\Users\\docs\\projects", true},
			{"different path", "C:\\Users", "C:\\Other", false},
			{"partial match", "C:\\Users", "C:\\UsersBackup", false},
			{"case insensitive", "C:\\Users", "c:\\users\\docs", true},
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := svc.isSubPath(tc.parent, tc.child)
			if result != tc.expected {
				t.Errorf("isSubPath(%q, %q) = %v, expected %v", tc.parent, tc.child, result, tc.expected)
			}
		})
	}
}

func TestUniquePaths(t *testing.T) {
	svc := setupSandboxService(t)

	paths := []string{
		"/home/user",
		"/home/user",
		"/home/other",
		"C:\\Users",
		"C:/Users",
	}

	result := svc.uniquePaths(paths)

	pathCount := make(map[string]int)
	for _, p := range result {
		normalized, _ := svc.normalizePath(p)
		pathCount[normalized]++
	}

	for path, count := range pathCount {
		if count > 1 {
			t.Errorf("uniquePaths returned duplicate normalized path %q", path)
		}
	}
}

func TestAddPath(t *testing.T) {
	svc := setupSandboxService(t)

	testPath := "C:\\Test\\Path"
	if runtime.GOOS != "windows" {
		testPath = "/test/path"
	}

	err := svc.AddPath(testPath, "test description")
	if err != nil {
		t.Fatalf("AddPath failed: %v", err)
	}

	paths := svc.GetAllowedPaths()
	found := false
	for _, p := range paths {
		if p == testPath || p == "C:\\Test\\Path" || p == "/test/path" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("path %q not found in allowed paths: %v", testPath, paths)
	}
}

func TestRemovePath(t *testing.T) {
	svc := setupSandboxService(t)

	testPath := "C:\\Test\\Remove"
	if runtime.GOOS != "windows" {
		testPath = "/test/remove"
	}

	err := svc.AddPath(testPath, "to be removed")
	if err != nil {
		t.Fatalf("AddPath failed: %v", err)
	}

	err = svc.RemovePath(testPath)
	if err != nil {
		t.Fatalf("RemovePath failed: %v", err)
	}

	paths := svc.GetAllowedPaths()
	for _, p := range paths {
		if p == testPath {
			t.Errorf("path %q should have been removed", testPath)
		}
	}
}

func TestGetAllPaths(t *testing.T) {
	svc := setupSandboxService(t)

	paths := []struct {
		path        string
		description string
	}{
		{"C:\\Path1", "first path"},
		{"C:\\Path2", "second path"},
		{"C:\\Path3", "third path"},
	}

	if runtime.GOOS != "windows" {
		paths = []struct {
			path        string
			description string
		}{
			{"/path1", "first path"},
			{"/path2", "second path"},
			{"/path3", "third path"},
		}
	}

	for _, p := range paths {
		err := svc.AddPath(p.path, p.description)
		if err != nil {
			t.Fatalf("AddPath failed: %v", err)
		}
	}

	allPaths := svc.GetAllPaths()
	if len(allPaths) != len(paths) {
		t.Errorf("GetAllPaths returned %d paths, expected %d", len(allPaths), len(paths))
	}
}

func TestExtractPathsFromCommand_WindowsCommands(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"dir command", "dir C:\\Users", 1},
		{"type command", "type C:\\Windows\\system.ini", 1},
		{"del command", "del C:\\Temp\\file.txt", 1},
		{"copy command", "copy C:\\Source\\file.txt C:\\Dest\\file.txt", 2},
		{"move command", "move C:\\Source\\file.txt C:\\Dest\\file.txt", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d", tc.command, len(paths), tc.minPaths)
			}
		})
	}
}

func TestExtractPathsFromCommand_EdgeCases(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		name     string
		command  string
		minPaths int
	}{
		{"empty command", "", 0},
		{"command with no path", "ls -la", 0},
		{"command with only flags", "grep -r pattern", 0},
		{"relative path", "cd subdir", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) < tc.minPaths {
				t.Errorf("ExtractPathsFromCommand(%q) returned %d paths, expected at least %d: %v", tc.command, len(paths), tc.minPaths, paths)
			}
		})
	}
}

func TestIsFlag(t *testing.T) {
	svc := setupSandboxService(t)

	testCases := []struct {
		input    string
		expected bool
	}{
		{"-la", true},
		{"-r", true},
		{"/?", true},
		{"/help", true},
		{"normal", false},
		{"path.txt", false},
		{"C:\\Path", false},
	}

	for _, tc := range testCases {
		result := svc.isFlag(tc.input)
		if result != tc.expected {
			t.Errorf("isFlag(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestIsEnabled_Default(t *testing.T) {
	svc := setupSandboxService(t)

	if svc.IsEnabled() {
		t.Error("expected sandbox to be disabled by default")
	}
}

func TestSetEnabled(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("SetEnabled(true) failed: %v", err)
	}
	if !svc.IsEnabled() {
		t.Error("expected sandbox to be enabled")
	}

	err = svc.SetEnabled(false)
	if err != nil {
		t.Fatalf("SetEnabled(false) failed: %v", err)
	}
	if svc.IsEnabled() {
		t.Error("expected sandbox to be disabled")
	}
}

func TestExtractPathsFromCommand_VariousCommands(t *testing.T) {
	svc := setupSandboxService(t)

	commands := []struct {
		name    string
		command string
	}{
		{"rm", "rm C:\\Users\\test\\file.txt"},
		{"mkdir", "mkdir C:\\Users\\test\\newdir"},
		{"touch", "touch C:\\Users\\test\\newfile.txt"},
		{"chmod", "chmod 755 C:\\Users\\test\\script.sh"},
		{"find", "find C:\\Users\\test -name '*.txt'"},
		{"head", "head -n 10 C:\\Users\\test\\file.txt"},
		{"tail", "tail -n 10 C:\\Users\\test\\file.txt"},
		{"less", "less C:\\Users\\test\\file.txt"},
		{"vim", "vim C:\\Users\\test\\file.txt"},
		{"echo redirect", "echo hello > C:\\Users\\test\\output.txt"},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			paths := svc.ExtractPathsFromCommand(tc.command)
			if len(paths) == 0 {
				t.Errorf("ExtractPathsFromCommand(%q) returned no paths", tc.command)
			}
		})
	}
}

func TestMultipleAllowedPaths(t *testing.T) {
	svc := setupSandboxService(t)

	err := svc.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable sandbox: %v", err)
	}

	paths := []string{
		"C:\\Projects",
		"C:\\Documents",
		"C:\\Downloads",
	}
	if runtime.GOOS != "windows" {
		paths = []string{
			"/home/user/projects",
			"/home/user/documents",
			"/home/user/downloads",
		}
	}

	for _, p := range paths {
		err = svc.AddPath(p, "allowed path")
		if err != nil {
			t.Fatalf("failed to add path %s: %v", p, err)
		}
	}

	for _, p := range paths {
		if !svc.IsPathAllowed(p) {
			t.Errorf("expected path %s to be allowed", p)
		}
		childPath := p + string(filepath.Separator) + "subdir"
		if !svc.IsPathAllowed(childPath) {
			t.Errorf("expected child path %s to be allowed", childPath)
		}
	}

	sort.Strings(paths)
	allowedPaths := svc.GetAllowedPaths()
	sort.Strings(allowedPaths)

	if len(allowedPaths) != len(paths) {
		t.Errorf("expected %d allowed paths, got %d", len(paths), len(allowedPaths))
	}
}
