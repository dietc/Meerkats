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
            fortest.IsReadOnly = true;
            fortest.Text = string.Empty;


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
            send_data_rev_data();
            
        }

        private void send_data_rev_data()
        {
            SocketTCPClient t1 = new SocketTCPClient();

            string testdata = null;
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
            // for text md5 = [1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]
            

            // convert to json list
            // Typ = 0x1 => file
            // Typ = 0x2 => Directory

            byte[] md11 = t1.HexStrTobyte(t1.GetMD5HashFromFile("F:\\fortest\\1.txt"));
            byte[] md22 = t1.HexStrTobyte(t1.GetMD5HashFromFile("F:\\fortest\\2.txt"));

            List<byte> file_1_md5 = new List<byte>();
            for (int i = 0; i < 16; i++)
                file_1_md5.Add(md11[i]);

            List<byte> file_2_md5 = new List<byte>();
            for (int i = 0; i < 16; i++)
                file_1_md5.Add(md22[i]);


            List <file_info_json> testdemo = new List<file_info_json>()
            {
                new file_info_json()
                {
                    Name = "1.txt",
                    Typ = 0x1,
                    Digest = file_1_md5
                },

                new file_info_json()
                {
                    Name = "2.txt",
                    Typ = 0x1,
                    Digest = file_2_md5
                }

            };

            // Json serialize
            testdata = JsonConvert.SerializeObject(testdemo);

            t1.CreateInstance();

            //byte[] test = t1.HexStrTobyte(testdata);

            byte[] test = System.Text.Encoding.Default.GetBytes(testdata);
            byte[] MessageBodyByte = new byte[test.Length + 30];

            MessageBodyByte = t1.BuildDataPackage_For_Pull(test, 0x2, Device_id);

            
            t1.SendMessage(MessageBodyByte);
            byte[] result = t1.ReceiveMessage();

            List<string> file_name = t1.Check_If_Upload(result);

            

            if (file_name.Count != 0)
            {
                fortest.Text = t1.Upload_File(file_name);
            }

            // string result_str = System.Text.Encoding.Default.GetString(result);
            // fortest.Text = result_str;


            // test for download
            byte[] test_for_download = null;
            byte[] MessageBodyByte_for_download = new byte[30];
            MessageBodyByte_for_download = t1.BuildDataPackage_For_Pull(test_for_download, 0x21, Device_id);
            t1.SendMessage(MessageBodyByte_for_download);

            string str_status_1 = t1.ReceiveMessage_For_download();
            string str_status_2 = t1.ReceiveMessage_For_download();

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
