package com.example.meerkats;

public class FileTypeTCP {

    public String name;

    public String ext;

    public int type;

    public FileTypeTCP(String name, String ext, int type) {
        this.name = name;
        this.ext = ext;
        this.type = type;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getExt() {
        return ext;
    }

    public void setExt(String ext) {
        this.ext = ext;
    }

    public int getType() {
        return type;
    }

    public void setType(int type) {
        this.type = type;
    }
}


