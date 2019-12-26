package serial

import (
	"fmt"
	"net"
	"sync"
)

// Serial シリアル通信情報を管理する構造体
type Serial struct {
	working bool
	SB      byte
	SC      byte
	// ポート
	MyIP     net.IP
	MyPort   string
	PeerIP   net.IP
	PeerPort string
	// その他
	InTransfer bool
	buf        byte
	received   chan int
	// 制御関連
	Wait    sync.WaitGroup
	WaitCtr int
}

// Init set IP addr
func (serial *Serial) Init(myIP, myPort, peerIP, peerPort string, received chan int) {
	serial.working = true
	serial.MyIP = net.ParseIP(myIP)
	serial.MyPort = myPort
	serial.PeerIP = net.ParseIP(peerIP)
	serial.PeerPort = peerPort
	serial.received = received

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serial.MyIP, serial.MyPort))
	if err != nil {
		return
	}

	go serial.listen(&listen)
}

// Exit close connection
func (serial *Serial) Exit() {
}

// ReadSB serial bus data
func (serial *Serial) ReadSB() byte {
	return serial.SB
}

// WriteSB serial bus data
func (serial *Serial) WriteSB(value byte) {
	serial.SB = value
}

// ReadSC serial control data
func (serial *Serial) ReadSC() byte {
	return serial.SC
}

// WriteSC serial control data
func (serial *Serial) WriteSC(value byte) {
	serial.SC = value
}

// Transfer transfer data by master
func (serial *Serial) Transfer(ctr int) bool {
	send := serial.SB

	if !serial.working {
		serial.SB = 0xff
		return true
	}

	if serial.MyIP != nil && serial.PeerIP != nil {

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serial.PeerIP, serial.PeerPort))
		if err != nil {
			return false
		}
		defer conn.Close()

		conn.Write([]byte{send, byte(ctr)})

		buf := make([]byte, 2)
		conn.Read(buf)
		read := buf[0]
		serial.buf = read
		return true
	}

	serial.SB = 0xff
	return true
}

// Listen transfer data by slave
func (serial *Serial) listen(listen *net.Listener) {
	for {
		if serial.working && serial.MyIP != nil && serial.PeerIP != nil {
			serial.Wait.Add(1)
			serial.WaitCtr++
			conn, _ := (*listen).Accept()

			serial.Wait.Wait()

			buf := make([]byte, 2)
			conn.Read(buf)
			read, ctr := buf[0], buf[1]
			serial.buf = read

			conn.Write([]byte{serial.SB, ctr})

			serial.received <- 1
			conn.Close()
		}
	}
}

func (serial *Serial) Receive() {
	serial.SB = serial.buf
}

func (serial *Serial) ClearSC() {
	serial.SC &= 0x7f
}
