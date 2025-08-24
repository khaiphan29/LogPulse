package helpers

import (
   "context"
   "time"
   "errors"
   "strconv"
)

type MockCacheClient struct {
}

var mockValue map[string]string = make(map[string]string)

func (m *MockCacheClient) SetRecord(ctx context.Context, key string, value string, exp time.Duration) error {
   mockValue[key] = value
   return nil
}

func (m *MockCacheClient) GetRecord(ctx context.Context, key string) (any, error) {
   if val, exists := mockValue[key]; !exists {
      return "0", errors.New("key does not exist")
   } else {
      return val, nil
   }
}

func (m *MockCacheClient) DeleteRecord(ctx context.Context, key string) error {
   if _, exists := mockValue[key]; !exists {
      return errors.New("key does not exist")
   }
   delete(mockValue, key)
   return nil
}

func (m *MockCacheClient) ErrKeyNotExists(err error) bool {
   if err != nil && err.Error() == "key does not exist" {
      return true
   }
   return false
}

func (m *MockCacheClient) IncreaseValue(ctx context.Context, key string, value int64) (int64, error) {
   val, exists := mockValue[key]
   if !exists {
      return 0, errors.New("key does not exist")
   }
   intVal, err := strconv.Atoi(val)
   if err != nil {
      return 0, err
   }
   intVal += int(value)
   mockValue[key] = strconv.Itoa(intVal)
   return int64(intVal), nil
}

func (m *MockCacheClient) DecreaseValue(ctx context.Context, key string, value int64) (int64, error) {
   val, exists := mockValue[key]
   if !exists {
      return 0, errors.New("key does not exist")
   }
   intVal, err := strconv.Atoi(val)
   if err != nil {
      return 0, err
   }
   intVal -= int(value)
   mockValue[key] = strconv.Itoa(intVal)
   return int64(intVal), nil
}
