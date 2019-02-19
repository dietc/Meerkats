package s

import(
    "net"
	//"bufio"
	//"log"
)

func ProcessTestingPacket(c net.Conn, data []byte, device byte) []byte{
	if string(data) == "hello" && device == 0x01 {
        return []byte("world")
	}
    return nil
}
