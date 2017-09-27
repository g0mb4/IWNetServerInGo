package iwnet

import (
	"fmt"
	"log"
	"net"
)

// basic server struct
type Server struct {
	Name   string       // name of the server, in case of other servers (log, CI, matchmaking)
	Conn   *net.UDPConn // UDP connection
	Logger *log.Logger  // logger
}

// close method
func (s *Server) Close() {
	s.Conn.Close()
	s.Logger.Println(s.Name + " closed")
}

// write message to a UDP address
func (s *Server) write_message(message []byte, addr *net.UDPAddr) {
	n, werr := s.Conn.WriteToUDP(message, addr)
	if werr != nil {
		s.Logger.Printf("["+s.Name+"] WriteToUDP() failed : %s\n", werr.Error())
		return
	}
	s.Logger.Printf("["+s.Name+"] written : %d bytes\n", n)
	s.print_message(message)
}

// print message in hex
func (s *Server) print_message(message []byte) {
	fmt.Printf("\n")
	for i, b := range message {
		if i == 0 {
			fmt.Printf("%s : ", s.Name)
		}

		fmt.Printf("%02X ", b)

		if i > 0 && (i+1)%16 == 0 {
			fmt.Printf("\n%s : ", s.Name)
		}
	}
	fmt.Printf("\n")
}
