package iwnet

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

const IP_SERVER_PORT = 1500

// client message
type IPRequestPacket struct {
	Type1, Type2, Type3   uint8
	RequestType, Sequence uint16
	XUID                  uint64
}

// server response
type IPResponsePacket struct {
	SourceAddr net.IP // [4] byte
	SourcePort int
	Sequence   uint16
	NatOpen    uint8
}

// create new server pointer
func NewIPServer(logger *log.Logger) *Server {
	// create UDP server
	srvaddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(IP_SERVER_PORT))
	if err != nil {
		logger.Fatal("[IPServer] ResolveUDPAddr() failed : " + err.Error())
	}

	// start UDP server
	conn, err := net.ListenUDP("udp", srvaddr)
	if err != nil {
		logger.Fatal("[IPServer] ListenUDP() failed : " + err.Error())
	}

	// fill up the server struct
	server := &Server{
		Name:   "IPServer",
		Conn:   conn,
		Logger: logger,
	}

	return server
}

// start the server
func (s *Server) StartAndRun() {
	s.Logger.Println(s.Name + " started")
	go s.handle_IP_request() // start client handle goroutine
}

// handling IP requests
func (s *Server) handle_IP_request() {
	var packet IPRequestPacket
	data := make([]byte, 64)
	buf := bytes.NewReader(data)
	for {
		_, addr, errr := s.Conn.ReadFromUDP(data) // read message
		if errr != nil {
			s.Logger.Println("[" + s.Name + "] ReadFromUDP() failed : " + errr.Error())
		}

		errc := binary.Read(buf, binary.LittleEndian, &packet) // put message in the stuct
		// in case of other message type
		if errc != nil {
			//s.Logger.Println("[" + s.Name + "] binary.Read() failed : " + errc.Error())
			continue
		}

		s.Logger.Printf(s.Name+" packet : %02X\n", packet.Type3)
		// only type we need
		if packet.Type3 == 0x14 {
			s.Logger.Printf("["+s.Name+"] Handling IP request from %X (%s)\n", packet.XUID, addr.IP.String())

			var response IPResponsePacket
			response.SourceAddr = addr.IP.To4()
			response.SourcePort = addr.Port
			response.Sequence = packet.Sequence
			response.NatOpen = 0

			// create the response message
			message := s.create_IP_response(response)
			s.write_message(message, addr)

			// uhm.... stuff, but needed
			if addr.Port >= 28960 && addr.Port <= 29960 {
				response.NatOpen = 1

				message := s.create_IP_response(response)
				s.write_message(message, addr)
			}
		}
	}
}

// pack the server message
func (s *Server) create_IP_response(response IPResponsePacket) []byte {
	message := new(bytes.Buffer)

	message.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF}) // magic bytes

	if response.NatOpen == 1 {
		message.Write([]byte("OpenNAT"))
	} else {
		message.Write([]byte("ipdetect"))
	}

	message.Write([]byte{0x00, 0x00}) // magic delimeter bytes ?

	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, response.Sequence) // big endian but why ...
	//binary.LittleEndian.PutUint16(b, response.Sequence)
	message.Write(bs)

	message.Write([]byte{0x00, 0x14, 0x1B}) // magic bytes

	// reverse !!
	// IP: 127.0.0.1 -> 01 00 00 7F
	message.Write([]byte{response.SourceAddr[3], response.SourceAddr[2], response.SourceAddr[1], response.SourceAddr[0]})

	// port: 28960 -> 20 71 00 00
	bp := make([]byte, 4)
	binary.LittleEndian.PutUint32(bp, uint32(response.SourcePort))
	message.Write(bp)

	message.Write([]byte{0x42, 0x37, 0x13, 0x37, 0x13, 0x42}) // magic bytes
	message.Write([]byte{0x40, 0x01})                         // magic bytes
	return message.Bytes()
}
