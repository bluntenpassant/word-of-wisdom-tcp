package transport

import (
	"encoding/binary"
	"errors"
	"github.com/bluntenpassant/word-of-wisdom-tcp/pow"
	"net"
	"time"
)

type PowClientConn struct {
	transportConn net.Conn
	isConnSecure  bool
}

func NewPowClient(transportConn net.Conn) *PowClientConn {
	return &PowClientConn{
		transportConn: transportConn,
	}
}

func (c *PowClientConn) EstablishSecureConnection() error {
	packetDataDifficultyBuf := make([]byte, maxPacketSize)
	n, err := c.transportConn.Read(packetDataDifficultyBuf)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes received from server during read data packet"))
	}

	dataLength := binary.BigEndian.Uint32(packetDataDifficultyBuf[:dataLengthSize])

	data := packetDataDifficultyBuf[dataLengthSize : dataLengthSize+dataLength]
	difficultyRaw := packetDataDifficultyBuf[dataLengthSize+dataLength : n]

	difficulty := binary.BigEndian.Uint32(difficultyRaw)

	nonceRaw, err := pow.CalculateNonce(difficulty, data)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	n, err = c.transportConn.Write(nonceRaw)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes sent to server during sending nonce"))
	}

	successByteBuf := make([]byte, 1)
	n, err = c.transportConn.Read(successByteBuf)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes received from server during read success byte"))
	}

	if successByteBuf[0] != successByte {
		errBuf := make([]byte, maxErrorSize)
		n, err = c.transportConn.Read(errBuf)
		if err != nil {
			return WrapErrWithReason(errEstablishSecureConn, err)
		}

		if n == 0 {
			return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes received from server during read error"))
		}

		return WrapErrWithReason(errEstablishSecureConn, errors.New(string(errBuf)))
	}

	c.isConnSecure = true

	return nil
}

func (c *PowClientConn) Write(data []byte) (int, error) {
	return c.transportConn.Write(data)
}

func (c *PowClientConn) Read(buf []byte) (int, error) {
	return c.transportConn.Read(buf)
}

func (c *PowClientConn) LocalAddr() net.Addr {
	return c.transportConn.LocalAddr()
}

func (c *PowClientConn) RemoteAddr() net.Addr {
	return c.transportConn.RemoteAddr()
}

func (c *PowClientConn) Close() error {
	return c.transportConn.Close()
}

func (c *PowClientConn) SetWriteDeadline(t time.Time) error {
	return c.transportConn.SetWriteDeadline(t)
}

func (c *PowClientConn) SetReadDeadline(t time.Time) error {
	return c.transportConn.SetReadDeadline(t)
}

func (c *PowClientConn) SetDeadline(t time.Time) error {
	return c.transportConn.SetDeadline(t)
}
