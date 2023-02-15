package manager

import (
	"errors"
	"fmt"
	"github.com/bluntenpassant/word-of-wisdom-tcp/transport"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

const defaultDifficulty = 1
const defaultReadDeadline = 5 * time.Minute

const defaultBurstCapacity = 1000

type Manager struct {
	currentConnectionBurst uint32
	quotes                 []string

	difficulty uint32
}

func NewManager(quotes []string) *Manager {
	return &Manager{
		quotes:                 quotes,
		difficulty:             defaultDifficulty,
		currentConnectionBurst: 0,
	}
}

func (m *Manager) Run(address string) error {
	connCh, err := listenClients(address)
	if err != nil {
		return err
	}

	for {
		conn := <-connCh

		difficulty := m.manageAndGetDifficulty()
		m.increaseBurst()

		powConn := transport.NewServerClient(conn, difficulty)

		err := powConn.EstablishSecureConnection()
		if err != nil {
			fmt.Println("error establishing secure conn: " + err.Error())
		}

		err = m.handleConn(conn)
		if err != nil {
			fmt.Println("error handling conn cause: " + err.Error())
		}
	}
}

func (m *Manager) manageAndGetDifficulty() uint32 {
	currentBurst := atomic.LoadUint32(&m.currentConnectionBurst)
	if currentBurst > defaultBurstCapacity {
		atomic.StoreUint32(&m.currentConnectionBurst, 0)
		atomic.AddUint32(&m.difficulty, 1)
	}

	return atomic.LoadUint32(&m.difficulty)
}

func (m *Manager) increaseBurst() {
	atomic.AddUint32(&m.currentConnectionBurst, 1)
}

func (m *Manager) handleConn(conn net.Conn) error {
	randQuote := m.quotes[rand.Intn(len(m.quotes))]
	n, err := conn.Write([]byte(randQuote))
	if err != nil {
		return errors.New("error handling conn cause: " + err.Error())
	}

	if n == 0 {
		return errors.New("zero bytes sent to client")
	}

	return nil
}

func listenClients(address string) (<-chan net.Conn, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	clientsCh := make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("error accept conn: " + err.Error())
			}

			err = conn.SetReadDeadline(time.Now().Add(defaultReadDeadline))
			if err != nil {
				panic(err)
			}

			clientsCh <- conn
		}
	}()

	return clientsCh, nil
}
