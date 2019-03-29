using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Meerkats_Win.File_json_info;

namespace Meerkats_Win.Class
{
    class Func
    {
        // stringhex(md5) => byte[]
        public byte[] HexStrTobyte(string hexString)
        {
            try
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
            catch (Exception ex)
            {
                throw new Exception(ex.Message);
            }
        }


        public string GetMD5HashFromFile(string fileName)
        {
            try
            {
                FileStream file = new FileStream(fileName, FileMode.Open);
                System.Security.Cryptography.MD5 md5 = new System.Security.Cryptography.MD5CryptoServiceProvider();
                byte[] retVal = md5.ComputeHash(file);
                file.Close();

                StringBuilder sb = new StringBuilder();
                for (int i = 0; i < retVal.Length; i++)
                {
                    sb.Append(retVal[i].ToString("x2"));
                }
                return sb.ToString();
            }
            catch (Exception ex)
            {
                throw new Exception("GetMD5HashFromFile() fail,error:" + ex.Message);
            }
        }

        public string GetMD5Hash(byte[] bytedata)
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

        public List<string> Get_file_block_md5(string Path,out string File_hash_str)
        {
            File_hash_str = null;
            List<string> md5_list =new List<string>();
            FileStream file = new FileStream(Path, FileMode.Open);
            file.Seek(0, SeekOrigin.Begin);
            long file_length = file.Length;
            long block_size = 1024;
            long round = file_length / block_size;
            for (int index = 0; index < round; index++)
            {
                byte[] file_data = new byte[block_size];
                file.Read(file_data, 0, file_data.Length);
                string md5 = GetMD5Hash(file_data);
                File_hash_str += md5;
                md5_list.Add(md5);
            }
            file.Flush();
            file.Close();
            return md5_list;
        }

        public void Search_block_index(string Path,List<string> md5_list, out differ_info_json_list diff_json, out byte[] filedata,string filename)
        {
            diff_json = new differ_info_json_list();
            filedata = null;
            FileStream file = new FileStream(Path, FileMode.Open);
            //differ_info_json_list diff_json = new differ_info_json_list();

            diff_json.List = new List<Differ_info_json>();
            diff_json.Name = filename;
            diff_json.Num = 1;
            diff_json.Idx = 0;

            // long block_size = 1024;
            int current_index = 0;
            int tmp_start = 0;
            bool flag = true;
            bool flag_ex = false;
            int matched = 0;
            byte[] data = new byte[file.Length];
            file.Read(data, 0, data.Length);
            
            while(current_index + 1023 <= data.Length - 1)
            {
                byte[] temp = new byte[1024];

                Buffer.BlockCopy(data, current_index, temp, 0, 1024);

                string md5 = GetMD5Hash(temp);

                for (int n = 0; n < md5_list.Count; n++)
                {
                    if (md5_list[n] == md5)
                    {
                        flag_ex = true;
                        matched = n;
                        break;
                    }


                }

                if (flag_ex)
                {
                    if (!flag)
                    {
                        flag = true;

                        //
                        diff_json.List.Add(new Differ_info_json()
                        {
                            Idx = tmp_start,
                            Typ = 1,
                            Len = current_index - tmp_start
                        });
                        byte[] tmp = new byte[current_index - tmp_start];
                        Buffer.BlockCopy(data, tmp_start, tmp, 0, tmp.Length);
                        if (filedata != null)
                            filedata = filedata.Concat(tmp).Where(a => '1' == '1').ToArray();
                        else
                            filedata = tmp;


                    }
                    diff_json.List.Add(new Differ_info_json()
                    {
                        Idx = matched,
                        Typ = 2,
                        Len = 1024

                    });

                    current_index += 1024;
                    tmp_start = current_index;
                }
                else
                {
                    if (flag)
                    {
                        tmp_start = current_index;
                        flag = false;
                    }
                    current_index++;
                }



            }
            if(current_index + 1023 > data.Length -1)
            {
                diff_json.List.Add(new Differ_info_json()
                {
                    Idx = tmp_start,
                    Typ = 1,
                    Len = data.Length-tmp_start

                });
                byte[] tmp = new byte[data.Length - tmp_start];
                Buffer.BlockCopy(data, tmp_start, tmp, 0, tmp.Length);
                if (filedata != null)
                    filedata = filedata.Concat(tmp).Where(a => '1' == '1').ToArray();
                else
                    filedata = tmp;
            }



        }

        public void Differ_modifer_file(string Path, differ_info_json_list diff_json, byte[] filedata)
        {
            string Path_new = Path + "_new";

            long block_size = 1024;
            FileStream file_old = new FileStream(Path, FileMode.Open);
            FileStream file_new = new FileStream(Path_new, FileMode.Create, FileAccess.Write);
            long last_index = 0;
            for (int i = 0; i < diff_json.List.Count; i++)
            {
                long len = 0;
                switch (diff_json.List[i].Typ)
                {
                    
                    case 1:

                        len = diff_json.List[i].Len;
                        byte[] temp_data_1 = new byte[len];
                        Buffer.BlockCopy(filedata, (int)last_index, temp_data_1, 0, (int)len);
                        file_new.Write(temp_data_1, 0, temp_data_1.Length);
                        last_index += len;
                        break;

                    case 2:
                        len = diff_json.List[i].Len;
                        long Idx = diff_json.List[i].Idx;
                        byte[] temp_data_2 = new byte[len];
                        // move pointer to Idx * block_size
                        file_old.Seek(Idx * block_size, SeekOrigin.Begin);
                        file_old.Read(temp_data_2, 0, (int)block_size);
                        file_new.Write(temp_data_2, 0, temp_data_2.Length);
                        break;
                }
            }

            file_old.Flush();
            file_new.Flush();
            file_old.Close();
            file_new.Close();

            File.Delete(Path);
            File.Move(Path_new,Path);

        }
    }
}
