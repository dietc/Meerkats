using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Meerkats_Win
{
    class File_json_info
    {
        public class file_index_json
        {
            public string Name { get; set; }
            public byte Num { get; set; }
            public byte Index { get; set; }

        }

        public class file_info_json
        {
            public string Name { get; set; }
            public byte Typ { get; set; }
            public List<byte> Digest { get; set; }

        }

        public class file_check
        {
            public List<File_type> upload { get; set; }
            public List<File_type> download { get; set; }
            public List<Rename> rename { get; set; }
            public List<Delete> delete { get; set; }
            public List<Backup> backup { get; set; }

        }
        
        public class File_type
        {
            public string Name { get; set; }

            public string Ext { get; set; }
            /// <summary>
            /// Type = 0x0 -> the whole upload/download
            /// Type = 0x1 -> differ upload/download
            /// </summary>
            public int Type { get; set; }

        }

        

        public class Rename
        {
            public string Name { get; set; }
            public string Ext { get; set; }
        }

        public class Delete
        {
            public string Name { get; set; }
        }

        public class Backup
        {
            public string Name { get; set; }
        }



    }
}