package testutil

// MockGitClient allows us to mock git operations for testing
type MockGitClient struct {
	isRepo       bool
	changedFiles []string
	fileDiffs    map[string]string
	repoError    error
	filesError   error
	diffError    error
}

// NewMockGitClient creates a new mock git client with default values
func NewMockGitClient() *MockGitClient {
	return &MockGitClient{
		fileDiffs: make(map[string]string),
	}
}

// WithIsRepo sets the mock to return the specified repo status
func (m *MockGitClient) WithIsRepo(isRepo bool) *MockGitClient {
	m.isRepo = isRepo
	return m
}

// WithChangedFiles sets the mock to return the specified changed files
func (m *MockGitClient) WithChangedFiles(files []string) *MockGitClient {
	m.changedFiles = files
	return m
}

// WithFileDiff sets a diff for a specific file
func (m *MockGitClient) WithFileDiff(filename, diff string) *MockGitClient {
	if m.fileDiffs == nil {
		m.fileDiffs = make(map[string]string)
	}
	m.fileDiffs[filename] = diff
	return m
}

// WithRepoError sets the mock to return the specified error for repo checks
func (m *MockGitClient) WithRepoError(err error) *MockGitClient {
	m.repoError = err
	return m
}

// WithFilesError sets the mock to return the specified error for file operations
func (m *MockGitClient) WithFilesError(err error) *MockGitClient {
	m.filesError = err
	return m
}

// WithDiffError sets the mock to return the specified error for diff operations
func (m *MockGitClient) WithDiffError(err error) *MockGitClient {
	m.diffError = err
	return m
}

// GetChangedFilesForTest returns the configured changed files for testing
func (m *MockGitClient) GetChangedFilesForTest() []string {
	return m.changedFiles
}

func (m *MockGitClient) IsGitRepo() (bool, error) {
	return m.isRepo, m.repoError
}

func (m *MockGitClient) GetChangedFiles() ([]string, error) {
	return m.changedFiles, m.filesError
}

func (m *MockGitClient) GetFileDiff(filename string) (string, error) {
	if diff, exists := m.fileDiffs[filename]; exists {
		return diff, m.diffError
	}
	return "", m.diffError
}
