package fnet

type ConnDataMap map[int32]*ConnData

type ConnData struct {
	NodeID     int32
	Port       int32
	IP         string
	ServerName string
	ServerType int32
	ConnState  ConnState

	netService INetService
}

func (m ConnDataMap) Add(connData *ConnData) {
	if m == nil {
		m = make(map[int32]*ConnData)
	}
	if _, ok := m[connData.NodeID]; !ok {
		m[connData.NodeID] = connData
	}
}

func (m ConnDataMap) Remove(serverID int32) {
	if _, ok := m[serverID]; !ok {
		return
	}
	delete(m, serverID)
}

func (m ConnDataMap) GetByServerID(serverID int32) *ConnData {
	if _, ok := m[serverID]; !ok {
		return nil
	}
	return m[serverID]
}

func (m ConnDataMap) GetByServerType(serverType int32) *ConnData {
	for _, connData := range m {
		if connData.ServerType == serverType {
			return connData
		}
	}
	return nil
}
