package com.example.meerkats;

import java.util.ArrayList;
import java.util.List;


public class FileInfoJson{
    String Name;
    byte Type;
    List<Byte> Digest;

    public FileInfoJson(String name, byte type, List<Byte> digest) {
        Name = name;
        Type = type;
        Digest = digest;
    }

    public String getName() {
        return Name;
    }

    public void setName(String name) {
        Name = name;
    }

    public byte getType() {
        return Type;
    }

    public void setType(byte type) {
        Type = type;
    }



    public List<Byte> getDigest() {
        return Digest;
    }

    public void setDigest(List<Byte> digest) {
        Digest = digest;
    }
}
