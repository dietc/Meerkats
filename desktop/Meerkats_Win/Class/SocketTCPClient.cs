using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;
using System.Security.Cryptography;
using System.IO;
using System.Threading;
using System.Diagnostics;
using Newtonsoft.Json.Linq;
using Newtonsoft.Json;
using System.Windows.Documents;
using static Meerkats_Win.File_json_info;
using Meerkats_Win.Class;

namespace Meerkats_Win.Class
{

    class SocketTCPClient
    {
        Func lib = new Func();
        // 10.40.157.245
        // 
        //private static string ip = "178.128.45.7";
        private static string ip = "10.40.157.245";
        private static int port = 4356;
        private static Socket socketClient;

        private byte Device_id = 0x2;


        // Socket_Buffer_Size
        int Buffer_Length_Max = 1024;


        // the path for stored data 
        private static string PATH = System.AppDomain.CurrentDomain.BaseDirectory + "sync_disk\\";
        private static string Backup_PATH = System.AppDomain.CurrentDomain.BaseDirectory + "back_history_file\\";

        public static List<string> listMessage = new List<string>();

        ///<summary>
        ///create a SocketClient Instance
        ///</summary>
        ///ip address port type = TCP
        public void CreateInstance()
        {
            socketClient = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            ConnectServer();

        }
        /// <summary>
        /// Connect to server
        /// </summary>
        public void ConnectServer()
        {
            try
            {
                socketClient.Connect(IPAddress.Parse(ip), port);

            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }

        }
        public void DisconnectServer()
        {
            try
            {
                socketClient.Close();
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }
        }
        /// <summary>
        /// receive Msg
        /// </summary>
        public byte[] ReceiveMessage()
        {
            try
            {
                // 8 + 2
                int Tcp_header_length = 10;
                int Tcp_body_length = 0;

                byte[] Tcp_header = new byte[Tcp_header_length];

                socketClient.Receive(Tcp_header, 0, Tcp_header_length, 0);
                Tcp_body_length = (Tcp_header[8] << 8) + (Tcp_header[9]);

                byte[] recvBytes = new byte[Tcp_body_length];

                // receive data
                int index = 0;

                // Buffer_Length_Max = 1024 * 6
                while (true)
                {
                    if (Tcp_body_length > Buffer_Length_Max)
                    {
                        socketClient.Receive(recvBytes, index * Buffer_Length_Max, Buffer_Length_Max, 0);
                        Tcp_body_length -= Buffer_Length_Max;
                    }

                    else
                    {
                        socketClient.Receive(recvBytes, index * Buffer_Length_Max, Tcp_body_length, 0);
                        break;
                    }

                    index++;

                    // wair for 1ms 
                    // .Net have the speed limit
                    // control received speed in order not to loss packet
                    System.Threading.Thread.Sleep(1);
                }

                byte[] md5 = new byte[16];
                socketClient.Receive(md5, 0, 16, 0);

                byte[] End_flag = new byte[2];
                socketClient.Receive(End_flag, 0, 2, 0);
                // check End_flag

                if (End_flag[0] == 0xff && End_flag[1] == 0xee)
                {
                    // socketClient.Close();
                    return (UnpackData(recvBytes));
                }

                else
                {
                    // socketClient.Close();
                    return null;
                }
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }


        }

