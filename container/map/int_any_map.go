package container

import (
	"encoding/json"
	"sync"
)

type IntAnyMap struct {
	mutext sync.RWMutex
	data   map[int]interface{}
}

func NewIntAnyMap() *IntAnyMap {
	return &IntAnyMap{
		mutext: sync.RWMutex{},
		data:   make(map[int]interface{}),
	}
}

func NewIntAnyMapFrom(data map[int]interface{}) *IntAnyMap {
	return &IntAnyMap{
		mutext: sync.RWMutex{},
		data:   data,
	}
}

func (m *IntAnyMap) Clone() *IntAnyMap {
	return NewIntAnyMapFrom(m.Copy())
}

func (m *IntAnyMap) Iterator(f func(key int, v interface{}) bool) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *IntAnyMap) Copy() map[int]interface{} {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	data := make(map[int]interface{})
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntAnyMap) Set(key int, val interface{}) {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	m.data[key] = val
}

func (m *IntAnyMap) Sets(data map[int]interface{}) {
	m.mutext.Lock()
	defer m.mutext.Unlock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
}

func (m *IntAnyMap) Search(key int) (interface{}, bool) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()
	if m.data == nil {
		return nil, false
	}
	val, found := m.data[key]
	return val, found
}

func (m *IntAnyMap) Get(key int) (val interface{}) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	if m.data != nil {
		val, _ = m.data[key]
	}
	return nil
}

func (m *IntAnyMap) Pop() (key int, val interface{}) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	for key, val = range m.data {
		delete(m.data, key)
	}
	return
}

func (m *IntAnyMap) Removes(keys []int) {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
}

func (m *IntAnyMap) Remove(key int) interface{} {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	if m.data != nil {
		if val, ok := m.data[key]; ok {
			delete(m.data, key)
			return val
		}
	}
	return nil
}

func (m *IntAnyMap) Keys() []int {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	if m.data != nil {
		keys := make([]int, 0)
		for key, _ := range m.data {
			keys = append(keys, key)
		}
		return keys
	}
	return nil
}

func (m *IntAnyMap) Values() []interface{} {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	if m.data != nil {
		values := make([]interface{}, 0)
		for _, value := range m.data {
			values = append(values, value)
		}
		return values
	}
	return nil
}

func (m *IntAnyMap) Contains(key int) bool {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	var ok bool
	if m.data != nil {
		_, ok = m.data[key]
	}
	return ok
}

func (m *IntAnyMap) Size() int {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	return len(m.data)
}

func (m *IntAnyMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *IntAnyMap) Clear() {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	m.data = make(map[int]interface{})
}

func (m *IntAnyMap) Replace(data map[int]interface{}) {
	m.mutext.Lock()
	defer m.mutext.Unlock()
	m.data = data
}

func (m *IntAnyMap) RLockFunc(f func(map[int]interface{})) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	f(m.data)
}

func (m *IntAnyMap) LockFunc(f func(map[int]interface{})) {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	f(m.data)
}

func (m *IntAnyMap) Merge(other *IntAnyMap) {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	if m.data == nil {
		m.data = other.Copy()
	}

	if other != m {
		other.mutext.RLock()
		defer other.mutext.RUnlock()
		for k, v := range other.data {
			m.data[k] = v
		}
	}
}

func (m *IntAnyMap) Marshal() ([]byte, error) {
	m.mutext.RLock()
	defer m.mutext.RUnlock()

	return json.Marshal(m.data)
}

func (m *IntAnyMap) Unmarshal(bytes []byte) error {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	var data interface{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	m.data, _ = data.(map[int]interface{})
	return nil
}
