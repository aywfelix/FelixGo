package container

type SetAny map[interface{}]struct{}
type SetInt map[int]struct{}
type SetStr map[string]struct{}
type SetUInt64 map[uint64]struct{}

type ISet interface {
	Add(key interface{})
	Remove(key interface{})
	Contains(key interface{}) bool
	ForEach(func(key interface{}) bool)
	ToArray() []interface{}
	Len() int
	Clear()
	IsEmpty() bool
}

func (s SetAny) Add(key interface{}) {
	if s == nil {
		s = make(map[interface{}]struct{})
	}

	if _, ok := s[key]; !ok {
		s[key] = struct{}{}
	}
}

func (s SetAny) Remove(key interface{}) {
	if s == nil {
		return
	}
	if _, ok := s[key]; !ok {
		return
	}
	delete(s, key)
}

func (s SetAny) Contains(key interface{}) bool {
	if s == nil {
		return false
	}
	if _, ok := s[key]; !ok {
		return false
	}
	return true
}

func (s SetAny) ForEach(f func(key interface{}) bool) {
	if s == nil || len(s) == 0 {
		return
	}

	for key, _ := range s {
		if !f(key) {
			break
		}
	}
}

func (s SetAny) ToArray() []interface{} {
	if s == nil || len(s) == 0 {
		return nil
	}
	list := make([]interface{}, 0)
	for key, _ := range s {
		list = append(list, key)
	}
	return list
}

func (s SetAny) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

func (s SetAny) Clear() {
	if s != nil {
		s = nil
	}
}

func (s SetAny) IsEmpty() bool {
	return len(s) == 0
}

//==============================================================================
func (s SetInt) Add(key int) {
	if s == nil {
		s = make(map[int]struct{})
	}

	if _, ok := s[key]; !ok {
		s[key] = struct{}{}
	}
}

func (s SetInt) Remove(key int) {
	if s == nil {
		return
	}
	if _, ok := s[key]; !ok {
		return
	}
	delete(s, key)
}

func (s SetInt) Contains(key int) bool {
	if s == nil {
		return false
	}
	if _, ok := s[key]; !ok {
		return false
	}
	return true
}

func (s SetInt) ToArray() []int {
	if s == nil || len(s) == 0 {
		return nil
	}
	list := make([]int, 0)
	for key, _ := range s {
		list = append(list, key)
	}
	return list
}

func (s SetInt) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

func (s SetInt) Clear() {
	if s != nil {
		s = nil
	}
}

//==============================================================================
func (s SetStr) Add(key string) {
	if s == nil {
		s = make(map[string]struct{})
	}

	if _, ok := s[key]; !ok {
		s[key] = struct{}{}
	}
}

func (s SetStr) Remove(key string) {
	if s == nil {
		return
	}
	if _, ok := s[key]; !ok {
		return
	}
	delete(s, key)
}

func (s SetStr) Contains(key string) bool {
	if s == nil {
		return false
	}
	if _, ok := s[key]; !ok {
		return false
	}
	return true
}

func (s SetStr) ToArr() []string {
	if s == nil || len(s) == 0 {
		return nil
	}
	list := make([]string, 0)
	for key, _ := range s {
		list = append(list, key)
	}
	return list
}

func (s SetStr) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

func (s SetStr) Clear() {
	if s != nil {
		s = nil
	}
}

//=====================================================================================
func (s SetUInt64) Add(key uint64) {
	if s == nil {
		s = make(map[uint64]struct{})
	}

	if _, ok := s[key]; !ok {
		s[key] = struct{}{}
	}
}

func (s SetUInt64) Remove(key uint64) {
	if s == nil {
		return
	}
	if _, ok := s[key]; !ok {
		return
	}
	delete(s, key)
}

func (s SetUInt64) Contains(key uint64) bool {
	if s == nil {
		return false
	}
	if _, ok := s[key]; !ok {
		return false
	}
	return true
}

func (s SetUInt64) ToArray() []uint64 {
	if s == nil || len(s) == 0 {
		return nil
	}
	list := make([]uint64, 0)
	for key, _ := range s {
		list = append(list, key)
	}
	return list
}

func (s SetUInt64) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

func (s SetUInt64) Clear() {
	if s != nil {
		s = nil
	}
}
