package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrKeyExists       = errors.New("key already exists")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrNilValue        = errors.New("value cannot be nil")
	ErrOperationFailed = errors.New("operation failed")
)

type Storage struct {
	data     sync.Map
	opsMutex sync.Mutex
	logger   *slog.Logger
}

type Operation struct {
	Command   string
	Key       string
	Value     *string
	OldValue  *string
	Timestamp time.Time
}

var defaultStorage = NewStorage()

func NewStorage() *Storage {
	return &Storage{
		logger: slog.Default().With("component", "storage"),
	}
}

func (s *Storage) ExecuteOperation(ctx context.Context, op Operation) (bool, error) {
	s.opsMutex.Lock()
	defer s.opsMutex.Unlock()

	s.logger.Info("executing operation",
		"command", op.Command,
		"key", op.Key,
		"time", op.Timestamp,
	)

	switch op.Command {
	case "CREATE":
		return s.handleCreate(op)
	case "DELETE":
		return s.handleDelete(op)
	case "UPDATE":
		return s.handleUpdate(op)
	case "CAS":
		return s.handleCAS(op)
	default:
		return false, fmt.Errorf("%w: %s", ErrInvalidCommand, op.Command)
	}
}

func (s *Storage) handleCreate(op Operation) (bool, error) {
	if op.Value == nil {
		return false, ErrNilValue
	}

	if _, loaded := s.data.LoadOrStore(op.Key, *op.Value); loaded {
		return false, ErrKeyExists
	}

	s.logger.Debug("created new record", "key", op.Key, "value", *op.Value)
	return true, nil
}

func (s *Storage) handleDelete(op Operation) (bool, error) {
	if _, ok := s.data.LoadAndDelete(op.Key); !ok {
		return false, nil
	}
	s.logger.Debug("deleted record", "key", op.Key)
	return true, nil
}

func (s *Storage) handleUpdate(op Operation) (bool, error) {
	if op.Value == nil {
		return false, ErrNilValue
	}

	if _, ok := s.data.Load(op.Key); !ok {
		return false, nil
	}

	s.data.Store(op.Key, *op.Value)
	s.logger.Debug("updated record", "key", op.Key, "new_value", *op.Value)
	return true, nil
}

func (s *Storage) handleCAS(op Operation) (bool, error) {
	if op.Value == nil || op.OldValue == nil {
		return false, ErrNilValue
	}

	success := s.data.CompareAndSwap(op.Key, *op.OldValue, *op.Value)
	if success {
		s.logger.Debug("выполнен CAS",
			"key", op.Key,
			"old_value", *op.OldValue,
			"new_value", *op.Value,
		)
	}
	return success, nil
}

func (s *Storage) Get(key string) (interface{}, bool) {
	return s.data.Load(key)
}

// Обертки для работы с глобальным хранилищем
func Get(key interface{}) (interface{}, bool) {
	return defaultStorage.Get(key.(string))
}

func ProcessWrite(command, key string, value, oldValue *string) (bool, error) {
	op := Operation{
		Command:   command,
		Key:       key,
		Value:     value,
		OldValue:  oldValue,
		Timestamp: time.Now(),
	}
	return defaultStorage.ExecuteOperation(context.Background(), op)
}
