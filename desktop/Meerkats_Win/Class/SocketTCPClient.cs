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

namespace Meerkats_Win.Class
{

    class SocketTCPClient
    {
        private static string ip = "178.128.45.7";
        private static int port = 4356;
        private static Socket socketClient;

        // the path for stored data 
        private static string PATH = "E:\\fortest\\";

        // public static List<string> listMessage = new List<string>();

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

            socketClient.Connect(IPAddress.Parse(ip), port);

            //Thread threadConnect = new Thread(new ThreadStart(ReceiveMessage));
            //threadConnect.Start();

        }
        /// <summary>
        /// receive Msg
        /// </summary>
        public byte[] ReceiveMessage()
        {

            int Tcp_header_length = 10;
            int Tcp_body_length = 0;

            byte[] Tcp_header = new byte[Tcp_header_length];

            socketClient.Receive(Tcp_header, 0, Tcp_header_length, 0);
            Tcp_body_length = (Tcp_header[8] * 0xff) + (Tcp_header[9]);

            byte[] recvBytes = new byte[Tcp_body_length];

            socketClient.Receive(recvBytes, 0, Tcp_body_length, 0);

            byte[] md5 = new byte[16];
            socketClient.Receive(md5, 0, 16, 0);

            byte[] End_flag = new byte[2];
            socketClient.Receive(End_flag, 0, 2, 0);
            // check End_flag
            if (End_flag[0] == 0xff && End_flag[1] == 0xee)
            {

                List<string> file_name = Check_Upload(UnpackData_For_Pull(recvBytes));
                if (file_name.Count == 0)
                {
                    socketClient.Close();
                    return null;
                }
                else
                {
                    // Upload_file part

                    Upload_File(file_name);


                    socketClient.Close();
                    return null;
                }
                //socketClient.Close();
                //return (UnpackData(recvBytes));


            }

            else
            {
                socketClient.Close();
                return null;
            }

        }
        ///<summary>
        ///Send Msg
        ///</summary>
        public void SendMessage(byte[] sendBytes)
        {
            //check if connected
            if (socketClient.Connected)
            {
                // get RemoteEndPoint info <ip,port>
                IPEndPoint ipe = (IPEndPoint)socketClient.RemoteEndPoint;
                socketClient.Send(sendBytes, sendBytes.Length, 0);
            }
        }

        public void Upload_File(List<string> file_name)
        {
            // file_list part 
            // need to add

            // Packet_type = 0x20
            byte Packet_type = 0x20;
            byte Device_id = 0x02;
            byte[] file_data;

            // 0xffff - 300 - 2 
            int MessageBody_Length_Max = 3;

            for (int i = 0; i < file_name.Count; i++)
            {
                //read file 
                //for test
                FileStream file = new FileStream(PATH + file_name[i], FileMode.Open);
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
                    // Context_Length = len(MessageBody) + 2
                    byte Context_Length = (byte)(MessageBody_Length + 2);
                    byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
                    Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

                    // Packet Type [ 1 byte ]
                    MessageBodyByte[10] = Packet_type;

                    // Packet ID [ 1 byte ]
                    MessageBodyByte[11] = Device_id;

                    byte[] file_json = new byte[300];

                    file_index_json f_js = new file_index_json() { Name = file_name[i], Num = 0x1, Index = 0 };
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

                    //read file data
                    file_data = new byte[MessageBody_Length];
                    file.Read(file_data, 0, (int)MessageBody_Length);
                    //copy file data  from index = 312
                    Buffer.BlockCopy(file_data, 0, MessageBodyByte, 12 + 300, MessageBody_Length);

                    // Check_sum
                    byte[] md5 = GetCheck_sum(Packet_type, Device_id, file_data, true);

                    Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + 300 + MessageBody_Length, md5.Length);

                    // end with 0xff 0xee
                    MessageBodyByte[MessageBody_Length + 28 + 300] = 0xff;
                    MessageBodyByte[MessageBody_Length + 29 + 300] = 0xee;

                    // SendMessage(MessageBodyByte);

                }

                else
                {

                    double Packet_num = Math.Ceiling((double)file_length / (double)MessageBody_Length_Max);
                    byte Packet_Num = (byte)Packet_num;
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
                        // Context_Length = len(MessageBody) + 2
                        byte Context_Length = (byte)(MessageBody_Length + 2);
                        byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
                        Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

                        // Packet Type [ 1 byte ]
                        MessageBodyByte[10] = Packet_type;

                        // Packet ID [ 1 byte ]
                        MessageBodyByte[11] = Device_id;

                        byte[] file_json = new byte[300];
                        file_index_json f_js = new file_index_json() { Name = file_name[i], Num = Packet_Num, Index = index };
                        string f_js_str = JsonConvert.SerializeObject(f_js);
                        byte[] f_js_bytes = System.Text.Encoding.Default.GetBytes(f_js_str);
                        Buffer.BlockCopy(f_js_bytes, 0, file_json, 0, f_js_str.Length);

                        //copy json data
                        Buffer.BlockCopy(file_json, 0, MessageBodyByte, 12, file_json.Length);


                        //read file
                        file_data = new byte[MessageBody_Length];
                        file.Read(file_data, 0, MessageBody_Length);
                        //copy file data
                        Buffer.BlockCopy(file_data, 0, MessageBodyByte, 12 + 300, MessageBody_Length);

                        // Check_sum
                        /*
                         * ** FILE data
                         */
                        byte[] md5 = GetCheck_sum(Packet_type, Device_id, file_data.Skip(index * MessageBody_Length_Max).Take(MessageBody_Length).ToArray(), true);

                        Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + 300 + MessageBody_Length, md5.Length);

