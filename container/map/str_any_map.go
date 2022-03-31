package container

import (
	"encoding/json"
	"sync"
)

type StrAnyMap struct {
	mutex sync.RWMutex
	data  map[string]interface{}
}

func NewStrAnyMap() *StrAnyMap {
	return &StrAnyMap{
		mutex: sync.RWMutex{},
		data:  make(map[string]interface{}),
	}
}

func NewStrAnyMapFrom(data map[string]interface{}) *StrAnyMap {
	return &StrAnyMap{
		mutex: sync.RWMutex{},
		data:  data,
	}
}

func (m *StrAnyMap) Iterator(f func(k string, v interface{}) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *StrAnyMap) Clone() *StrAnyMap {
	return NewStrAnyMapFrom(m.MapCopy())
}

func (m *StrAnyMap) MapCopy() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *StrAnyMap) FilterEmpty() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		if v == nil {
			delete(m.data, k)
		}
	}
}

func (m *StrAnyMap) FilterNil() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		if v == nil {
			delete(m.data, k)
		}
	}
}

func (m *StrAnyMap) Set(key string, val interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = val
}

func (m *StrAnyMap) Sets(data map[string]interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
}

func (m *StrAnyMap) Search(key string) (value interface{}, found bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data != nil {
		value, found = m.data[key]
	}
	return
}

func (m *StrAnyMap) Get(key string) (value interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data != nil {
		value, _ = m.data[key]
	}
	return
}

func (m *StrAnyMap) Pop() (key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *StrAnyMap) Pops(size int) map[string]interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(map[string]interface{}, size)
	)
	for k, v := range m.data {
		delete(m.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

func (m *StrAnyMap) doSetWithLockCheck(key string, value interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = value
	}
	return value
}

func (m *StrAnyMap) GetOrSet(key string, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *StrAnyMap) GetOrSetFunc(key string, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *StrAnyMap) GetOrSetFuncLock(key string, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (m *StrAnyMap) SetIfNotExist(key string, value interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *StrAnyMap) SetIfNotExistFunc(key string, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *StrAnyMap) SetIfNotExistFuncLock(key string, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (m *StrAnyMap) Removes(keys []string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
}

func (m *StrAnyMap) Remove(key string) (value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data != nil {
		var ok bool
		if value, ok = m.data[key]; ok {
			delete(m.data, key)
		}
	}
	return
}

func (m *StrAnyMap) Keys() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var (
		keys  = make([]string, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	return keys
}

func (m *StrAnyMap) Values() []interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var (
		values = make([]interface{}, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	return values
}

func (m *StrAnyMap) Contains(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var ok bool
	if m.data != nil {
		_, ok = m.data[key]
	}
	return ok
}

func (m *StrAnyMap) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	length := len(m.data)
	return length
}

func (m *StrAnyMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *StrAnyMap) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[string]interface{})
}

func (m *StrAnyMap) Replace(data map[string]interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = data
}

func (m *StrAnyMap) LockFunc(f func(m map[string]interface{})) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	f(m.data)
}

func (m *StrAnyMap) RLockFunc(f func(m map[string]interface{})) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	f(m.data)
}

func (m *StrAnyMap) Merge(other *StrAnyMap) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mutex.RLock()
		defer other.mutex.RUnlock()
		for k, v := range other.data {
			m.data[k] = v
		}
	}

}

func (m *StrAnyMap) String() string {
	b, _ := m.Marshal()
	return string(b)
}

func (m *StrAnyMap) Marshal() ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return json.Marshal(m.data)
}

func (m *StrAnyMap) Unmarshal(b []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var data interface{}
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data, _ = data.(map[string]interface{})
	return nil
}
