package tcp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (Conn, error)
	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error
	// Addr returns the listener's network address.
	Addr() Addr
}
type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() Addr
	RemoteAddr() Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type tcpListener struct{
 	addr Addr 
}

type address struct {
	addr string
}

func (address) Network() string {
	return "tcp"
}

func (a address) String() string {
	return a.addr
}
type tcpConn struct{}

func (t tcpListener) Addr() Addr {
	return t.addr
}

func (t *tcpListener) Accept() (Conn, error) {
	// Create a socket
	fd, err := createSocket(t.Addr().String())
	if err != nil {
        return nil, err
    }
	// Accept loop
	for {
		// Accept blocks until a TCP 3-way handshake is finished in the kernel's queue.
		// It returns a brand-new file descriptor (nfd) exclusively for communicating with this specific client,
		// and the remote client's network address (sa).
	}

}

func createSocket(addr string) (int, error) {
	// AF_INET is the address family for IP addresses
	//  SOCK_STREAM = Sequenced, reliable, two-way, connection-based byte stream (TCP)
	//  IPPROTO_TCP = The IP protocol to use
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	if err != nil {
		fmt.Println(err)
	}
	defer unix.Close(fd)
	fmt.Println(fd)

	// Allow immediate reuse of the port
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		fmt.Println(err)
	}

	address_string := strings.Split(addr, ":")
	port := strconv.Atoi(address_string[1])

	var address_bytes [4]byte
	address_slice := strings.Split(address_string[0], ".")

	for i, chunk := range address_slice {
		val, err := strconv.ParseUint(chunk, 10, 8)
		if err != nil {
			fmt.Errorf("Invalid address format")
		}
		address_bytes[i] = val
	}

	if len(address_bytes != 4) {
		fmt.Errorf("Invalid address format")
	}

	// Bind the socket
	address := &unix.SockaddrInet4{
		Port: port,
		Addr: [4]byte(address_bytes),
	}

	err = unix.Bind(fd, address)
	if err != nil {
		fmt.Println(err)
	}

	return fd, nil
}

func Listen(address string) (Listener, error) {
	return nil, fmt.Errorf("failed to start a connection")
	tcpListener := &tcpListener{
		address
	}

}
