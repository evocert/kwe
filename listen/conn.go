package listen

import (
	"bufio"
	"io"
	"net"
	"time"
)

type Conn struct {
	conn     net.Conn
	p        []byte
	initrdur time.Duration
	rdur     time.Duration
	rcnt     int32
	crntrdur time.Duration
	bufr     *bufio.Reader
	bufrw    *bufio.ReadWriter
	wcnt     int32
	initwdur time.Duration
	wdur     time.Duration
	crntwdur time.Duration
	bufw     *bufio.Writer
	n        int
	err      error
	//prdrr *io.PipeReader
	//prdrw *io.PipeWriter
	//pwtrr *io.PipeReader
	//pwtrw *io.PipeWriter
}

func NewCon(connn net.Conn) (conn *Conn) {
	conn = &Conn{conn: connn, initrdur: 30 * time.Second, rdur: 10 * time.Millisecond, crntrdur: 0, crntwdur: 0}
	conn.bufr = bufio.NewReaderSize(conn.conn, 1024*1024)
	conn.bufw = bufio.NewWriterSize(conn.conn, 1024*1024)
	conn.bufrw = bufio.NewReadWriter(conn.bufr, conn.bufw)
	return
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (conn *Conn) Read(b []byte) (n int, err error) {
	if bl := len(b); bl > 0 {
		conn.bufrw.Flush()
		conn.crntwdur = 0
		if conn.crntrdur == 0 {
			conn.crntrdur = conn.initrdur
		} else if conn.crntrdur != conn.rdur {
			conn.crntrdur = conn.rdur
		}
		for n < bl && err == nil {
			rn, rerr := conn.bufrw.Read(b[n : n+(bl-n)])
			if rn > 0 {
				conn.conn.SetReadDeadline(time.Now().Add(conn.rdur))
				n += rn
			}
			if rn == 0 && rerr == nil {
				err = io.EOF
				break
			} else if rerr != nil {
				break
			}
		}
		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (conn *Conn) Write(b []byte) (n int, err error) {
	if bl := len(b); bl > 0 {
		if conn.crntwdur == 0 {
			conn.crntwdur = conn.initwdur
		} else if conn.crntwdur != conn.wdur {
			conn.crntwdur = conn.wdur
		}
		conn.crntrdur = 0
		for n < bl && err == nil {
			wn, werr := conn.bufrw.Write(b[n : n+(bl-n)])
			if werr == nil {
				if err = conn.bufrw.Flush(); err == nil {
					if wn > 0 {
						n += wn
						break
					}
				}
			}
			if wn == 0 && err == nil {
				break
			}
		}
	}
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (conn *Conn) Close() (err error) {
	err = conn.conn.Close()
	return
}

// LocalAddr returns the local network address, if known.
func (conn *Conn) LocalAddr() (addr net.Addr) {
	addr = conn.conn.LocalAddr()
	return
}

// RemoteAddr returns the remote network address, if known.
func (conn *Conn) RemoteAddr() (addr net.Addr) {
	addr = conn.conn.RemoteAddr()
	return
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (conn *Conn) SetDeadline(t time.Time) (err error) {
	err = conn.conn.SetDeadline(t)
	return
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (conn *Conn) SetReadDeadline(t time.Time) (err error) {
	err = conn.conn.SetReadDeadline(t)
	return
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (conn *Conn) SetWriteDeadline(t time.Time) (err error) {
	err = conn.conn.SetWriteDeadline(t)
	return
}
