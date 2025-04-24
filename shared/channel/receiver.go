package channel

import (
	"fmt"
	"net"
	"sync"
)

type fragmentsData struct {
	totalFragments uint16
	received       map[uint16][]byte
	receivedCount  int
}

type Receiver struct {
	conn *net.UDPConn

	mu           sync.Mutex
	fragmentsMap map[uint32]*fragmentsData
	successSubs  []func(MessageType, []byte)
	errorSubs    []func(error)
	stopChan     chan struct{}
	stoppedChan  chan struct{}
}

func NewReceiver(conn *net.UDPConn) *Receiver {
	return &Receiver{
		conn:         conn,
		fragmentsMap: make(map[uint32]*fragmentsData),
		stopChan:     make(chan struct{}),
		stoppedChan:  make(chan struct{}),
	}
}

func (r *Receiver) OnSuccess(sub func(MessageType, []byte)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.successSubs = append(r.successSubs, sub)
}

func (r *Receiver) OnError(sub func(error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errorSubs = append(r.errorSubs, sub)
}

func (r *Receiver) triggerSuccess(messageType MessageType, payload []byte) {
	r.mu.Lock()
	subscribers := append([]func(MessageType, []byte){}, r.successSubs...)
	r.mu.Unlock()

	for _, sub := range subscribers {
		go sub(messageType, payload)
	}
}

func (r *Receiver) triggerError(err error) {
	r.mu.Lock()
	subscribers := append([]func(error){}, r.errorSubs...)
	r.mu.Unlock()

	for _, sub := range subscribers {
		go sub(err)
	}
}

func (r *Receiver) Run() {
	buffer := make([]byte, 65535)

	for {
		select {
		case <-r.stopChan:
			close(r.stoppedChan)
			return
		default:
		}

		n, _, err := r.conn.ReadFromUDP(buffer)
		if err != nil {
			r.triggerError(fmt.Errorf("ReadFromUDP error: %w", err))
			continue
		}

		pkt, err := Deserialize(buffer[:n])
		if err != nil {
			r.triggerError(fmt.Errorf("Deserialize error: %w", err))
			continue
		}

		r.mu.Lock()
		fd, exists := r.fragmentsMap[pkt.MessageID]
		if !exists {
			fd = &fragmentsData{
				totalFragments: pkt.TotalFragments,
				received:       make(map[uint16][]byte),
				receivedCount:  0,
			}
			r.fragmentsMap[pkt.MessageID] = fd
		}

		if _, ok := fd.received[pkt.FragmentIndex]; !ok {
			fd.received[pkt.FragmentIndex] = pkt.Payload
			fd.receivedCount++
		}

		if fd.receivedCount == int(fd.totalFragments) {
			var fullPayload []byte
			for i := uint16(0); i < fd.totalFragments; i++ {
				part, ok := fd.received[i]
				if !ok {
					r.mu.Unlock()
					r.triggerError(fmt.Errorf("missing fragment %d for messageID %d, discarding", i, pkt.MessageID))
					r.mu.Lock()
					delete(r.fragmentsMap, pkt.MessageID)
					r.mu.Unlock()
					goto CONTINUE_LOOP
				}
				fullPayload = append(fullPayload, part...)
			}

			if len(fullPayload) == 0 {
				r.mu.Unlock()
				r.triggerError(fmt.Errorf("reassembled payload empty for messageID %d", pkt.MessageID))
				r.mu.Lock()
				delete(r.fragmentsMap, pkt.MessageID)
				r.mu.Unlock()
				goto CONTINUE_LOOP
			}

			messageType := MessageType(fullPayload[0])
			payloadWithoutType := fullPayload[1:]

			delete(r.fragmentsMap, pkt.MessageID)
			r.mu.Unlock()

			r.triggerSuccess(messageType, payloadWithoutType)
		} else {
			r.mu.Unlock()
		}

	CONTINUE_LOOP:
	}
}
