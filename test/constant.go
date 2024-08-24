package test

const (
	// Error Message Formats.
	ExpectedErrorMessageFormat = "location: %s notification: %s"
	ExpectedStatusOk           = "Expected status OK but got %v"
	ExpectedBody               = "Expected body %v but got %v"
	EqualMessage               = "they should be equal"

	// Failure/Success Messages.
	FailureMessage      = "result should be a failure"
	NotFailureMessage   = "result should not be a failure"
	DataNotNilMessage   = "data should be nil"
	ErrorNilMessage     = "error should not be nil"
	CorrectErrorMessage = "error message should be correct"

	// Response Messages.
	ResponseBodyMismatch = "response body mismatch"
	StatusCodeMismatch   = "status code mismatch"

	// JSON Response Messages.
	AlreadyLoggedInMessage = `{"error":{"notification":"Already logged in. This action is not allowed."},"status":"fail"}`
	NotLoggedInMessage     = `{"error":{"notification":"You are not logged in."},"status":"fail"}`
	SuccessResponse        = `{"message":"success"}`
	Message                = "message"

	// Test URL.
	TestURL       = "/test"
	Localhost     = "localhost:8080"
	LocalhostTest = "http://localhost:8080/test"

	// Keys for JWT token handling.
	PublicKey  = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwZ0grQ3Z6Tkg5Z0ZYVVUxQ0lVbApadjdwZExBRFZpZVV1TnZ3eVBnV1NxNklnY0xJeGVjMDRJV1E0R1dZc0lhcWlzQzFqUVp6dWd1V1NrZ0xyNDkrCnlZWm5tYzBXd0tndFRUQllyamtzTFdQNmREVFhzUTZnZ0RDbnJ0ZWY4V3NiTTJGdWxoNVBGYUgzbmlDeGpqWnoKcU5TN2NtdnNaQUFDeHEwZUZ4ZTh0aVR3S0ZPbzhWWTNoelZ1YXhSVHhyWlRaTEZiaHc4TXRGWUxuaitjaTl0egppOXhMY2VWbGpJL0thVU8vU2lCQmFBUVQrTmJCb0dnbDdNemMvRVBSMXBjZnlMdXlPRUp1WVg3TW5aaW9jSHJkCkFrVUttNDc5MzNybnJXcVVSdEFDN2F4ZEdBWHFXNjVqSXNUcU9yVUhmQ2VLMEFqWHBzOGRNZnJaQ085L3NHMUIKb3dJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
	PrivateKey = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFNBZjRLL00wZjJBVmQKUlRVSWhTVm0vdWwwc0FOV0o1UzQyL0RJK0JaS3JvaUJ3c2pGNXpUZ2haRGdaWml3aHFxS3dMV05Cbk82QzVaSwpTQXV2ajM3SmhtZVp6UmJBcUMxTk1GaXVPU3d0WS9wME5OZXhEcUNBTUtldTE1L3hheHN6WVc2V0hrOFZvZmVlCklMR09Obk9vMUx0eWEreGtBQUxHclI0WEY3eTJKUEFvVTZqeFZqZUhOVzVyRkZQR3RsTmtzVnVIRHd5MFZndWUKUDV5TDIzT0wzRXR4NVdXTWo4cHBRNzlLSUVGb0JCUDQxc0dnYUNYc3pOejhROUhXbHgvSXU3STRRbTVoZnN5ZAptS2h3ZXQwQ1JRcWJqdjNmZXVldGFwUkcwQUx0ckYwWUJlcGJybU1peE9vNnRRZDhKNHJRQ05lbXp4MHgrdGtJCjczK3diVUdqQWdNQkFBRUNnZ0VBRWF3ZlA3ZDBYNGlqTXUwZkFGK0wvVFhZV1h4eVcyNnJRajhuN1JHTGRxOW4KUjF3bjN4ZU15SlFVMC8xWXN3b3lFY2tUdmhGYjdiMEo0YWhjYTJLczdiS0V4MW1OMzVxSGJXWnpILzRwckl3cwpTRmttQ1gxTW5sejV6Mm5QeU5ZVmpPWlhFd1RyN01zYmRsQVVBUDZ1RHZnUDZob1E0MzFvdm1WVkVlWnFkLzFPCm9OOFFxTVg5YkhtZ1htNmlsZ2p3eWFiTWJFRXZnWml4TnBZcVNaNHNLKzVCVm5tYXMrZjlFLzdFZzJBYU9SRDUKbGMyREhUdmozbFgvd3h4aUg0Nkd0bk9VOUNWVzBtVm83ZHR1S2EyaWIxY2JraUZSaEdSdGRkMVZ6RzFjbEZRegovR2JEUDFDbnJoOUZiRWRmZWlpY2ZYNVVnVjZxKzAxdk9ENmUwU2RmMFFLQmdRRDZWWFpUcjJFRzBwM2p6SUNGCjR6Slg3TmlxdkhUN3hNTWoyRlM4TlM5d1hiTXg0MFZLbGVkaGxFYlJGUG1ERVh1Qis2OFBkVEppV0NMUEEwczQKT0EyeWV2UlFJbFdLdFd0OXJ0bjFSTi9WMXpMRGVsMlpTdENxWnJUZGZHWnNOc0VKQXg4b2lDRGlMTWNGQ2p0NgpsNVJWL2IxZmxIZkUwY0FaZ0NSUVA2QmJtUUtCZ1FEV3d0MURFYWJyUGNIaCtoMTdVeXpFL2NpNEU2TUpQWXhlCmYrWjQ3V3RFbEdiSExSb2gyekUwNDlsbGYxaTFkdFR5Q2pHNTdybmVCMnIzZ0M2MnJpWmJTMW4wZ054ZG9IVy8KWktIVGtNN2ZZM0h6aEt4M05VYmRDVTZpSlROSFFvTEhUOFd0UWx2dHV3ckdQN1dza29DYjBlNkozUzc5ck1FKwo3QzhUYWdPc213S0JnUUN3ek5jck4zd2hZM01ieGYwbmtsU21BS0x0d3Zna01NMVpiWm82NnAwOGtSRFVOUjFsCkZnWTZ4b3hWY3FqZVJ1U2g0dTI2enh6c2xDN1JZaFFuK243Q0JWQ3puK3dtY1FjZjF2UWM0NjNxeTNnUTAwVnoKMUlFWE9ENlpCeGtYYUh4aEx4RThnUmdvWlZPU1hhMndZWW5rU2JjTDRFSE9nZzFZZFVZd1h4K1VVUUtCZ0ZNQwo5MnVaUXgvaXV6S1I3eHRnUnduTjN1dm9DempqSllMUmhWQncxT21wUXlEeCtndmtJZDBEeFdCS0hRdm5adUEzCnVJamFFZFlVbi9BVEIvdHN6VDYwbll5NDBuVU9OUFZKL0pNK2dmZ3ZCRGpRcTZsWVdvL05yU3RYbmI2Sm91dFAKbG1VbUpVcDY1ZXREYlFITGp4S3J6cnhUVm5xUGNCTFdVRXY4eW5iSkFvR0JBUFR2YzFLb3lCNklvSHJPcWJ3RQoxbG1jQjVoZTVIbG5ncVFYSm85cG9DRmdGaEp2RFNVb3Z3K3Zvdm01Y0VkazFsb3lCWGh6UmYveHA4TExBL1B5CjJPYmFlaEovNStjZStTOTFBUnRMQUpINlNpZnhkV0JORGk4K1NLZ2hQTEpWcUdTQnNUSE4xZG5iK3NXUjZScksKdzB3ZUxDVmdVUTZzaFF5Zmg0QTlCanhSCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"

	// Secure headers.
	ContentSecurityPolicyHeader   = "Content-Security-Policy"
	ContentSecurityPolicyValue    = "default-src 'self';"
	StrictTransportSecurityHeader = "Strict-Transport-Security"
	StrictTransportSecurityValue  = "max-age=31536000; includeSubDomains"
	XContentTypeOptionsHeader     = "X-Content-Type-Options"
	XContentTypeOptionsValue      = "nosniff"
)
