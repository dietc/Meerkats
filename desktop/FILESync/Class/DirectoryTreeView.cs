using System.IO;
using System.Windows.Controls;

namespace FILESync.Class
{
    class DirectoryTreeView : TreeView
    {
        /// <summary>
        /// Constructor builds
        /// </summary>
        public DirectoryTreeView()
        {
            RefreshTree();
        }

        public void RefreshTree()
        {
            BeginInit();
            Items.Clear();

            //Obtain the disk drivers
            DriveInfo[] drivers = DriveInfo.GetDrives();

            foreach (DriveInfo drive in drivers)
            {
                
                char chDrive = drive.Name.ToUpper()[0];
                if (chDrive == 'F')
                {
                    DirectoryTreeViewItem item = new DirectoryTreeViewItem(drive.RootDirectory);

                    //display ...
                    if (chDrive != 'A' && chDrive != 'B' && drive.IsReady && drive.VolumeLabel.Length > 0)
                    {
                        item.Text = string.Format("{0}  ({1})", drive.VolumeLabel, drive.Name);
                    }
                    else
                    {
                        item.Text = string.Format("{0}  ({1})", drive.DriveType, drive.Name);
                    }

                    Items.Add(item);

                    if (chDrive != 'A' && chDrive != 'B' && drive.IsReady)
                    {
                        item.Populate();
                    }
                }
                else
                    continue;
            }

            EndInit();
        }
    }
}

