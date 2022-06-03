package testutils

import (
	"net"
	"time"
)

type FakePacketConn struct {
	WriteToCalled  bool
	WriteToContent []byte
}

func NewFakePacketConn() *FakePacketConn {
	return &FakePacketConn{
		WriteToCalled: false,
	}
}

func (fpc *FakePacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return 0, &net.IPAddr{}, nil
}
func (fpc *FakePacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	fpc.WriteToCalled = true
	fpc.WriteToContent = p

	return len(p), nil
}
func (fpc *FakePacketConn) Close() error {
	return nil
}
func (fpc *FakePacketConn) LocalAddr() net.Addr {
	return &net.IPAddr{}
}
func (fpc *FakePacketConn) SetDeadline(t time.Time) error {
	return nil
}
func (fpc *FakePacketConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (fpc *FakePacketConn) SetWriteDeadline(t time.Time) error {
	return nil
}
