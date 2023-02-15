package transport

import (
	"encoding/binary"
	"errors"
	"github.com/bluntenpassant/word-of-wisdom-tcp/pow"
	"math/rand"
	"net"
	"time"
)

type PowServerConn struct {
	transportConn net.Conn
	isConnSecure  bool
	difficulty    uint32
}

func NewServerClient(transportConn net.Conn, difficulty uint32) *PowServerConn {
	return &PowServerConn{
		transportConn: transportConn,
		difficulty:    difficulty,
	}
}

func (c *PowServerConn) EstablishSecureConnection() error {
	randomData, err := generateRandomData(minDataLength, maxDataLength)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	dataDifficultyBuf := make([]byte, dataLengthSize+len(randomData)+difficultySize)
	dataLengthSizeBuf := make([]byte, dataLengthSize)
	binary.BigEndian.PutUint32(dataLengthSizeBuf, uint32(len(randomData)))

	copy(dataDifficultyBuf, dataLengthSizeBuf)
	copy(dataDifficultyBuf[dataLengthSize:], randomData)

	rawDifficultyBuf := make([]byte, difficultySize)
	binary.BigEndian.PutUint32(rawDifficultyBuf, c.difficulty)

	copy(dataDifficultyBuf[dataLengthSize+len(randomData):], rawDifficultyBuf)

	n, err := c.transportConn.Write(dataDifficultyBuf)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes sent to client during write data packet"))
	}

	nonceBuf := make([]byte, nonceLength)

	n, err = c.transportConn.Read(nonceBuf)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes got from client during reading nonce"))
	}

	err = pow.IsValidNonce(randomData, nonceBuf, c.difficulty)
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	n, err = c.transportConn.Write([]byte{successByte})
	if err != nil {
		return WrapErrWithReason(errEstablishSecureConn, err)
	}

	if n == 0 {
		return WrapErrWithReason(errEstablishSecureConn, errors.New("zero bytes sent to client during write success byte"))
	}

	c.isConnSecure = true

	return nil
}

func (c *PowServerConn) Write(data []byte) (int, error) {
	return c.transportConn.Write(data)
}

func (c *PowServerConn) Read(buf []byte) (int, error) {
	return c.transportConn.Read(buf)
}

func (c *PowServerConn) LocalAddr() net.Addr {
	return c.transportConn.LocalAddr()
}

func (c *PowServerConn) RemoteAddr() net.Addr {
	return c.transportConn.RemoteAddr()
}

func (c *PowServerConn) Close() error {
	return c.transportConn.Close()
}

func (c *PowServerConn) SetWriteDeadline(t time.Time) error {
	return c.transportConn.SetWriteDeadline(t)
}

func (c *PowServerConn) SetReadDeadline(t time.Time) error {
	return c.transportConn.SetReadDeadline(t)
}

func (c *PowServerConn) SetDeadline(t time.Time) error {
	return c.transportConn.SetDeadline(t)
}

func generateRandomData(min int, max int) ([]byte, error) {
	dataLength := rand.Intn(max-min) + min

	randDataBuf := make([]byte, dataLength)

	n, err := rand.Read(randDataBuf)
	if err != nil {
		return nil, errors.New("error generating random data: " + err.Error())
	}

	if n == 0 {
		return nil, errors.New("error generating random data: zero bytes received from rand reader")
	}

	return randDataBuf, nil
}
