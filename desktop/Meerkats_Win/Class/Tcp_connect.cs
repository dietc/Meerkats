using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;
using System.Security.Cryptography;

namespace Meerkats_Win.Class
{
    class Tcp_connect
    {
        public string tcp_send(string ip,int port,string data)
        {
            TcpClient tcpClient = new TcpClient();
            try
            {
                tcpClient.Connect(IPAddress.Parse(ip), port);
                NetworkStream ntwStream = tcpClient.GetStream();
                if (ntwStream.CanWrite)
                {
                    uint content_length = (uint)(data.Length + 2);
                    uint total_length = (uint)(content_length + 28);
                    string testdata = null;
                    //initiator  -- 8 
                    byte[] initiator = { 0x11, 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e };
                    testdata += System.Text.Encoding.Default.GetString(initiator);
                    //length -- 2 
                    byte[] length = { (byte)(uint)(content_length >> 8), (byte)(uint)(content_length & 0xff) };
                    testdata += System.Text.Encoding.Default.GetString(length);
                    // type  -- 1
                    byte[] type = {0x1};
                    testdata += System.Text.Encoding.Default.GetString(type);
                    //device id -- 1
                    byte[] device_id = { 0x1 };
                    testdata += System.Text.Encoding.Default.GetString(device_id);
                    //context
                    testdata += data;
                    //checksum
                    string checksum = data + "aaaaa";
                    testdata += GetMd5Str(checksum);
                    //end
                    byte[] end = { 0xff, 0xee };
                    testdata += System.Text.Encoding.Default.GetString(end);

                    Byte[] bytSend = Encoding.UTF8.GetBytes(testdata);
                    ntwStream.Write(bytSend, 0, bytSend.Length);
                }
                else
                {
                    return("false");
                }

                return ("success");


            }
            catch (Exception ex)
            {
                return ("Error:" + ex.Message);
            }
        }

        /// <summary>
        /// MD5(16位加密)
        /// </summary>
        /// <param name="ConvertString">需要加密的字符串</param>
        /// <returns>MD5加密后的字符串</returns>
        public string GetMd5Str(string ConvertString)
        {
            string md5Pwd = string.Empty;

            //使用加密服务提供程序
            MD5CryptoServiceProvider md5 = new MD5CryptoServiceProvider();

            //将指定的字节子数组的每个元素的数值转换为它的等效十六进制字符串表示形式。
            md5Pwd = BitConverter.ToString(md5.ComputeHash(UTF8Encoding.Default.GetBytes(ConvertString)), 4, 8);

            md5Pwd = md5Pwd.Replace("-", "");

            return md5Pwd;
        }

    }
}
