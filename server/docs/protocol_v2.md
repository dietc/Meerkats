download 服务

发送0x21包 内容为空

服务端会返回一系列连续的0x21包（两个文件一起发送）

前十个包是一个文件 最后一个包是一个文件

格式

多包文件

0x21 |{"Name":"large.jpeg","Num":10,"Index":0}00000000000(一共300字节) | 65334 bytes

0x21 | 65534 bytes 

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 65534 bytes

0x21 | 5679  bytes

注意：只有第一包有json，当你读出json的Num为10的时候，再连续读取后面的9个包就可以了

单包文件

0x21 | {"Name":"large.jpeg","Num":10,"Index":0} | yes you did it

注意：单包文件直接写文件 多包文件用追加的方式
