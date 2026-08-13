package logger

import "time"

func AuthRequestRejected(reason, method, path string, statusCode int, duration time.Duration, message string) {
	Warn("auth_request_rejected", "request rejected due to authentication failure",
		Field{Key: "reason", Value: reason},
		Field{Key: "method", Value: method},
		Field{Key: "path", Value: path},
		Field{Key: "status_code", Value: statusCode},
		Field{Key: "duration", Value: duration.String()},
		Field{Key: "message", Value: message},
	)
}

func RequestHandled(method, path string, statusCode int, duration time.Duration, message string) {
	Info("request_handled", "request handled successfully",
		Field{Key: "method", Value: method},
		Field{Key: "path", Value: path},
		Field{Key: "status_code", Value: statusCode},
		Field{Key: "duration", Value: duration.String()},
		Field{Key: "message", Value: message},
	)
}