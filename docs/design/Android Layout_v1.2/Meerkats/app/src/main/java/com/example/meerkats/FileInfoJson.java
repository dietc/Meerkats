package com.example.meerkats;

import java.util.ArrayList;


public class FileInfoJson {
    String Name;
    Byte Type;


    public String getName() {
        return Name;
    }

    public void setName(String name) {
        Name = name;
    }

    public Byte getType() {
        return Type;
    }

    public void setType(Byte type) {
        Type = type;
    }

    public ArrayList<Byte> getDigest() {
        return Digest;
    }

    public void setDigest(ArrayList<Byte> digest) {
        Digest = digest;
    }

    public ArrayList<Byte> Digest;


}
