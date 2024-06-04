// Package redisfsm contains redis storage.
package redisfsm

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

type keyType string

const (
	stateKey keyType = "state"
	dataKey  keyType = "data"
)

type Storage struct {
	rds  redis.UniversalClient
	pref StorageSettings
}

var _ fsm.Storage = (*Storage)(nil)

type StorageSettings struct {
	// Prefix for records in Redis.
	// Default is "fsm".
	Prefix string

	// TTL for state.
	// Default is 0 (no ttl).
	TTLState time.Duration

	// TTL for state data.
	// Default is 0 (no ttl).
	TTLData time.Duration

	// Batch size for reset data.
	// Default is 0 (no batching).
	ResetDataBatchSize int64
}

const defaultPrefix = "fsm"

// NewStorage returns new redis storage.
func NewStorage(client redis.UniversalClient, opts ...OptionFunc) *Storage {
	pref := StorageSettings{Prefix: defaultPrefix}
	if len(opts) != 0 {
		prefPtr := &pref
		for _, opt := range opts {
			opt(prefPtr)
		}
	}

	return &Storage{
		rds:  client,
		pref: pref,
	}
}

func (s *Storage) State(ctx context.Context, key fsm.StorageKey) (fsm.State, error) {
	val, err := s.rds.Get(ctx, s.generateKey(key, stateKey)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fsm.DefaultState, nil
		}
		return fsm.DefaultState, wrapError(err, "get state")
	}

	return fsm.State(val), nil
}

func (s *Storage) SetState(ctx context.Context, key fsm.StorageKey, state fsm.State) error {
	err := s.rds.Set(
		ctx,
		s.generateKey(key, stateKey),
		string(state),
		s.pref.TTLState,
	).Err()
	return wrapError(err, "set state")
}

func (s *Storage) ResetState(ctx context.Context, key fsm.StorageKey, withData bool) error {
	err := s.rds.Del(ctx, s.generateKey(key, stateKey)).Err()
	if err != nil {
		return wrapError(err, "reset state")
	}

	if withData {
		if err := s.resetData(ctx, key); err != nil {
			return wrapError(err, "reset data")
		}
	}
	return nil
}

func (s *Storage) resetData(ctx context.Context, key fsm.StorageKey) error {
	var cursor uint64
	var keys []string

	redisKey := s.generateKey(key, dataKey, "*")

	for {
		var err error
		keys, cursor, err = s.rds.
			Scan(ctx, cursor, redisKey, s.pref.ResetDataBatchSize).
			Result()
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		if len(keys) > 0 {
			if err := s.rds.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("delete keys: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *Storage) UpdateData(ctx context.Context, targetKey fsm.StorageKey, key string, data interface{}) error {
	redisKey := s.generateKey(targetKey, dataKey, key)

	if data == nil {
		err := s.rds.Del(ctx, redisKey).Err()
		return wrapError(err, "delete data")
	}

	encodedData, err := s.encode(data)
	if err != nil {
		return wrapError(err, "encode data")
	}

	err = s.rds.
		Set(ctx, redisKey, encodedData, s.pref.TTLData).
		Err()

	return wrapError(err, "set data")
}

func (s *Storage) Data(ctx context.Context, targetKey fsm.StorageKey, key string, to interface{}) error {
	dataBytes, err := s.rds.
		Get(ctx, s.generateKey(targetKey, dataKey, key)).
		Bytes()

	if errors.Is(err, redis.Nil) {
		return fsm.ErrNotFound
	}
	if err != nil {
		return wrapError(err, "get data")
	}

	if err := s.decode(dataBytes, to); err != nil {
		return wrapError(err, "decode data")
	}
	return nil
}

func (s *Storage) Close() error {
	return s.rds.Close()
}

func (s *Storage) generateKey(key fsm.StorageKey, keyType keyType, keys ...string) string {
	const (
		prefixPart     = 1
		baseTargetPart = 3 // bot + chat + user
		keyTypePart    = 1
	)
	keyPartsCount := prefixPart + baseTargetPart + keyTypePart + len(keys)
	if key.ThreadID != 0 {
		keyPartsCount++
	}
	parts := make([]string, 0, keyPartsCount) // prefix + key parts + keyType + keys
	parts = append(
		parts,
		s.pref.Prefix,
		strconv.FormatInt(key.BotID, 10),
		strconv.FormatInt(key.ChatID, 10),
		strconv.FormatInt(key.UserID, 10),
	)

	if key.ThreadID != 0 {
		parts = append(parts, strconv.FormatInt(key.ThreadID, 10))
	}

	parts = append(parts, string(keyType))

	if len(keys) > 0 {
		parts = append(parts, keys...)

	}

	return strings.Join(parts, ":")
}

func (s *Storage) encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Storage) decode(data []byte, to interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(to)
}

type ErrOperation struct {
	Operation string
	Err       error
}

func (e ErrOperation) Unwrap() error { return e.Err }
func (e ErrOperation) Error() string {
	return fmt.Sprintf("fsm-telebot/storage/redis: %s: %v", e.Operation, e.Err)
}

func wrapError(err error, op string) error {
	if err == nil {
		return nil
	}
	return &ErrOperation{Operation: op, Err: err}
}
