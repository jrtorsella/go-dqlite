package protocol

import "net"

func (m *Message) Body() ([]byte, int) {
	return m.body.Bytes, m.body.Offset
}

func (m *Message) Rewind() {
	m.body.Offset = 0
}

func NewProtocol(version uint64, conn net.Conn) *Protocol {
	return &Protocol{version: version, conn: conn}
}
