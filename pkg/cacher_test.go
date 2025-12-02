package pkg

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacher_Basic_Success(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		counter++
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 100*time.Millisecond)
	assert.NoError(t, err)

	// Первое чтение — обновление
	val, err := c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	// Второе чтение до истечения TTL — кэш
	val, err = c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	// Ждем истечения TTL
	time.Sleep(120 * time.Millisecond)
	val, err = c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
}

func TestCacher_Concurrent_Success(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		time.Sleep(50 * time.Millisecond) // имитируем долгую работу
		counter++
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 100*time.Millisecond)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	results := make([]int, 10)
	errs := make([]error, 10)

	// Принудительно делаем данные устаревшими
	c.lastSync = time.Now().Add(-200 * time.Millisecond)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			val, err := c.GetData(t.Context())
			results[i] = val
			errs[i] = err
		}(i)
	}
	wg.Wait()

	// Все должны получить одно и то же обновленное значение
	for i, val := range results {
		assert.Equal(t, 2, val, "goroutine %d: expected 2", i)
		assert.NoError(t, errs[i], "goroutine %d: unexpected error", i)
	}

	// Проверяем, что updateFn вызвался только один раз для обновления
	assert.Equal(t, 2, counter)
}

func TestCacher_InvalidUpdate_Error(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		counter++
		if counter == 2 {
			return 0, errors.New("update failed")
		}
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 50*time.Millisecond)
	assert.NoError(t, err)

	// Первая успешная загрузка
	val, err := c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	time.Sleep(60 * time.Millisecond)

	// Вторая попытка — ошибка, должны вернуть старые данные
	val, err = c.GetData(t.Context())
	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.Equal(t, 1, val)
}

func TestCacher_ForceSync_Success(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		counter++
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 1*time.Hour)
	assert.NoError(t, err)

	// Изначально данные не загружены
	val, err := c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	// Принудительно обновляем данные
	val, err = c.ForceSync(t.Context())
	assert.NoError(t, err)

	// После ForceSync значение должно увеличиться
	assert.Equal(t, 2, val)

	val, err = c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
}

func TestCacher_ForceSync_Error(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		counter++
		if counter == 2 {
			return 0, errors.New("force sync failed")
		}
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 1*time.Hour)
	assert.NoError(t, err)

	// Первоначальная загрузка
	val, err := c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	// Вторая попытка ForceSync вернет ошибку
	val, err = c.ForceSync(t.Context())
	assert.Error(t, err)
	assert.EqualError(t, err, "force sync failed")

	// Значение при ошибке не меняется
	val, err = c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestCacher_ForceSync_Concurrent(t *testing.T) {
	counter := 0
	updateFn := func(ctx context.Context) (int, error) {
		time.Sleep(50 * time.Millisecond)
		counter++
		return counter, nil
	}

	c, err := NewCacher(t.Context(), updateFn, 1*time.Hour)
	assert.NoError(t, err)

	// Первоначальная загрузка
	val, err := c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	var wg sync.WaitGroup
	errs := make([]error, 5)

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := c.ForceSync(t.Context())
			errs[i] = err
		}(i)
	}
	wg.Wait()

	// Все вызовы ForceSync должны завершиться без ошибок
	for i, err := range errs {
		assert.NoError(t, err, "goroutine %d: unexpected error", i)
	}

	// Значение должно увеличиться ровно на 1 после всех ForceSync
	val, err = c.GetData(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
}
