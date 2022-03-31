package container

import (
	"encoding/json"
	"sync"
)

type anymap map[interface{}]interface{}

type AnyAnyMap struct {
	mutex sync.RWMutex
	data  anymap
}

func NewAnyAnyMap() *AnyAnyMap {
	return &AnyAnyMap{
		mutex: sync.RWMutex{},
		data:  make(anymap),
	}
}

func NewAnyAnyMapFrom(data anymap) *AnyAnyMap {
	return &AnyAnyMap{
		mutex: sync.RWMutex{},
		data:  data,
	}
}

func (m *AnyAnyMap) Iterator(f func(k interface{}, v interface{}) bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *AnyAnyMap) Clone() *AnyAnyMap {
	return NewAnyAnyMapFrom(m.MapCopy())
}

func (m *AnyAnyMap) MapCopy() anymap {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data := make(anymap, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *AnyAnyMap) FilterEmpty() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		if v == nil {
			delete(m.data, k)
		}
	}
}

func (m *AnyAnyMap) FilterNil() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		if v == nil {
			delete(m.data, k)
		}
	}
}

func (m *AnyAnyMap) Set(key interface{}, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = make(map[interface{}]interface{})
	}
	m.data[key] = value
}

func (m *AnyAnyMap) Sets(data map[interface{}]interface{}) {
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

func (m *AnyAnyMap) Search(key interface{}) (value interface{}, found bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data != nil {
		value, found = m.data[key]
	}
	return
}

func (m *AnyAnyMap) Get(key interface{}) (value interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data != nil {
		value, _ = m.data[key]
	}
	return
}

func (m *AnyAnyMap) Pop() (key, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *AnyAnyMap) Pops(size int) map[interface{}]interface{} {
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
		newMap = make(map[interface{}]interface{}, size)
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

func (m *AnyAnyMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = make(map[interface{}]interface{})
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

func (m *AnyAnyMap) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *AnyAnyMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *AnyAnyMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (m *AnyAnyMap) SetIfNotExist(key interface{}, value interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *AnyAnyMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *AnyAnyMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (m *AnyAnyMap) Remove(key interface{}) (value interface{}) {
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

func (m *AnyAnyMap) Removes(keys []interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
}

func (m *AnyAnyMap) Keys() []interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var (
		keys  = make([]interface{}, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	return keys
}

func (m *AnyAnyMap) Values() []interface{} {
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

func (m *AnyAnyMap) Contains(key interface{}) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var ok bool
	if m.data != nil {
		_, ok = m.data[key]
	}
	return ok
}

func (m *AnyAnyMap) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	length := len(m.data)
	return length
}

func (m *AnyAnyMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *AnyAnyMap) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[interface{}]interface{})
}

func (m *AnyAnyMap) Replace(data map[interface{}]interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = data
}

func (m *AnyAnyMap) LockFunc(f func(m map[interface{}]interface{})) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	f(m.data)
}

func (m *AnyAnyMap) RLockFunc(f func(m map[interface{}]interface{})) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	f(m.data)
}

func (m *AnyAnyMap) Merge(other *AnyAnyMap) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mutex.RUnlock()
		defer other.mutex.RUnlock()
		for k, v := range other.data {
			m.data[k] = v
		}
	}
}

func (m *AnyAnyMap) String() string {
	b, _ := m.Marshal()
	return string(b)
}

func (m *AnyAnyMap) Marshal() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m *AnyAnyMap) Unmarshal(bytes []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.data = make(map[interface{}]interface{})
	}
	var data interface{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	m.data, _ = data.(map[interface{}]interface{})
	return nil
}
