### Meerkats Communication Protocol Over TCP

- Packet Initiator [ 8 bytes ]
**0x11 0xff 0x6c 0x6f 0x6e 0x64 0x6f 0x6e**
(which is "/x11/xff london")

- Packet Content Length [ 2 bytes ]
uint16 cl; 
// cl = len(data) + 1, data is what you want to transmit in this packet, //1-byte-space is for byte type
uint8 cl[2];
cl[0] = uint8(cl >> 8);
cl[1] = uint8(cl & 0xff)

- Packet Type [ 1 byte ]
uint8

- Packet Content [ len(data) bytes ]
uint8 buf[len(data)]

- Md5 Checksum [16 bytes]

  buffer[0]  packet type
  
  buffer[1]  device id
  
  buffer[2] ~ buffer[n]  packet content
  
  buffer[n+1] ~ buffer[n+5] private key ("aaaaa" for testing)
  
  checksum = md5(buffer)

- Packet Terminator [ 2 bytes ]
  **0xff 0xee**

### A Structured Packet(sending from client)
| initiator(8) |  length(2) | packet type(1) | device id(1) | packet content(n) | checksum(16) | terminator(2)|

example:
| 11ff6c6f6e646f6e | 0007 | 01 | 01 | 68656c6c6f | e623934ab23681c042e3ee5eae36b518 | ffee |

A packet that comes from server is similar in structure except that there is no device id (1 byte)  

    
### Pseudocode For Enpacketing

```
    *uint8 enpacket(uint8 data[], uint8 p_type, uint8 device_id) {
        uint16 content_length = len(data) + 2
        uint16 total_length = content_length + 28
        uint8 buffer[total_length]
        uint8 checksum_source[len(data) + 6]
        uint8 checksum[16]
        /** initiator **/
        buffer[0] = 0x11
        buffer[1] = 0xff
        buffer[2] = 0x6c
        buffer[3] = 0x6f
        buffer[4] = 0x6e
        buffer[5] = 0x64
        buffer[6] = 0x6f
        buffer[7] = 0x6e
        /** length **/
        buffer[8] = uint8(content_length >> 8)
        buffer[9] = uint8(content_length & 0xff)
        /** type **/
        buffer[10] = p_type
        /** device id **/
        buffer[11] = device_id
        /** content **/
        memcpy(buffer + 12, data, len(data))
        /** calculate checksum and put it inside packet **/
        checksum_source[0] = p_type
        checksum_source[1] = device_id
        memcpy(checksum_source + 2, data, len(data))
        checksum_source[len(data)+2] = 'a'  // "aaaaa" private key
        checksum_source[len(data)+3] = 'a'
        checksum_source[len(data)+4] = 'a'
        checksum_source[len(data)+5] = 'a'
        checksum_source[len(data)+6] = 'a'
        checksum = md5(checksum_source)
        memcpy(buffer + 12 + len(data), checksum, len(checksum))
        buffer[total_length-2] = 0xff
        buffer[total_length-1] = 0xee  
        return buffer
    }
```
### Pseudocode For Depacketing

```
       Connection con; /** established connection **/
       uint8 byte 
       int state = 0x00
       uint16 length = 0
       uint8 content_buffer[]
       while(1) {
           if state != 0x0a && state != 0x0b {
               try{
                      byte = con.ReadByte()
               } catch (Exception EOF) {
                       return 
               }
            }
            switch state{
                case 0x00:
                    if byte == 0x11 {
                        state ++
                    }
                    break
                case 0x01:
                    if byte == 0xff {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x02:
                    if byte == 0x6c {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x03:
                    if byte == 0x6f {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x04:
                    if byte == 0x6e {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x05:
                    if byte == 0x64 {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x06:
                    if byte == 0x6f {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x07:
                    if byte == 0x6e {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x08:
                    length += uint16(byte) * 256
                    state++
                    break
                case 0x09:
                    length += uint16(byte)
                    state++
                    content_buffer = malloc(sizeof(uint8)*length)
                    break
                case 0x0a:
                    data = con.ReadBytesByLength(length)                    memcpy(content_buffer, data, length)
                    /** so packet_type + packet_content is stored in content_buffer now, then you can resolve it **/
                    state++
                    break
                case 0x0b:
                    received_checksum = con.ReadBytesByLength(length)
                    calculated_checksum = md5(content_buffer + "aaaaa")
                    if received_checksum == calculated_checksum {
                        state ++
                    } else {
                        state = 0x00
                        length = 0
                        memset(content_buffer, 0) /** clear buffer **/
                    }
                    break
                case 0x0c:
                    if byte == 0xff {
                        state ++
                    } else {
                        state = 0x00
                    }
                    break
                case 0x0d:
                    if byte == 0xee {
                        process_data(content_buffer)
                    }
                    break
            }
            if state == 0x0d{
                 state = 0x00
                 length = 0
                 memset(content_buffer, 0) /** clear buffer **/
            }
       }
```
### Start From Now
socket programming over TCP:
ip: 178.128.45.7
port: 4356

Try to connect to the server and send a packet which is enpacketed with 
- data "hello"
- packet type 0x01
- device id  0x01

After that, catch a packet and resolve it. 
If packet type = 0x01 and packet content = "world", then we can continue.

