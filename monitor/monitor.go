package monitor

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/naeimc/health/api"
	"github.com/naeimc/health/api/info"
	"github.com/naeimc/health/internal/network"
)

type Monitor struct {
	Log      func(any)
	emitters map[string]func()
}

func NewMonitor() *Monitor {
	return &Monitor{func(a any) {}, make(map[string]func())}
}

func (m *Monitor) Handle(packet network.Packet) ([]byte, error) {
	form := make(map[string]any)

	if err := json.Unmarshal(packet.Data(), &form); err != nil {
		return []byte("cannot unmarshal json\n"), nil
	}

	if err := api.IsValid(form); err != nil {
		return []byte(err.Error() + "\n"), nil
	}

	switch form[api.API_Key_Command] {
	case api.API_Command_StartEmitter:
		return m.handleStartEmitter(packet, form)
	case api.API_Command_StopEmitter:
		return m.handleStopEmitter(packet, form)
	}

	return []byte("not implemented\n"), nil
}

func (m *Monitor) handleStartEmitter(packet network.Packet, form map[string]any) ([]byte, error) {
	address := form[api.API_Key_Address].(string)
	duration, _ := time.ParseDuration(form[api.API_Key_Duration].(string))
	ticker := time.NewTicker(duration)
	connection, _ := net.Dial("udp", address)

	m.emitters[address] = ticker.Stop
	go m.startEmitter(connection, ticker, form)

	m.Log(fmt.Sprintf("(%s) START_EMITTER %d OK", packet.RemoteAddress(), len(packet.Data())))
	return []byte("OK"), nil
}

func (m *Monitor) handleStopEmitter(packet network.Packet, form map[string]any) ([]byte, error) {
	address := form[api.API_Key_Address].(string)
	m.stopEmitter(address)

	m.Log(fmt.Sprintf("(%s) STOP_EMITTER %d OK", packet.RemoteAddress(), len(packet.Data())))
	return []byte("OK"), nil
}

func (m *Monitor) startEmitter(connection net.Conn, ticker *time.Ticker, form map[string]any) {
	for range ticker.C {
		data, err := json.MarshalIndent(info.Complete(form), "", "    ")
		if err != nil {
			return
		}
		size, _ := connection.Write(data)
		m.Log(fmt.Sprintf("(%s) EMIT %d", connection.RemoteAddr(), size))
	}
}

func (m *Monitor) stopEmitter(address string) {
	stop, ok := m.emitters[address]
	if !ok {
		return
	}
	stop()
	delete(m.emitters, address)
}
