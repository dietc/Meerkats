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
    }
}