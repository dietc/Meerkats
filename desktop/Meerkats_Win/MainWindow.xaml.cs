using Meerkats_Win.Class;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using static Meerkats_Win.File_json_info;

namespace Meerkats_Win
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        public ObservableCollection<FileInfo_cs> List = new ObservableCollection<FileInfo_cs>();

        // delegate
        public delegate string FuncHandle();
        FuncHandle fh;

        /// <summary>
        /// desktop client id = 0x2
        /// </summary>
        private byte Device_id = 0x2;

        public MainWindow()
        {
            InitializeComponent();
            init_status();
        }

        public void init_status()
        {

            file_tree.Items.Clear();
            file_tree.ItemsSource = List;

            //// Json Filedata
            //string File_json = "{\"1\":\"2\",\"3\":{\"4\":\"5\",\"6\":\"7\"}}";

            //var js_obj = JObject.Parse(File_json);
            ////创建TreeView的数据源
            //file_tree.ItemsSource = js_obj.Children().Select(c => JsonHeaderLogic.FromJToken(c));

            DirectoryTreeView mainTree = new DirectoryTreeView();
            mainTree.SelectedItemChanged += MainTree_SelectedItemChanged;
            file_tree_grid.Children.Add(mainTree);

            //右键菜单
            ContextMenu myContext = new ContextMenu();

            MenuItem myMUItem = new MenuItem();
            myMUItem.Header = "Open";
            myMUItem.Name = "Menu01";
            myContext.Items.Add(myMUItem);

            myMUItem = new MenuItem();
            myMUItem.Header = "View";
            myMUItem.Name = "Menu02";
            //myMUItem.Click += FileLook_Click;
            myContext.Items.Add(myMUItem);

            myMUItem = new MenuItem();
            myMUItem.Header = "Refresh";
            myMUItem.Name = "Menu03";
            myContext.Items.Add(myMUItem);

            myMUItem = new MenuItem();
            myMUItem.Header = "Rename";
            myMUItem.Name = "Menu04";
            myContext.Items.Add(myMUItem);

            myMUItem = new MenuItem();
            myMUItem.Header = "Delete";
            myMUItem.Name = "Menu05";
            myContext.Items.Add(myMUItem);


            myMUItem = new MenuItem();
            myMUItem.Header = "New directory";
            myMUItem.Name = "Menu06";
            myContext.Items.Add(myMUItem);

            //myMUItem = new MenuItem();
            //myMUItem.Header = "Upload file";
            //myMUItem.Click += upload_file;
            //myMUItem.Name = "Menu07";
            //myContext.Items.Add(myMUItem);

            file_info.ContextMenu = myContext;

        }

        private void upload_file(object sender, RoutedEventArgs e)
        {
            //
        }

        /// <summary>
        /// 文件夹树改变时，查找文件夹下是否存在文件，如果存在，则显示
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        private void MainTree_SelectedItemChanged(object sender, RoutedPropertyChangedEventArgs<object> e)
        {
            DirectoryTreeViewItem item = e.NewValue as DirectoryTreeViewItem;

            //stack.Children.Clear();
            file_info.Items.Clear();

            FileInfo[] fileInfos;

            try
            {
                fileInfos = item.DirInfo.GetFiles();
            }
            catch
            {
                return;
            }

            foreach (FileInfo info in fileInfos)
            {
                FileInfo_cs myFile = new FileInfo_cs();
                myFile.strFileName = info.Name;
                myFile.strFileType = info.Extension;
                myFile.strFileSize = info.Length.ToString();
                myFile.strlastModifyTime = info.LastAccessTime.ToString();
                file_info.Items.Add(myFile);
            }
        }


        // connect to the server or refresh
        private void Conect_btn_Click(object sender, RoutedEventArgs e)
        {

            fh = new FuncHandle(this.send_data_rev_data);
            AsyncCallback callback = new AsyncCallback(this.AsyncCallbackImpl);
            fh.BeginInvoke(callback, null);
            //send_data_rev_data();
        }

        private string send_data_rev_data()
        {

            SocketTCPClient t1 = new SocketTCPClient();

            t1.CreateInstance();
           
            //
            t1.SendMessage(Get_local_File_info("F:\\fortest\\"));

            byte[] result = t1.ReceiveMessage();

            // get file cmd_flag < 6 operation >
            file_check f_check = new file_check();
            t1.Check_cmd_flag(result, f_check);

            
            string status_file = null;
            status_file = t1.Upload_File(f_check.upload);

            // string result_str = System.Text.Encoding.Default.GetString(result);
            // fortest.Text = result_str;


            // test for download
            byte[] test_for_download = null;
            byte[] MessageBodyByte_for_download = new byte[30];
            MessageBodyByte_for_download = t1.BuildDataPackage_For_Pull(test_for_download, 0x21, Device_id);
            t1.SendMessage(MessageBodyByte_for_download);

            string str_status_1 = t1.ReceiveMessage_For_download(0);
            string str_status_2 = t1.ReceiveMessage_For_download(0);

            t1.DisconnectServer();
            return status_file + " + " + str_status_1 + " + " + str_status_2;

        }

        public void AsyncCallbackImpl(IAsyncResult ar)
        {
            string re = fh.EndInvoke(ar);
            MessageBox.Show(re + ar.AsyncState);
        }

        public byte[] Get_local_File_info(string PATH)
        {
            DirectoryInfo dir = new DirectoryInfo(PATH);

            FileSystemInfo[] fsinfos = dir.GetFileSystemInfos();
            SocketTCPClient t1 = new SocketTCPClient();

            List<file_info_json> file_json = new List<file_info_json>();

            /**
             * json file_info demo
             *
             *  [
             *      {
             *          "Name":"1.txt",
             *          "Typ":1,
             *          "Digest":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]
             *      },
             *      {
             *          "Name":"2.txt",
             *          "Typ":1,
             *          "Digest":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]
             *      }
             *  ]
             *  
            **/
            foreach (FileSystemInfo file in fsinfos)
            {
                byte[] file_md5 = t1.HexStrTobyte(t1.GetMD5HashFromFile(file.FullName));
                List<byte> file__md5_list = new List<byte>();
                for (int i = 0; i < 16; i++)
                    file__md5_list.Add(file_md5[i]);

                file_json.Add(new file_info_json() {
                    Name = file.Name,
                    Typ = 0x1,
                    Digest = file__md5_list
                });
             
            }
            // Json serialize
            string MsgBody_str = JsonConvert.SerializeObject(file_json);

            byte[] MsgBody = System.Text.Encoding.Default.GetBytes(MsgBody_str);
            // encapsulate packets
            byte[] MessageBodyByte = new byte[MsgBody.Length + 30];
            MessageBodyByte = t1.BuildDataPackage_For_Pull(MsgBody, 0x2, Device_id);

            return MessageBodyByte;
        }

        private static void listDirectory(string path, int leval, List<file_info_json> file_json)
        {
            DirectoryInfo theFolder = new DirectoryInfo(@path);

            leval++;

            //遍历文件
            foreach (FileInfo NextFile in theFolder.GetFiles())
            {
                SocketTCPClient t1 = new SocketTCPClient();
                byte[] file_md5 = t1.HexStrTobyte(t1.GetMD5HashFromFile(NextFile.FullName));
                List<byte> file__md5_list = new List<byte>();
                for (int i = 0; i < 16; i++)
                    file__md5_list.Add(file_md5[i]);

                file_json.Add(new file_info_json()
                {
                    Name = NextFile.Name,
                    Typ = 0x1,
                    Digest = file__md5_list
                });
            }

            //遍历文件夹
            foreach (DirectoryInfo NextFolder in theFolder.GetDirectories())
            {
                listDirectory(NextFolder.FullName, leval, file_json);
            }
        }




        /**
         * // string => byte[]：

              byte[] byteArray = System.Text.Encoding.Default.GetBytes ( str );

           //byte[] => string：
              string str = System.Text.Encoding.Default.GetString ( byteArray );
         *
         */
    }
}
