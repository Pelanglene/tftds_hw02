package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_ExecuteOperation(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	t.Run("Операции создания", func(t *testing.T) {
		testCases := []struct {
			name       string
			op         Operation
			wantResult bool
			wantErr    error
		}{
			{
				name: "успешное создание",
				op: Operation{
					Command:   "CREATE",
					Key:       "test1",
					Value:     strPtr("value1"),
					Timestamp: time.Now(),
				},
				wantResult: true,
				wantErr:    nil,
			},
			{
				name: "дублирование ключа",
				op: Operation{
					Command:   "CREATE",
					Key:       "test1",
					Value:     strPtr("value2"),
					Timestamp: time.Now(),
				},
				wantResult: false,
				wantErr:    ErrKeyExists,
			},
			{
				name: "nil значение",
				op: Operation{
					Command:   "CREATE",
					Key:       "test2",
					Value:     nil,
					Timestamp: time.Now(),
				},
				wantResult: false,
				wantErr:    ErrNilValue,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := storage.ExecuteOperation(ctx, tc.op)
				assert.Equal(t, tc.wantResult, result)
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Операции обновления и удаления", func(t *testing.T) {
		// Подготовка данных
		key := "test_key"
		initValue := "initial"
		newValue := "updated"

		createOp := Operation{
			Command:   "CREATE",
			Key:       key,
			Value:     &initValue,
			Timestamp: time.Now(),
		}
		_, err := storage.ExecuteOperation(ctx, createOp)
		require.NoError(t, err)

		// Тестирование UPDATE
		updateOp := Operation{
			Command:   "UPDATE",
			Key:       key,
			Value:     &newValue,
			Timestamp: time.Now(),
		}
		result, err := storage.ExecuteOperation(ctx, updateOp)
		assert.True(t, result)
		assert.NoError(t, err)

		// Проверка значения
		val, exists := storage.Get(key)
		assert.True(t, exists)
		assert.Equal(t, newValue, val)

		// Тестирование DELETE
		deleteOp := Operation{
			Command:   "DELETE",
			Key:       key,
			Timestamp: time.Now(),
		}
		result, err = storage.ExecuteOperation(ctx, deleteOp)
		assert.True(t, result)
		assert.NoError(t, err)

		// Проверка удаления
		_, exists = storage.Get(key)
		assert.False(t, exists)
	})

	t.Run("Операция CAS", func(t *testing.T) {
		key := "cas_test"
		oldValue := "old"
		newValue := "new"
		wrongValue := "wrong"

		// Создание начального значения
		createOp := Operation{
			Command:   "CREATE",
			Key:       key,
			Value:     &oldValue,
			Timestamp: time.Now(),
		}
		_, err := storage.ExecuteOperation(ctx, createOp)
		require.NoError(t, err)

		// Успешный CAS
		casOp := Operation{
			Command:   "CAS",
			Key:       key,
			Value:     &newValue,
			OldValue:  &oldValue,
			Timestamp: time.Now(),
		}
		result, err := storage.ExecuteOperation(ctx, casOp)
		assert.True(t, result)
		assert.NoError(t, err)

		// Неуспешный CAS
		casOp.Value = &wrongValue
		casOp.OldValue = &oldValue
		result, err = storage.ExecuteOperation(ctx, casOp)
		assert.False(t, result)
		assert.NoError(t, err)
	})
}

func strPtr(s string) *string {
	return &s
}
