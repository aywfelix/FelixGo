package eventsystem

type EventType int

const (
	ET_Unknown       EventType = 0
	ET_Login         EventType = 1 // 玩家登录进入游戏
	ET_Logout        EventType = 2 // 玩家退出游戏
	ET_EnemyKilled   EventType = 3 //
	ET_SoldierKilled EventType = 4 //
)

type IEvent interface {
	Release()
	RegisterObserver(eventType EventType, observer IEventObserver)
	NotifySubject(eventType EventType, params ...interface{})
}

type Event struct {
	events map[EventType]IEventSubject
}

func NewEvent() *Event {
	return &Event{
		events: make(map[EventType]IEventSubject),
	}
}

func (e *Event) RegisterObserver(eventType EventType, observer IEventObserver) {
	subject := e.getEventSubject(eventType)
	if subject == nil {
		return
	}
	subject.Attach(observer)
	observer.SetSubject(subject)
}

func (e *Event) NotifySubject(eventType EventType, params ...interface{}) {
	if subject, ok := e.events[eventType]; ok {
		subject.SetParams(params...)
	}
}

func (e *Event) Release() {
	e.events = nil
}

func (e *Event) getEventSubject(eventType EventType) IEventSubject {
	if iSubject, ok := e.events[eventType]; ok {
		return iSubject
	}
	var subject IEventSubject
	switch eventType {
	case ET_Login:
		subject = NewLoginSubject()
	case ET_Logout:
		subject = NewLogoutSubject()
	case ET_EnemyKilled:
		subject = NewEnemyKilledSubject()
	case ET_SoldierKilled:
		subject = NewSoldierKilledSubject()
	default:
	}
	if subject != nil {
		e.events[eventType] = subject
	}
	return subject
}
