package logger

import "github.com/rs/zerolog"

type Field struct {
	Key   string
	Value any
}

func applyFields(evt *zerolog.Event, fields ...Field) *zerolog.Event {
	for _, field := range fields {
		switch v := field.Value.(type) {
		case string:
			evt = evt.Str(field.Key, v)
		case uint:
			evt = evt.Uint(field.Key, v)
		case uint8:
			evt = evt.Uint8(field.Key, v)
		case uint16:
			evt = evt.Uint16(field.Key, v)
		case uint32:
			evt = evt.Uint32(field.Key, v)
		case uint64:
			evt = evt.Uint64(field.Key, v)
		case int:
			evt = evt.Int(field.Key, v)
		case int8:
			evt = evt.Int8(field.Key, v)
		case int16:
			evt = evt.Int16(field.Key, v)
		case int32:
			evt = evt.Int32(field.Key, v)
		case int64:
			evt = evt.Int64(field.Key, v)
		case float32:
			evt = evt.Float32(field.Key, v)
		case float64:
			evt = evt.Float64(field.Key, v)
		case bool:
			evt = evt.Bool(field.Key, v)
		case error:
			evt = evt.Err(v)
		default:
			evt = evt.Interface(field.Key, v)
		}
	}
	return evt
}

func Info(eventName, msg string, fields ...Field) {
	applyFields(Log.Info().Str("event", eventName), fields...).Msg(msg)
}

func Warn(eventName, msg string, fields ...Field) {
	applyFields(Log.Warn().Str("event", eventName), fields...).Msg(msg)
}

func Error(eventName, msg string, fields ...Field) {
	applyFields(Log.Error().Str("event", eventName), fields...).Msg(msg)
}
