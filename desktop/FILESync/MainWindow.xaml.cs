using FILESync.Class;
using System;
using System.Collections.ObjectModel;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Windows;
using System.Windows.Controls;

namespace FILESync
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {

        public ObservableCollection<FileInfo_cs> List = new ObservableCollection<FileInfo_cs>();
        

        string server_ip = "127.0.0.1";
        int server_port = 8888;
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
            TcpClient tcpClient = new TcpClient();
            try
            {
                tcpClient.Connect(IPAddress.Parse(server_ip), server_port);
                NetworkStream ntwStream = tcpClient.GetStream();
                if (ntwStream.CanWrite)
                {
                    string data=null;
                    Byte[] bytSend = Encoding.UTF8.GetBytes(data);
                    ntwStream.Write(bytSend, 0, bytSend.Length);
                }
                else
                {
                    MessageBox.Show("Error");

                    ntwStream.Close();
                    tcpClient.Close();

                    return;
                }


            }
            catch (Exception ex)
            {
                MessageBox.Show("Error:" + ex.Message);
            }

        }
    }
}
