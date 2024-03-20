package redisfsm

import "time"

type OptionFunc func(settings *StorageSettings)

func WithPrefix(prefix string) OptionFunc {
	return func(settings *StorageSettings) {
		settings.Prefix = prefix
	}
}

func WithTTLForStates(ttlForStates time.Duration) OptionFunc {
	return func(settings *StorageSettings) {
		settings.TTLState = ttlForStates
	}
}
func WithTTLForData(ttlForData time.Duration) OptionFunc {
	return func(settings *StorageSettings) {
		settings.TTLData = ttlForData
	}
}
func WithResetDataBatchSize(batchSize int64) OptionFunc {
	return func(settings *StorageSettings) {
		settings.ResetDataBatchSize = batchSize
	}
}

func FromOptions(opts ...OptionFunc) StorageSettings {
	s := StorageSettings{}
	sPtr := &s
	for _, opt := range opts {
		opt(sPtr)
	}
	return s
}
