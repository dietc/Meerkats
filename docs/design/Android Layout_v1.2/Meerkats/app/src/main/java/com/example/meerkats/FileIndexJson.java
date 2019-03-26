package com.example.meerkats;


public class FileIndexJson {

    public FileIndexJson(String name, Byte num, Byte index) {
        Name = name;
        Num = num;
        Index = index;
    }

    String Name;
    Byte Num;
    Byte Index;

    String getName() {
        return Name;
    }

    public void setName(String name) {
        Name = name;
    }

    public Byte getNum() {
        return Num;
    }

    public void setNum(Byte num) {
        Num = num;
    }

    public Byte getIndex() {
        return Index;
    }

    public void setIndex(Byte index) {
        Index = index;
    }
}



