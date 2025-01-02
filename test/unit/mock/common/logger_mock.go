package common

type MockLogger struct {
	LastTrace error
	LastDebug error
	LastInfo  error
	LastWarn  error
	LastError error
	LastFatal error
	LastPanic error
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (mockLogger *MockLogger) Trace(err error) {
	mockLogger.LastTrace = err
}

func (mockLogger *MockLogger) Debug(err error) {
	mockLogger.LastDebug = err
}

func (mockLogger *MockLogger) Info(err error) {
	mockLogger.LastInfo = err
}

func (mockLogger *MockLogger) Warn(err error) {
	mockLogger.LastWarn = err
}

func (mockLogger *MockLogger) Error(err error) {
	mockLogger.LastError = err
}

func (mockLogger *MockLogger) Fatal(err error) {
	mockLogger.LastFatal = err
}

func (mockLogger *MockLogger) Panic(err error) {
	mockLogger.LastPanic = err
}
