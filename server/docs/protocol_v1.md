### 此文档会不定期更新

## 03/03 TASK
发送pull包 -> 根据响应包的指令操作本地目录（目前只有cmd=1的文件上传操作） -> 检查文件名与checksum是否一致

    ->如果与本地同文件名文件的checksum不一致 -> 不做任何操作（因为文件在你发pull包之后有新的修改）

    ->如果与本地同文件名文件的checksum一致 -> 发送单个或多个upload包上传该文件 -> 接收到upload响应包 -> 任务完成


## 协议细节 

log：


03/03/2019 服务器端增加列表上传功能，文件上传功能



### 今后的测试，客户端使用不同的device id：


- desktop 0x02


- android 0x03



### 客户端发包格式：

| 包名           | 功能            | typeid           | 内容                    |
| ------------- | -------------  | ---------------- | ----------------------- |
| pull          | 上传本地目录列表  |  0x02            | 扫描目录结构得到的json[1]  |
| upload        | 上传完整文件     |  0x20            | 文件信息json+文件内容字节[2]|


> note

1.上表内容是指 typeid|deviceid 之后， checksum 之前的部分

> detail

[1]example：

- 本地的真实目录结构：

根目录下有 **文件夹dir1** ， **文件夹dir2** 和 **文件pic1.jpg** ,其中文件夹dir1中有文件 **doc1**

- 生成的json字符串形式和字节形式如下

```
## json 
## Typ1->文件  Typ2->目录
## Digest->md5(文件)

{
	"dir": {
		"Name": "dir",
		"Typ": 2,
		"Digest": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
	},
	"dir/doc1": {
		"Name": "dir/doc1",
		"Typ": 1,
		"Digest": [62, 238, 247, 95, 134, 12, 98, 97, 1, 23, 50, 153, 132, 75, 115, 133]
	},
	"dir2": {
		"Name": "dir2",
		"Typ": 2,
		"Digest": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
	},
	"pic1.jpg": {
		"Name": "pic1.jpg",
		"Typ": 1,
		"Digest": [141, 221, 57, 244, 109, 128, 93, 240, 169, 147, 223, 23, 244, 177, 183, 2]
	}
}

## 字节（给你测试用）
7b22646972223a7b224e616d65223a22646972222c22547970223a322c22446967657374223a5b302c302c302c
302c302c302c302c302c302c302c302c302c302c302c302c305d7d2c226469722f646f6331223a7b224e616d65
223a226469722f646f6331222c22547970223a312c22446967657374223a5b36322c3233382c3234372c39352c
3133342c31322c39382c39372c312c32332c35302c3135332c3133322c37352c3131352c3133335d7d2c226469
7232223a7b224e616d65223a2264697232222c22547970223a322c22446967657374223a5b302c302c302c302c
302c302c302c302c302c302c302c302c302c302c302c305d7d2c22706963312e6a7067223a7b224e616d65223a
22706963312e6a7067222c22547970223a312c22446967657374223a5b3134312c3232312c35372c3234342c31
30392c3132382c39332c3234302c3136392c3134372c3232332c32332c3234342c3137372c3138332c325d7d7d

## 目录为空时
{}
7b7d

```

[2]example

```
1、对于收到pull响应包中 cmd=1 的文件，发送upload包
2、对于一个文件，可能发送一个upload包，也可能发送多个upload包，取决于文件的大小
3、一个upload包传送的最大文件字节是65233，例如：一个文件只有10个字节，那么发送一个upload包即可;如果一个文件有65243个字节，
  则需要发送两个upload包，第一个65233个字节，第二个10个字节
4、对于单个文件拆成的多个包，必须单次连续发出
5、每个upload包前面都包含一个文件信息json（一定不足300字节），并用0x00在json后填充完整300字节
6、上述json样例
   ## Num 是这个文件的分包数量
   ## Index 是当前包的位置索引
   单包发送的文件：
   {
	    "Name": "dir/doc1",
	    "Num": 1,           
	    "Index": 0
   }
   多包发送的文件:
   {
	    "Name": "pic1.jpg",
	    "Num": 2,       
	    "Index": 0
   }
   {
	    "Name": "pic1.jpg",
	    "Num": 2,
	    "Index": 1
   }
 7、一个完整的upload包的内容部分是
  |一个json的字节形式| 0x00 0x00 0x00|   文件字节  |
  {           300bytes            }{  <=65233  }
```






### 客户端接收包格式：

| 包名           | 功能                                         | typeid           | 内容                    |
| ------------- | --------------------------------------------| ---------------- | ----------------------- |
| pull响应包     | 发送pull包之后收到，得到本地操作指令               |  0x02            | 本地操作指令的json[1]      |
| upload响应包   | 发送upload包之后收到,收到表示上传成功，不用做任何操作 |  0x20           |  "ok"                    |

> detail

1.example：

```
##用于本地操作的指令json
## Cmd: = 1  -> 整体上传该文件（发送upload包）
## Ext目前为空，是预留字段，以后会用到，这次测试不用考虑

{
	"dir/doc1": {
		"Name": "dir/doc1",
		"Digest": [62, 238, 247, 95, 134, 12, 98, 97, 1, 23, 50, 153, 132, 75, 115, 133],
		"Cmd": 1,
		"Ext": ""
	},
	"pic1.jpg": {
		"Name": "pic1.jpg",
		"Digest": [141, 221, 57, 244, 109, 128, 93, 240, 169, 147, 223, 23, 244, 177, 183, 2],
		"Cmd": 1,
		"Ext": ""
	}
}

##对应字节
7b226469722f646f6331223a7b224e616d65223a226469722f646f6331222c22446967657374223a5b36322c3233
382c3234372c39352c3133342c31322c39382c39372c312c32332c35302c3135332c3133322c37352c3131352c31
33335d2c22436d64223a312c22457874223a22227d2c22706963312e6a7067223a7b224e616d65223a2270696331
2e6a7067222c22446967657374223a5b3134312c3232312c35372c3234342c3130392c3132382c39332c3234302c
3136392c3134372c3232332c32332c3234342c3137372c3138332c325d2c22436d64223a312c22457874223a2222
7d7d

## 同样需要处理指令为空的情况(不在本地做任何操作)
{}
7b7d
```
