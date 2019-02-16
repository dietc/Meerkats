using System.Collections.ObjectModel; //ObservableCollection
using System.ComponentModel; //INotifyPropertyChanged

namespace FILESync.Class
{
    /// <summary>
    /// 自定义的文件信息类，用于绑定到ListViewItem中
    /// </summary>
    /// 
    public class FileInfo_cs : INotifyPropertyChanged //通知接口
    {
        public event PropertyChangedEventHandler PropertyChanged;
        string _strFileName;
        string _strFileType;
        string _strFileSize;
        string _strlastModifyTime;


        /// <summary>
        /// 文件名
        /// </summary>
        public string strFileName
        {
            get { return _strFileName; }
            set {
                _strFileName = value;
                PropertyChanged?.Invoke(this, new PropertyChangedEventArgs("strFileName"));
            }
        }

        /// <summary>
        /// 文件类型
        /// </summary>
        public string strFileType
        {
            get { return _strFileType; }
            set {
                _strFileType = value;
                PropertyChanged?.Invoke(this, new PropertyChangedEventArgs("strFileType"));
            }
        }

        /// <summary>
        /// 文件大小
        /// </summary>
        public string strFileSize
        {
            get { return _strFileSize; }
            set {
                _strFileSize = value;
                PropertyChanged?.Invoke(this, new PropertyChangedEventArgs("strFileSize"));
            }
        }

        /// <summary>
        /// 最后一次修改时间
        /// </summary>
        public string strlastModifyTime
        {
            get { return _strlastModifyTime; }
            set {
                _strlastModifyTime = value;
                PropertyChanged?.Invoke(this, new PropertyChangedEventArgs("strlastModifyTime"));
            }
        }

       

        //public FileInfo_cs(string _strFileName, string _strFileType, string _strFileSize, string _strlastModifyTime)
        //{
        //    _strFileName = strFileName;
        //    _strFileType = strFileType;
        //    _strFileSize = strFileSize;
        //    _strlastModifyTime = strlastModifyTime;
            
        //}
    }

 

}
