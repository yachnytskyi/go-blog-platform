package mock

// MockLogger is a mock implementation of a logger for testing purposes.
type MockLogger struct {
	LastTrace error
	LastDebug error
	LastInfo  error
	LastWarn  error
	LastError error
	LastFatal error
	LastPanic error
}

func (mock *MockLogger) Trace(data error) {
	mock.LastTrace = data
}

func (mock *MockLogger) Debug(data error) {
	mock.LastDebug = data
}

func (mock *MockLogger) Info(data error) {
	mock.LastInfo = data
}

func (mock *MockLogger) Warn(data error) {
	mock.LastWarn = data
}

func (mock *MockLogger) Error(data error) {
	mock.LastError = data
}

func (mock *MockLogger) Fatal(data error) {
	mock.LastFatal = data
}

func (mock *MockLogger) Panic(data error) {
	mock.LastPanic = data
}
