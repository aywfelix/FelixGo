package eventsystem

type IEventSubject interface {
	Attach(observer IEventObserver)
	Detach(observer IEventObserver)
	Notify()
	SetParams(params ...interface{})
	GetParams() []interface{}
	GetCount() int
}

type EventSubject struct {
	observers []IEventObserver
	params    []interface{}

	count int
}

func NewEventSubject() IEventSubject {
	return &EventSubject{
		observers: make([]IEventObserver, 0),
	}
}

func (s *EventSubject) Attach(observer IEventObserver) {
	s.observers = append(s.observers, observer)
}

func (s *EventSubject) Detach(observer IEventObserver) {
	for k, v := range s.observers {
		if v == observer {
			s.observers = append(s.observers[:k], s.observers[k+1:]...)
			break
		}
	}
}
func (s *EventSubject) Notify() {
	for _, observer := range s.observers {
		observer.Update()
	}
}
func (s *EventSubject) SetParams(params ...interface{}) {
	s.params = params
	s.Notify()
}
func (s *EventSubject) GetParams() []interface{} {
	return s.params
}
func (s *EventSubject) GetCount() int {
	return s.count
}

type LoginSubject struct {
	EventSubject
}

func NewLoginSubject() IEventSubject {
	return &LoginSubject{}
}

type LogoutSubject struct {
	EventSubject
}

func NewLogoutSubject() IEventSubject {
	return &LogoutSubject{}
}

type EnemyKilledSubject struct {
	EventSubject
}

func NewEnemyKilledSubject() IEventSubject {
	return &EnemyKilledSubject{}
}

func (s *EnemyKilledSubject) SetParams(params ...interface{}) {
	s.params = params
	s.count++
	s.Notify()
}

type SouldierKilledSubject struct {
	EventSubject
}

func NewSoldierKilledSubject() IEventSubject {
	return &SouldierKilledSubject{}
}

func (s *SouldierKilledSubject) SetParams(params ...interface{}) {
	s.params = params
	s.count++
	s.Notify()
}
