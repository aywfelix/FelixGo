package fnet

import (
	"errors"
	"sync/atomic"
)

var (
	errSessionNotFound = errors.New("session not found")
)

type ISessionManager interface {
	Add(session ISession)
	Remove(session ISession)
	RemoveById(sessionID int64)
	Len() int
	Clear()
	GetSession(sessionID int64) (ISession, error)
	ClearSession(sID int64)
}

type SessionManager struct {
	atomicValue atomic.Value
}

func NewSessionManager() *SessionManager {
	cm := &SessionManager{}
	sessionMap := make(map[int64]ISession)
	cm.atomicValue.Store(sessionMap)
	return cm
}

func (cm *SessionManager) Add(session ISession) {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	sessionMap[session.GetID()] = session
	cm.atomicValue.Store(sessionMap)
}

func (cm *SessionManager) Remove(session ISession) {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	sID := session.GetID()
	if _, ok := sessionMap[sID]; ok {
		delete(sessionMap, sID)
		cm.atomicValue.Store(sessionMap)
	}
}

func (cm *SessionManager) RemoveById(sessionID int64) {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	if _, ok := sessionMap[sessionID]; ok {
		delete(sessionMap, sessionID)
		cm.atomicValue.Store(sessionMap)
	}
}

func (cm *SessionManager) Len() int {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	return len(sessionMap)
}

func (cm *SessionManager) Clear() {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	for sessionID, session := range sessionMap {
		session.Stop()
		delete(sessionMap, sessionID)
	}
	cm.atomicValue.Store(sessionMap)
}

func (cm *SessionManager) GetSession(sessionID int64) (ISession, error) {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	if session, ok := sessionMap[sessionID]; ok {
		return session, nil
	}
	return nil, errSessionNotFound
}

func (cm *SessionManager) ClearSession(sID int64) {
	sessionMap := cm.atomicValue.Load().(map[int64]ISession)
	if session, ok := sessionMap[sID]; ok {
		session.Stop()
		delete(sessionMap, sID)
		cm.atomicValue.Store(sessionMap)
	}
}
