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

namespace Meerkats_Win.Class
{

    class SocketTCPClient
    {
        private static string ip = "178.128.45.7";
        private static int port = 4356;
        private static Socket socketClient;


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
            int block_size = 1024;
            byte[] recvBytes = new byte[block_size];
            int tran_size = 0;

            while(true)
            {
                tran_size = socketClient.Receive(recvBytes, block_size, 0);
                if (tran_size < block_size)
                    break;
                else
                    continue;      
            }
            
            socketClient.Close();
            return (recvBytes);



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

        /// <summary>
        /// BuildDataPackage
        /// </summary>
        /// <param name="MessageBody">data</param>
        /// <param name="Packet_type">packet type</param>
        /// <param name="Device_id">device id</param>
        /// <returns></returns>
        public byte[] BuildDataPackage(byte[] MessageBody, byte Packet_type, byte Device_id)
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
            byte[] length = { (byte)( Context_Length >> 8), (byte)(Context_Length & 0xff) };
            Buffer.BlockCopy(length, 0, MessageBodyByte, 8, length.Length);

            // Packet Type [ 1 byte ]
            MessageBodyByte[10] = Packet_type;

            // Packet ID [ 1 byte ]
            MessageBodyByte[11] = Device_id;

            // Packet Content [ len(data) bytes ]
            Buffer.BlockCopy(MessageBody, 0, MessageBodyByte, 12, MessageBody.Length);

            // Md5 Checksum [16 bytes]
            // private_key = 'aaaaa'
            byte[] private_key = {0x61, 0x61, 0x61, 0x61, 0x61 };
            byte[] Check_sum = new byte[MessageBody_Length + private_key.Length + 2 ];

            // Check_sum = {type + id + MessageBody + private_key}
            Check_sum[0] = Packet_type;
            Check_sum[1] = Device_id;
            Buffer.BlockCopy(MessageBody, 0, Check_sum, 2, MessageBody.Length);
            Buffer.BlockCopy(private_key, 0, Check_sum, MessageBody.Length+2, private_key.Length);

            // MD5
            byte[] md5 = HexStrTobyte(GetMD5Hash(Check_sum));

            Buffer.BlockCopy(md5, 0, MessageBodyByte, 12 + MessageBody.Length, md5.Length);
            // end with 0xff 0xee
            MessageBodyByte[MessageBody_Length + 28] = 0xff;
            MessageBodyByte[MessageBody_Length + 29] = 0xee;
            return (MessageBodyByte);
        }

        public byte[] UnpackData(byte[] recvMsg)
        {
            int MessageBody_len = (int)recvMsg[9] - 1;
            byte[] Msg = new byte[MessageBody_len];
            Buffer.BlockCopy(recvMsg, 10, Msg, 0, MessageBody_len);
            return (Msg);

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
        private static byte[] HexStrTobyte(string hexString)
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
