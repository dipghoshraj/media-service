package sni

// Package internal provides a ConnBuffer type that wraps a net.Conn
// and allows reading from a buffer before falling back to the connection.
// It implements the net.Conn interface, allowing it to be used in place of a net.Conn.
// It is useful for scenarios where you want to read data from a buffer first,

import (
	"net"
	"time"
)

type ConnBuffer struct {
	buf  []byte
	conn net.Conn
}

func (cb *ConnBuffer) Read(p []byte) (int, error) {
	if len(cb.buf) > 0 {
		n := copy(p, cb.buf)
		cb.buf = cb.buf[n:]
		return n, nil
	}
	return cb.conn.Read(p)
}

func (cb *ConnBuffer) Write(p []byte) (int, error) {
	return cb.conn.Write(p)
}

func (cb *ConnBuffer) Close() error {
	return cb.conn.Close()
}
func (cb *ConnBuffer) LocalAddr() net.Addr                { return cb.conn.LocalAddr() }
func (cb *ConnBuffer) RemoteAddr() net.Addr               { return cb.conn.RemoteAddr() }
func (cb *ConnBuffer) SetDeadline(t time.Time) error      { return cb.conn.SetDeadline(t) }
func (cb *ConnBuffer) SetReadDeadline(t time.Time) error  { return cb.conn.SetReadDeadline(t) }
func (cb *ConnBuffer) SetWriteDeadline(t time.Time) error { return cb.conn.SetWriteDeadline(t) }
