package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samverrall/review-ui/internal/git"
	"github.com/samverrall/review-ui/internal/git/testutil"
)


// Helper function to create a test model with mocked dependencies
func createTestModel(mock git.GitClient) model {
	ti := textinput.New()
	ti.Placeholder = "Enter your comment..."
	ti.CharLimit = 200
	ti.Width = 80

	// Initialize viewport with basic dimensions
	vp := viewport.New(80, 20)

	// Get changed files from mock for testing
	var changedFiles []string
	if mockClient, ok := mock.(*testutil.MockGitClient); ok {
		changedFiles = mockClient.GetChangedFilesForTest()
	}

	return model{
		gitClient:    mock,
		changedFiles: changedFiles,
		currentIndex: 0,
		diffs:        make(map[string]string),
		viewport:     vp,
		commentInput: ti,
		commentMode:  false,
		comments:     make(map[string][]string),
	}
}

func TestModelInitialization(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() git.GitClient
		expectError bool
	}{
		{
			name: "git repo check succeeds when in repo",
			setupMock: func() git.GitClient {
				return testutil.NewMockGitClient().WithIsRepo(true)
			},
			expectError: false,
		},
		{
			name: "git repo check fails when git command errors",
			setupMock: func() git.GitClient {
				return testutil.NewMockGitClient().
					WithRepoError(&mockError{"git check failed"})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()

			// Test IsGitRepo behavior
			_, err := mock.IsGitRepo()
			if err != nil && !tt.expectError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectError {
				t.Errorf("expected error but got none")
			}
		})
	}
}

func TestFileNavigation(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go", "file3.go"})

	m := createTestModel(mock)

	// Test initial state
	if m.currentIndex != 0 {
		t.Errorf("expected initial currentIndex 0, got %d", m.currentIndex)
	}

	// Test next file navigation
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.currentIndex != 1 {
		t.Errorf("expected currentIndex 1 after next, got %d", m.currentIndex)
	}

	// Test next file again
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.currentIndex != 2 {
		t.Errorf("expected currentIndex 2 after next, got %d", m.currentIndex)
	}

	// Test wraparound to beginning
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.currentIndex != 0 {
		t.Errorf("expected currentIndex 0 after wraparound, got %d", m.currentIndex)
	}

	// Test previous file navigation
	m.currentIndex = 2 // Set to last file
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = updatedModel.(model)
	if m.currentIndex != 1 {
		t.Errorf("expected currentIndex 1 after previous, got %d", m.currentIndex)
	}

	// Test wraparound to end
	m.currentIndex = 0 // Set to first file
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = updatedModel.(model)
	if m.currentIndex != 2 {
		t.Errorf("expected currentIndex 2 after wraparound to end, got %d", m.currentIndex)
	}
}

func TestCursorMovement(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go"})

	m := createTestModel(mock)

	// Mock viewport with some content
	// Note: This is a simplified test since viewport behavior is complex
	// In a real test, you'd need to properly initialize the viewport

	// Test initial cursor position
	if m.cursorLine != 0 {
		t.Errorf("expected initial cursorLine 0, got %d", m.cursorLine)
	}

	// Test cursor movement keys are handled (without full viewport setup)
	// We test that the keys don't cause errors and the model updates
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if updatedModel == nil {
		t.Errorf("expected model to be returned after 'j' key")
	}

	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if updatedModel == nil {
		t.Errorf("expected model to be returned after 'k' key")
	}
}

func TestSelectionMode(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go"})

	m := createTestModel(mock)

	// Test initial state - not in selection mode
	if m.selectionMode {
		t.Errorf("expected initial selectionMode false, got %v", m.selectionMode)
	}

	// Enter selection mode
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")})
	m = updatedModel.(model)
	if !m.selectionMode {
		t.Errorf("expected selectionMode true after 'v', got %v", m.selectionMode)
	}
	if m.selectionStart != m.cursorLine {
		t.Errorf("expected selectionStart to equal cursorLine")
	}

	// Exit selection mode
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")})
	m = updatedModel.(model)
	if m.selectionMode {
		t.Errorf("expected selectionMode false after second 'v', got %v", m.selectionMode)
	}

	// Test selection mode exit with escape
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("v")}) // Enter selection
	m = updatedModel.(model)
	if !m.selectionMode {
		t.Errorf("expected selectionMode true")
	}
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")})
	m = updatedModel.(model)
	if m.selectionMode {
		t.Errorf("expected selectionMode false after 'esc', got %v", m.selectionMode)
	}
}

