using System;
using System.IO;
using System.Windows;
using System.Windows.Media.Imaging;

namespace FILESync.Class
{
    /// <summary>
    /// 文件夹列表item，继承自ImagedTreeViewItem
    /// </summary>
    class DirectoryTreeViewItem : ImagedTreeViewItem
    {
        DirectoryInfo dir;

        //Constructor requires DirectoryInfo object
        public DirectoryTreeViewItem(DirectoryInfo pDir)
        {
            this.dir = pDir;
            Text = pDir.Name;

            SelectedImage = new BitmapImage(new Uri("pack://application:,,/img/OPEN.BMP"));
            UnselectedImage = new BitmapImage(new Uri("pack://application:,,/img/CLOSED.BMP"));
        }

        /// <summary>
        /// public property to obtain DirectoryInfo
        /// </summary>
        public DirectoryInfo DirInfo
        {
            get { return dir; }
        }

        /// <summary>
        /// public mathod to populate wtih items
        /// </summary>
        public void Populate()
        {
            DirectoryInfo[] dirs;

            try
            {
                dirs = dir.GetDirectories();
            }
            catch
            {
                return;
            }

            foreach (DirectoryInfo dirChild in dirs)
            {
                Items.Add(new DirectoryTreeViewItem(dirChild));
            }
        }

        /// <summary>
        /// event override to populate subitem
        /// </summary>
        /// <param name="e"></param>
        protected override void OnExpanded(RoutedEventArgs e)
        {
            base.OnExpanded(e);

            foreach (object obj in Items)
            {
                DirectoryTreeViewItem item = obj as DirectoryTreeViewItem;
                item.Populate();
            }
        }
    }
}

