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

func (mock *MockLogger) Trace(err error) {
	mock.LastTrace = err
}

func (mock *MockLogger) Debug(err error) {
	mock.LastDebug = err
}

func (mock *MockLogger) Info(err error) {
	mock.LastInfo = err
}

func (mock *MockLogger) Warn(err error) {
	mock.LastWarn = err
}

func (mock *MockLogger) Error(err error) {
	mock.LastError = err
}

func (mock *MockLogger) Fatal(err error) {
	mock.LastFatal = err
}

func (mock *MockLogger) Panic(err error) {
	mock.LastPanic = err
}