        public string ReceiveMessage_For_download(int file_type)
        {

            try
            {

                // timer
                System.Diagnostics.Stopwatch stopwatch = new Stopwatch();
                stopwatch.Start(); //  start

                byte[] recvBytes;

                byte Packet_Num;
                string File_Name;

                // 8 + 2
                int Tcp_header_length = 10;
                int Tcp_body_length = 0;

                byte[] Tcp_header = new byte[Tcp_header_length];

                socketClient.Receive(Tcp_header, 0, Tcp_header_length, 0);

                Tcp_body_length = ((int)Tcp_header[8] << 8) + ((int)Tcp_header[9]);

                // Context Length
                byte[] Packet_type = new byte[1];
                int File_data_Length;

                byte[] file_json = new byte[300];
                // 65234
                File_data_Length = Tcp_body_length - 300 - 1;
                recvBytes = new byte[File_data_Length];
                socketClient.Receive(Packet_type, 0, 1, 0);
                socketClient.Receive(file_json, 0, 300, 0);

                // get Packet Num
                int End_file_json_flag = Array.IndexOf(file_json, (byte)0x00);
                byte[] File_json = new byte[End_file_json_flag];
                Buffer.BlockCopy(file_json, 0, File_json, 0, End_file_json_flag);
                file_index_json F_json = JsonConvert.DeserializeObject<file_index_json>(System.Text.Encoding.Default.GetString(File_json));
                // return File_Name + Packet_Num
                File_Name = F_json.Name;
                Packet_Num = F_json.Num;

                //create file or find file
                string filePath = PATH + File_Name;
                FileStream fs = new FileStream(filePath, FileMode.Create, FileAccess.Write);


                // receive file data
                int index = 0;
                // set Buffer size = 1024

                while (true)
                {
                    if (File_data_Length > Buffer_Length_Max)
                    {
                        socketClient.Receive(recvBytes, index * Buffer_Length_Max, Buffer_Length_Max, 0);
                        File_data_Length -= Buffer_Length_Max;
                    }

                    else
                    {
                        socketClient.Receive(recvBytes, index * Buffer_Length_Max, File_data_Length, 0);
                        break;
                    }

                    index++;

                    // wair for 1ms 
                    // .Net have the speed limit
                    // control received speed in order not to loss packet
                    System.Threading.Thread.Sleep(1);
                }


                byte[] md5 = new byte[16];
                socketClient.Receive(md5, 0, 16, 0);

                byte[] End_flag = new byte[2];
                socketClient.Receive(End_flag, 0, 2, 0);
                // check End_flag


                if (End_flag[0] == 0xff && End_flag[1] == 0xee)
                {
                    // socketClient.Close();
                    switch (file_type)
                    {
                        case 0:
                            // the whole file --download
                            fs.Write(recvBytes, 0, recvBytes.Length);
                            fs.Position = fs.Length;
                            break;
                        case 1:
                            // differ download

                            fs.Write(recvBytes, 0, recvBytes.Length);
                            fs.Position = fs.Length;
                            break;
                    }
                    

                    if (Packet_Num == 1)
                    {
                        stopwatch.Stop(); //  stop watch
                        TimeSpan timespan = stopwatch.Elapsed; // Get the elapsed time as a TimeSpan value.
                                                               // double hours = timespan.TotalHours; // hours
                                                               // double minutes = timespan.TotalMinutes;  // Minutes
                                                               // double seconds = timespan.TotalSeconds;  //  Seconds
                        double milliseconds = timespan.TotalMilliseconds;  //  Milliseconds
                        fs.Flush();
                        double fs_size = fs.Length;
                        double speed = fs_size / milliseconds * 1000 / 1024; // kb/s 

                        fs.Close();

                        return File_Name + ": download successfully";
                    }

                    else
                    {
                        while (Packet_Num > 1)
                        {
                            byte[] file_data = UnpackData(ReceiveMessage());
                            switch (file_type)
                            {
                                case 0:
                                    // the whole file --download
                                    fs.Write(file_data, 0, file_data.Length);
                                    fs.Position = fs.Length;
                                    break;
                                case 1:
                                    // differ download
                                    fs.Write(file_data, 0, file_data.Length);
                                    fs.Position = fs.Length;
                                    break;
                            }
                            
                            Packet_Num--;
                        }

                        stopwatch.Stop(); // stop watch
                        TimeSpan timespan = stopwatch.Elapsed; // Get the elapsed time as a TimeSpan value.
                                                               // double hours = timespan.TotalHours; // hours
                                                               // double minutes = timespan.TotalMinutes;  // Minutes
                                                               // double seconds = timespan.TotalSeconds;  //  Seconds
                        double milliseconds = timespan.TotalMilliseconds;  //  Milliseconds
                        fs.Flush();
                        double fs_size = fs.Length;
                        double speed = fs_size / milliseconds * 1000 / 1024; // kb/s 

                        fs.Close();

                        return File_Name + ": download successfully";

                    }
                
                }

                else
                {
                    stopwatch.Stop(); //  stop watch
                    // socketClient.Close();
                    return null;
                }
            }
            catch (Exception ex)
            {

                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }

        }