func TestCommentFunctionality(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go"})

	m := createTestModel(mock)

	// Test initial state
	if m.commentMode {
		t.Errorf("expected initial commentMode false, got %v", m.commentMode)
	}
	if len(m.comments) != 0 {
		t.Errorf("expected no initial comments, got %d", len(m.comments))
	}

	// Enter comment mode
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = updatedModel.(model)
	if !m.commentMode {
		t.Errorf("expected commentMode true after 'c', got %v", m.commentMode)
	}

	// Test adding a comment (simplified - normally this would be done through text input)
	// This is hard to test fully without mocking the text input behavior
	// In practice, you'd test the comment storage logic separately

	// Test comment key generation
	key := m.getCommentKey(5)
	expectedKey := "file1.go:5"
	if key != expectedKey {
		t.Errorf("expected comment key %s, got %s", expectedKey, key)
	}

	// Test range comment key generation
	rangeKey := m.getCommentKeyForRange(5, 10)
	expectedRangeKey := "file1.go:5-10"
	if rangeKey != expectedRangeKey {
		t.Errorf("expected range comment key %s, got %s", expectedRangeKey, rangeKey)
	}

	// Test range comment key generation with swapped order
	rangeKey2 := m.getCommentKeyForRange(10, 5)
	if rangeKey2 != expectedRangeKey {
		t.Errorf("expected range comment key %s even with swapped args, got %s", expectedRangeKey, rangeKey2)
	}
}

func TestExportFunctionality(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go"})

	m := createTestModel(mock)

	// Test export with no comments
	exported := m.exportComments()
	expected := "No comments to export."
	if exported != expected {
		t.Errorf("expected '%s', got '%s'", expected, exported)
	}

	// Add some test comments
	m.comments["file1.go:5"] = []string{"This line needs improvement"}
	m.comments["file1.go:10-15"] = []string{"This block could be refactored"}
	m.comments["file2.go:20"] = []string{"Consider error handling"}

	// Test export with comments
	exported = m.exportComments()
	if exported == "No comments to export." {
		t.Errorf("expected comments to be exported, but got no comments message")
	}

	// Check that the export contains expected content
	if !contains(exported, "# Code Review Comments") {
		t.Errorf("expected export to contain header")
	}
	if !contains(exported, "file1.go") {
		t.Errorf("expected export to contain file1.go")
	}
	if !contains(exported, "file2.go") {
		t.Errorf("expected export to contain file2.go")
	}
	if !contains(exported, "This line needs improvement") {
		t.Errorf("expected export to contain the comment")
	}
}

func TestCommentRangeHandling(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go"})

	m := createTestModel(mock)

	// Test getSelectionRange with normal order
	m.selectionStart = 5
	m.cursorLine = 10
	start, end := m.getSelectionRange()
	if start != 5 || end != 10 {
		t.Errorf("expected range 5-10, got %d-%d", start, end)
	}

	// Test getSelectionRange with reversed order
	m.selectionStart = 10
	m.cursorLine = 5
	start, end = m.getSelectionRange()
	if start != 5 || end != 10 {
		t.Errorf("expected range 5-10 (normalized), got %d-%d", start, end)
	}

	// Test range comment key generation
	key := m.getCommentKeyForRange(5, 10)
	expectedKey := "file1.go:5-10"
	if key != expectedKey {
		t.Errorf("expected key %s, got %s", expectedKey, key)
	}

	// Test range comment key generation with swapped args (should normalize)
	key2 := m.getCommentKeyForRange(10, 5)
	if key2 != expectedKey {
		t.Errorf("expected key %s even with swapped args, got %s", expectedKey, key2)
	}
}

func TestCommentStorageAndRetrieval(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go"})

	m := createTestModel(mock)

	// Test adding single line comments
	m.comments["file1.go:5"] = []string{"First comment"}
	m.comments["file1.go:5"] = append(m.comments["file1.go:5"], "Second comment")
	m.comments["file2.go:10"] = []string{"Comment on different file"}

	// Test retrieving comments
	if len(m.comments["file1.go:5"]) != 2 {
		t.Errorf("expected 2 comments for line 5, got %d", len(m.comments["file1.go:5"]))
	}

	if m.comments["file1.go:5"][0] != "First comment" {
		t.Errorf("expected first comment 'First comment', got '%s'", m.comments["file1.go:5"][0])
	}

	if len(m.comments["file2.go:10"]) != 1 {
		t.Errorf("expected 1 comment for file2 line 10, got %d", len(m.comments["file2.go:10"]))
	}
}