                        // end with 0xff 0xee
                        MessageBodyByte[MessageBody_Length + 28 + 300] = 0xff;
                        MessageBodyByte[MessageBody_Length + 29 + 300] = 0xee;


                        Packet_num--;
                        index++;

                        // SendMessage(MessageBodyByte);

                    }

                }

                file.Close();

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
            // MessageBody Length
            int MessageBody_Length = 0;
            MessageBody_Length = MessageBody.Length;

            // 8 + 2 + 1 + 1 + len(MessageBody) + 16 +2
            byte[] MessageBodyByte = new byte[MessageBody_Length + 30];

            // Packet Initiator[ 8 bytes] 0x11 0xff 0x6c 0x6f 0x6e 0x64 0x6f 0x6e
            byte[] Initiator = { 0x11, 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e };
            Buffer.BlockCopy(Initiator, 0, MessageBodyByte, 0, Initiator.Length);

            // Packet Content Length [ 2 bytes ]
            // Context_Length = len(MessageBody) + 2
            byte Context_Length = (byte)(MessageBody_Length + 2);
            byte[] length = { (byte)(Context_Length >> 8), (byte)(Context_Length & 0xff) };
            Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

            // Packet Type [ 1 byte ]
            MessageBodyByte[10] = Packet_type;

            // Packet ID [ 1 byte ]
            MessageBodyByte[11] = Device_id;

            // Packet Content [ len(data) bytes ]
            Buffer.BlockCopy(MessageBody, 0, MessageBodyByte, 12, MessageBody_Length);

            // Check_sum
            byte[] md5 = GetCheck_sum(Packet_type, Device_id, MessageBody, true);

            Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + MessageBody.Length, md5.Length);
            // end with 0xff 0xee
            MessageBodyByte[MessageBody_Length + 28] = 0xff;
            MessageBodyByte[MessageBody_Length + 29] = 0xee;
            return (MessageBodyByte);
        }

        public byte[] UnpackData_For_Pull(byte[] recvMsg)
        {
            if (recvMsg.Length != 0)
            {
                int Tcp_body_length = recvMsg.Length;
                int Context_Length = Tcp_body_length - 1;
                int Packet_type = recvMsg[0];

                byte[] Msg = new byte[Context_Length];
                Buffer.BlockCopy(recvMsg, 1, Msg, 0, Context_Length);
                return (Msg);
            }
            else
                return null;
        }

        public List<string> Check_Upload(byte[] recvMsg)
        {
            string result_str = System.Text.Encoding.Default.GetString(recvMsg);
            string cmd_flag = null;
            List<string> file_name = new List<string>();
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
                cmd_flag = ja["Cmd"].ToString();
                if (cmd_flag == "1")
                {
                    file_name.Add(ja["Name"].ToString());
                }
                else
                    continue;
            }

            return file_name;
        }

        private byte[] GetCheck_sum(byte Packet_type, byte Device_id, byte[] Msg, bool Send_Or_Recv)
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

            byte[] private_key = { 0x61, 0x61, 0x61, 0x61, 0x61 };

            // for Send part
            // Check_sum = { type + id + Msg + private_key}

            // for Recv part
            // Check_sum = { type + Msg + private_key}
            byte[] Check_sum = new byte[index + Msg.Length + private_key.Length];

            Check_sum[0] = Packet_type;

            if (Send_Or_Recv)
            {
                Check_sum[1] = Device_id;
            }

            Buffer.BlockCopy(Msg, 0, Check_sum, index, Msg.Length);
            Buffer.BlockCopy(private_key, 0, Check_sum, Msg.Length + index, private_key.Length);

            byte[] md5 = new byte[16];
            md5 = HexStrTobyte(GetMD5Hash(Check_sum));

            return (md5);
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
        public byte[] HexStrTobyte(string hexString)
        {
            hexString = hexString.Replace(" ", "");
            if ((hexString.Length % 2) != 0)
                hexString += " ";
            byte[] returnBytes = new byte[hexString.Length / 2];
            for (int i = 0; i < returnBytes.Length; i++)
                returnBytes[i] = Convert.ToByte(hexString.Substring(i * 2, 2).Trim(), 16);
            return returnBytes;

            ////byte[] => stringhex(md5)

            //byte[] buffer = {};

            //StringBuilder strBuider = new StringBuilder();
            //for (int index = 0; index < count; index++)
            //{
            //    strBuider.Append(((int)buffer[index]).ToString("X2"));
            //}
        }



    }
}