        ///<summary>
        ///Send Msg
        ///</summary>
        public void SendMessage(byte[] sendBytes)
        {
            try
            {
                //check if connected
                if (socketClient.Connected)
                {
                    // get RemoteEndPoint info <ip,port>
                    IPEndPoint ipe = (IPEndPoint)socketClient.RemoteEndPoint;
                    socketClient.Send(sendBytes, sendBytes.Length, 0);
                }
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }

        }

        public string Download_File(string file_name, int download_type)
        {

            
            byte[] Msgbody = System.Text.Encoding.Default.GetBytes(file_name);
            byte[] MessageBodyByte_for_download = new byte[30 + Msgbody.Length];
            MessageBodyByte_for_download = BuildDataPackage_For_Pull(Msgbody, 0x21, Device_id);
            SendMessage(MessageBodyByte_for_download);

            ReceiveMessage_For_download(download_type);
    
            return "Download success";
        }

        public string Rename_File(string file_name_origin,string file_name)
        {
            FileInfo fi = new FileInfo(PATH + file_name);
            var di = fi.Directory;
            if (!di.Exists)
                di.Create();

            File.Move(PATH + file_name_origin, PATH + file_name);

            return "rename successfully";
        }

        public string Backup_File(string file_name)
        {

            if (File.Exists(PATH + file_name))
            // overwrite = true
            {
                System.DateTime currentTime = new System.DateTime();
                string newPath = (Backup_PATH + file_name).Replace(".", currentTime.ToString("s")+".");
                File.Move(PATH + file_name, Backup_PATH + file_name + ".bak");
            }
            return "backup successfully";
        }

