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

        public class file_info_json_pull
        {
            public string Name { get; set; }
            public byte Typ { get; set; }
            public List<byte> Digest { get; set; }

        }
        public class differ_info_json_list
        {
            public string Name { get; set; }
            public int Num { get; set; }
            public int Idx { get; set; }
            public List<Differ_info_json> List { get; set; }
        }


        public class Differ_info_json
        {
            public byte Typ { get; set; }
            public long Idx { get; set; }
            public long Len { get; set; }
        }

        

    }
}