func TestExportWithMultipleFilesAndComments(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go"})

	m := createTestModel(mock)

	// Add comments to multiple files and ranges
	m.comments["file1.go:5"] = []string{"Single line comment"}
	m.comments["file1.go:10-15"] = []string{"Range comment on file1"}
	m.comments["file2.go:20"] = []string{"Comment on file2 line 20"}
	m.comments["file2.go:25"] = []string{"Another comment on file2"}

	exported := m.exportComments()

	// Check that export contains expected elements
	if !contains(exported, "# Code Review Comments") {
		t.Errorf("expected header in export")
	}

	if !contains(exported, "## File: file1.go") {
		t.Errorf("expected file1.go section")
	}

	if !contains(exported, "## File: file2.go") {
		t.Errorf("expected file2.go section")
	}

	if !contains(exported, "Line 6") {
		t.Errorf("expected line 6 reference (1-indexed)")
	}

	if !contains(exported, "Lines 11-16") {
		t.Errorf("expected range 11-16 reference (1-indexed)")
	}

	if !contains(exported, "Single line comment") {
		t.Errorf("expected single line comment content")
	}

	if !contains(exported, "Range comment on file1") {
		t.Errorf("expected range comment content")
	}
}

func TestFileListMode(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go", "file3.go"})

	m := createTestModel(mock)

	// Test entering file list mode
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("tab")})
	m = updatedModel.(model)
	if !m.fileListMode {
		t.Errorf("expected fileListMode true after tab, got %v", m.fileListMode)
	}

	// Test navigating in file list
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = updatedModel.(model)
	if m.fileListCursor != 1 {
		t.Errorf("expected fileListCursor 1 after j, got %d", m.fileListCursor)
	}

	// Test selecting file from list
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("enter")})
	m = updatedModel.(model)
	if m.fileListMode {
		t.Errorf("expected fileListMode false after enter, got %v", m.fileListMode)
	}
	if m.currentIndex != 1 {
		t.Errorf("expected currentIndex 1 after selection, got %d", m.currentIndex)
	}

	// Test exiting file list mode with escape
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("tab")}) // Enter file list
	m = updatedModel.(model)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")})
	m = updatedModel.(model)
	if m.fileListMode {
		t.Errorf("expected fileListMode false after esc, got %v", m.fileListMode)
	}
}

func TestQuitFunctionality(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go"})

	m := createTestModel(mock)

	// Test quit with 'q'
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Errorf("expected quit command, got nil")
	}
	// In Bubble Tea, tea.Quit is a special command that causes the program to exit

	// Test quit with Ctrl+C
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Errorf("expected quit command for Ctrl+C, got nil")
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   containsMiddle(s, substr)
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestEdgeCases(t *testing.T) {
	// Test with empty file list
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{})

	m := createTestModel(mock)

	// Test navigation with empty file list
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.currentIndex != 0 {
		t.Errorf("expected currentIndex 0 with empty file list, got %d", m.currentIndex)
	}

	// Test with single file
	mock = testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"single.go"})
	m = createTestModel(mock)

	// Test navigation wraparound with single file
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.currentIndex != 0 {
		t.Errorf("expected currentIndex 0 after next with single file, got %d", m.currentIndex)
	}

	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = updatedModel.(model)
	if m.currentIndex != 0 {
		t.Errorf("expected currentIndex 0 after previous with single file, got %d", m.currentIndex)
	}
}

func TestCommentKeyGeneration(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go", "file2.go"})

	m := createTestModel(mock)

	// Test comment key generation bounds checking
	m.currentIndex = -1
	key := m.getCommentKey(5)
	if key != "" {
		t.Errorf("expected empty key for invalid currentIndex, got %s", key)
	}

	m.currentIndex = 10 // Beyond array bounds
	key = m.getCommentKey(5)
	if key != "" {
		t.Errorf("expected empty key for out of bounds currentIndex, got %s", key)
	}

	// Reset to valid index
	m.currentIndex = 0

	// Test valid key generation
	key = m.getCommentKey(5)
	expected := "file1.go:5"
	if key != expected {
		t.Errorf("expected key %s, got %s", expected, key)
	}

	// Test range key bounds checking
	m.currentIndex = -1
	rangeKey := m.getCommentKeyForRange(5, 10)
	if rangeKey != "" {
		t.Errorf("expected empty range key for invalid currentIndex, got %s", rangeKey)
	}
}

func TestStatusMessages(t *testing.T) {
	mock := testutil.NewMockGitClient().
		WithIsRepo(true).
		WithChangedFiles([]string{"file1.go"})

	m := createTestModel(mock)

	// Test that status messages are cleared on navigation
	m.statusMessage = "Some status"
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = updatedModel.(model)
	if m.statusMessage != "" {
		t.Errorf("expected status message to be cleared on navigation, got %s", m.statusMessage)
	}

	// Test that status messages are cleared on file switch
	m.statusMessage = "Some status"
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = updatedModel.(model)
	if m.statusMessage != "" {
		t.Errorf("expected status message to be cleared on file switch, got %s", m.statusMessage)
	}
}

// mockError implements error interface for testing
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}