        public string Delete_File(string file_name)
        {
            
            if (File.Exists(PATH + file_name))
                File.Delete(PATH + file_name);
            
            return "delete successfully";
        }
        public string Upload_File(string file_name,int upload_type)
        {
            try
            {
                // file_list part 
                // need to add

                // Packet_type = 0x20
                byte Packet_type = 0x20;
                byte Device_id = 0x02;
                byte[] file_data;

                // 0xffff - 300 - 2 
                int MessageBody_Length_Max = 65233;
                   
                //read file 
                //for test
                FileStream file = new FileStream(PATH + file_name, FileMode.Open);
                file.Seek(0, SeekOrigin.Begin);
                long file_length = file.Length;

                file_data = null;

                // MessageBody Length
                int MessageBody_Length = 0;

                if (file_length <= MessageBody_Length_Max)
                {

                    MessageBody_Length = (int)file_length;

                    // 8 + 2 + 1 + 1 + 300 + len(MessageBody) + 16 +2
                    byte[] MessageBodyByte = new byte[MessageBody_Length + 30 + 300];

                    // Packet Initiator[ 8 bytes] 0x11 0xff 0x6c 0x6f 0x6e 0x64 0x6f 0x6e
                    byte[] Initiator = { 0x11, 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e };
                    Buffer.BlockCopy(Initiator, 0, MessageBodyByte, 0, Initiator.Length);

                    // Packet Content Length [ 2 bytes ]
                    // Context_Length = 2 + 300 <file_index_json> + len(MessageBody) 
                    int Context_Length = (2 + 300 + MessageBody_Length);
                    byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
                    Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

                    // Packet Type [ 1 byte ]
                    MessageBodyByte[10] = Packet_type;

                    // Packet ID [ 1 byte ]
                    MessageBodyByte[11] = Device_id;

                    byte[] file_json = new byte[300];

                    file_index_json f_js = new file_index_json() { Name = file_name, Num = 0x1, Index = 0 };
                    string f_js_str = JsonConvert.SerializeObject(f_js);
                    byte[] f_js_bytes = System.Text.Encoding.Default.GetBytes(f_js_str);
                    Buffer.BlockCopy(f_js_bytes, 0, file_json, 0, f_js_str.Length);

                    /** for 1.txt file (cmd = 1)
                        * demo
                        * {
                            "Name": "1.txt",
                            "Num": 2,       
                            "Index": 0
                        }
                        {
                            "Name": "1.txt",
                            "Num": 2,
                            "Index": 1
                        }
                        * 
                        */

                    //copy json data
                    Buffer.BlockCopy(file_json, 0, MessageBodyByte, 12, file_json.Length);

                    switch (upload_type) {

                        case 0:
                            //read file data
                            file_data = new byte[MessageBody_Length];
                            file.Read(file_data, 0, (int)MessageBody_Length);
                            //copy file data  from index = 312
                            Buffer.BlockCopy(file_data, 0, MessageBodyByte, 12 + 300, MessageBody_Length);
                            break;
                        case 1: 
                            //read file data -- rsync
                            file_data = new byte[MessageBody_Length];
                            file.Read(file_data, 0, (int)MessageBody_Length);
                            //copy file data  from index = 312
                            Buffer.BlockCopy(file_data, 0, MessageBodyByte, 12 + 300, MessageBody_Length);
                            break;
                        }

                    // Check_sum
                    byte[] all_data = new byte[file_data.Length + file_json.Length];
                    Buffer.BlockCopy(file_json, 0, all_data, 0, file_json.Length);
                    Buffer.BlockCopy(file_data, 0, all_data, file_json.Length, file_data.Length);
                    byte[] md5 = GetCheck_sum(Packet_type, Device_id, all_data, true);

                    Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + 300 + MessageBody_Length, md5.Length);

                    // end with 0xff 0xee
                    MessageBodyByte[MessageBody_Length + 28 + 300] = 0xff;
                    MessageBodyByte[MessageBody_Length + 29 + 300] = 0xee;

                    SendMessage(MessageBodyByte);

                    byte[] res_flag = ReceiveMessage();
                    // result == "ok"
                    if (res_flag[0] == 0x6f && res_flag[1] == 0x6b)
                    {
                        file.Close();
                    }

                    else
                    {
                        file.Close();
                        return "failed upload" + file_name;
                    }

                }

                else
                {
                    // Rounds 
                    double Packet_num = Math.Ceiling((double)file_length / (double)MessageBody_Length_Max);

                    // Json - Packet Num
                    byte Packet_Num = (byte)Packet_num;
                    // Json - Packet Index
                    byte index = 0;
                    file_data = null;

                    byte[] MessageBodyByte;
                    while (Packet_num > 0)
                    {
                        if (Packet_num != 1)
                            MessageBody_Length = MessageBody_Length_Max;
                        else
                            MessageBody_Length = (int)(file_length - MessageBody_Length_Max * index);

                        MessageBodyByte = new byte[MessageBody_Length + 30 + 300];
                        // Packet Initiator[ 8 bytes] 0x11 0xff 0x6c 0x6f 0x6e 0x64 0x6f 0x6e
                        byte[] Initiator = { 0x11, 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e };
                        Buffer.BlockCopy(Initiator, 0, MessageBodyByte, 0, Initiator.Length);

                        // Packet Content Length [ 2 bytes ]
                        // Context_Length = 2 + 300 <file_index_json> + len(MessageBody) 
                        int Context_Length = (2 + 300 + MessageBody_Length);
                        byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
                        Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

                        // Packet Type [ 1 byte ]
                        MessageBodyByte[10] = Packet_type;

                        // Packet ID [ 1 byte ]
                        MessageBodyByte[11] = Device_id;

                        byte[] file_json = new byte[300];
                        file_index_json f_js = new file_index_json() { Name = file_name, Num = Packet_Num, Index = index };
                        string f_js_str = JsonConvert.SerializeObject(f_js);
                        byte[] f_js_bytes = System.Text.Encoding.Default.GetBytes(f_js_str);
                        Buffer.BlockCopy(f_js_bytes, 0, file_json, 0, f_js_str.Length);

                        //copy json data
                        Buffer.BlockCopy(file_json, 0, MessageBodyByte, 12, file_json.Length);


                        //read file
                        file_data = new byte[MessageBody_Length];
                        //it will continue reading from last position
                        file.Read(file_data, 0, MessageBody_Length);
                        //copy file data
                        Buffer.BlockCopy(file_data, 0, MessageBodyByte, 12 + 300, MessageBody_Length);

                        // Check_sum
                        /*
                            * ** FILE data
                            */
                        // Check_sum
                        byte[] all_data = new byte[file_data.Length + file_json.Length];
                        Buffer.BlockCopy(file_json, 0, all_data, 0, file_json.Length);
                        Buffer.BlockCopy(file_data, 0, all_data, file_json.Length, file_data.Length);
                        byte[] md5 = GetCheck_sum(Packet_type, Device_id, all_data, true);

                        Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + 300 + MessageBody_Length, md5.Length);

                        // end with 0xff 0xee
                        MessageBodyByte[MessageBody_Length + 28 + 300] = 0xff;
                        MessageBodyByte[MessageBody_Length + 29 + 300] = 0xee;


                        Packet_num--;
                        index++;
                        SendMessage(MessageBodyByte);


                    }

                    

                    byte[] res_flag = ReceiveMessage();
                    if (res_flag[0] == 0x6f && res_flag[1] == 0x6b)
                    {
                        file.Close();
                    }

                    else
                    {
                        file.Close();
                        return "failed upload" + file_name;
                    }

                }

                
                // socketClient.Close();
                return "success upload";
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }
        }

        /// <summary>
        /// BuildDataPackage
        /// </summary>
        /// <param name="MessageBody">data</param>
        /// <param name="Packet_type">packet type</param>
        /// <param name="Device_id">device id</param>
        /// <returns></returns>
        public byte[] BuildDataPackage_For_Pull(byte[] MessageBody, byte Packet_type, byte Device_id)
        {
            try
            {
                // MessageBody Length
                int MessageBody_Length = 0;

                if (MessageBody == null)
                    MessageBody_Length = 0;
                else
                    MessageBody_Length = MessageBody.Length;

                // 8 + 2 + 1 + 1 + len(MessageBody) + 16 +2
                byte[] MessageBodyByte = new byte[MessageBody_Length + 30];

                // Packet Initiator[ 8 bytes] 0x11 0xff 0x6c 0x6f 0x6e 0x64 0x6f 0x6e
                byte[] Initiator = { 0x11, 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e };
                Buffer.BlockCopy(Initiator, 0, MessageBodyByte, 0, Initiator.Length);

                // Packet Content Length [ 2 bytes ]
                // Context_Length = len(MessageBody) + 2
                int Context_Length = (int)(MessageBody_Length + 2);

                byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
                Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

                // Packet Type [ 1 byte ]
                MessageBodyByte[10] = Packet_type;

                // Packet ID [ 1 byte ]
                MessageBodyByte[11] = Device_id;

                // Packet Content [ len(data) bytes ]
                if (MessageBody_Length != 0)
                    Buffer.BlockCopy(MessageBody, 0, MessageBodyByte, 12, MessageBody_Length);

                // Check_sum
                byte[] md5 = GetCheck_sum(Packet_type, Device_id, MessageBody, true);

                Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + MessageBody_Length, md5.Length);
                // end with 0xff 0xee
                MessageBodyByte[MessageBody_Length + 28] = 0xff;
                MessageBodyByte[MessageBody_Length + 29] = 0xee;
                return (MessageBodyByte);
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }
        }

        public byte[] UnpackData(byte[] recvMsg)
        {
            try
            {
                if (recvMsg != null)
                {
                    // get type
                    byte Packet_type = recvMsg[0];
                    int Tcp_body_length = recvMsg.Length;
                    int Context_Length;

                    Context_Length = Tcp_body_length - 1;

                    byte[] Msg = new byte[Context_Length];
                    Buffer.BlockCopy(recvMsg, 1, Msg, 0, Context_Length);
                    return (Msg);
                }
                else
                    return null;
            }
            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }
        }
        /// <summary>
        /// Check the cmd_flag in file_info_json
        /// </summary>
        /// <param name="recvMsg">Msg body</param>
        /// <param name="F_check">return json data including file operation</param>
        public void Check_cmd_flag(byte[] recvMsg)
        {
            try
            {

                string result_str = System.Text.Encoding.Default.GetString(recvMsg);
                int cmd_flag = 0;

                // convert json to JArray
                JArray jArray = (JArray)JsonConvert.DeserializeObject(result_str);
                /*
                 * Response-json
                         [
                            {
                                "Name":"1.txt",
                                "Digest":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1],
                                "Cmd":1,
                                "Ext":""
                            },
                            {
                                "Name":"2.txt",
                                "Digest":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1],
                                "Cmd":1,
                                "Ext":""
                            }
                         ]
                 */
                foreach (JObject ja in jArray)
                {
                    cmd_flag = int.Parse(ja["Cmd"].ToString());
                    // the whole upload
                    switch (cmd_flag)
                    {
                        case 1:
                            // upload the whole file 
                            Upload_File(ja["Name"].ToString(), 0);
                            break;
                        case 2:
                            // download the whole file
                            Upload_File(ja["Name"].ToString(), 0);
                            break;
                        case 3:
                            // rename
                            Rename_File(ja["Name"].ToString(), ja["Ext"].ToString());
                            break;
                        case 4:
                            // differ upload
                            Upload_File(ja["Name"].ToString(), 1);
                            break;
                        case 5:
                            // differ download
                            Upload_File(ja["Name"].ToString(), 1);
                            break;
                        case 6:
                            // delete
                            Delete_File(ja["Name"].ToString());
                            break;
                        case 7:
                            // backup
                            Backup_File(ja["Name"].ToString());
                            break;

                    }
                    
                  
                }

            }

            catch (Exception ex)
            {
                listMessage.Add(ex.ToString());
                throw new Exception(ex.Message);
            }

        }

        private byte[] GetCheck_sum(byte Packet_type, byte Device_id, byte[] Msg, bool Send_Or_Recv)
        {
            try
            {
                // Md5 Checksum [16 bytes]
                // private_key = 'aaaaa'
                // send is true
                // recv is false
                int index = 1;
                if (Send_Or_Recv)
                {
                    index++;
                }
                int Msg_length;
                byte[] private_key = { 0x61, 0x61, 0x61, 0x61, 0x61 };
                if (Msg == null)
                    Msg_length = 0;
                else
                    Msg_length = Msg.Length;
                // for Send part
                // Check_sum = { type + id + Msg + private_key}

                // for Recv part
                // Check_sum = { type + Msg + private_key}
                byte[] Check_sum = new byte[index + Msg_length + private_key.Length];

                Check_sum[0] = Packet_type;

                if (Send_Or_Recv)
                {
                    Check_sum[1] = Device_id;
                }
                if (Msg != null)
                    Buffer.BlockCopy(Msg, 0, Check_sum, index, Msg_length);

                Buffer.BlockCopy(private_key, 0, Check_sum, Msg_length + index, private_key.Length);

                byte[] md5 = new byte[16];
                md5 = lib.HexStrTobyte(GetMD5Hash(Check_sum));

                return (md5);
            }
            catch (Exception ex)
            {
                throw new Exception(ex.Message);
            }
        }

        private static string GetMD5Hash(byte[] bytedata)
        {
            try
            {
                System.Security.Cryptography.MD5 md5 = new System.Security.Cryptography.MD5CryptoServiceProvider();
                byte[] retVal = md5.ComputeHash(bytedata);

                StringBuilder sb = new StringBuilder();
                for (int i = 0; i < retVal.Length; i++)
                {
                    sb.Append(retVal[i].ToString("x2"));
                }
                return sb.ToString();
            }
            catch (Exception ex)
            {
                throw new Exception("GetMD5Hash() fail,error:" + ex.Message);
            }
        }

        // stringhex(md5) => byte[]



    }
}