package p

import (
    "net"
	"log"
	"bufio"
	"io"
    "util"
    "s"
)

var num int = 0

func Handle(c net.Conn){
	var state byte =  0x00
	var length uint16 = 0
	var err error
	var rbyte byte
	var rbuffer []byte
	var cbuffer []byte = make([]byte, 16)
    var flag bool
    //256*256*256
	reader := bufio.NewReaderSize(c, 16777216)
    log.Printf("Serving %s\n", c.RemoteAddr().String())
    _ = rbuffer
    _ = cbuffer
	for{
		if state != 0x0a && state != 0x0b {
            rbyte, err = reader.ReadByte()
        }
        if err != nil {
            if err == io.EOF {
                log.Printf("Connection closed by %s\n", c.RemoteAddr().String())
                c.Close()
                return
            }
            continue
        }
		switch state {
		    case 0x00:
                state = updateState(rbyte, 0x11, state)
			case 0x01:
                state = updateState(rbyte, 0xff, state)
            case 0x02:
                state = updateState(rbyte, 0x6c, state)
            case 0x03:
                state = updateState(rbyte, 0x6f, state)
            case 0x04:
                state = updateState(rbyte, 0x6e, state)
            case 0x05:
                state = updateState(rbyte, 0x64, state)
            case 0x06:
                state = updateState(rbyte, 0x6f, state)
            case 0x07:
                state = updateState(rbyte, 0x6e, state)
            case 0x08:
                /** packet length **/
                length += uint16(rbyte) * 256
                state++
            case 0x09:
                length += uint16(rbyte)
                rbuffer = make([]byte, length)
                state++
            case 0x0a:
                log.Println(int(length))
                rbuffer, err = reader.Peek(int(length))
                //buffer full?
                if err != nil {
                    log.Println(1, err)
                    state, length, rbuffer = restoreState()
                    continue
                }
                log.Println(len(rbuffer))
                _, err := reader.Discard(int(length))
                if err != nil {
                    log.Println(2, err)
                    state, length, rbuffer = restoreState()
                    continue
                }
                state++
            case 0x0b:
                cbuffer, err = reader.Peek(16)
                //fmt.Println("received checksum:", cbuffer)
                if err != nil {
                    log.Println(3, err)
                    state, length, rbuffer = restoreState() 
                    continue
                }
                log.Printf("%x", cbuffer)
                flag = util.CheckMd5(cbuffer, rbuffer, []byte("aaaaa"))
                if flag {
                    discarded, err := reader.Discard(16)
                    if err != nil || discarded != 16 {
                        state, length, rbuffer = restoreState() 
                        continue
                    }
                    state++
                } else {
                    state, length, rbuffer = restoreState()
                    continue
                }
            case 0x0c:
                state = updateState(rbyte, 0xff, state)
            case 0x0d:
                if rbyte == 0xee{
                     handlePacket(c, rbuffer)
                }
                state, length, rbuffer = restoreState() 
		}
	}
}

func handlePacket(c net.Conn, packet []byte) {
    switch packet[0] {
        case 0x01:
        //testing packet
        send(c, enPacket(s.ProcessTestingPacket(c, packet[2:], packet[1]), 0x01))
        case 0x02:
        //pull packet
        send(c, enPacket(s.ProcessLocalListPacket(packet[2:], packet[1]), 0x02))
        case 0x20:
        //client upload file
        var flag bool 
        flag = s.ProcessFileAssembly(packet[2:], packet[1]) 
        if flag {
            send(c, enPacket([]byte("ok"), 0x20))
        }
        case 0x21:        
        //client download
        var files [][]byte = s.ProcessFileDownload(packet[2:], packet[1])
        for _,part := range files{
           //log.Println(len(part))
           send(c, enPacket(part, 0x21))
        }
        case 0x22:
        //client upload modified file(rsync)
        //flag := s.ProcessFileAssembly(packet[2:], packet[1])
        //if flag {
        //    send(c, enPacket([]byte("ok"), 0x22))
        //}
        var files [][]byte = s.ProcessModifiedFileUpload(packet[2:], packet[1])
        for _,part := range files{
           send(c, enPacket(part, 0x22))
        }
        case 0x23:
        //client download modified file(rsync)
        var files [][]byte = s.ProcessFileDownload(packet[2:], packet[1])
        for _,part := range files{
           //log.Println(len(part))
           send(c, enPacket(part, 0x23))
        }
    }
}

func enPacket(data []byte, typ byte) []byte{
     var cl uint16 = uint16(len(data) + 1)
     var tl int = int(cl) + 28
     var res []byte = make([]byte, tl)
     res[0] = 0x11
     res[1] = 0xff
     res[2] = 0x6c
     res[3] = 0x6f
     res[4] = 0x6e
     res[5] = 0x64
     res[6] = 0x6f
     res[7] = 0x6e
     res[8] = byte(cl >> 8)
     res[9] = byte(cl & 0xFF)
     //packet type 
     res[10] = typ
     copy(res[11:], data)
     var checksum [16]byte = util.Md5(res[10:10+uint32(cl)], []byte("aaaaa"))
     for i:=0;i<16;i++ {
         res[tl-18+i] = checksum[i]
     }
     res[tl-2] = 0xff
     res[tl-1] = 0xee
     return res
}

func send(c net.Conn, b []byte) {
    bufferWriter := bufio.NewWriter(c)
    bufferWriter.Write(b)
    bufferWriter.Flush()
}

func updateState(r byte, target byte, state byte) (retstate byte) {
    if r == target {
        state++
    } else {
        state = 0x00
    }
    retstate = state
    return
}

func restoreState() (a byte, b uint16, c []byte) {
    a = 0 
    b = 0
    c = nil
    return
}