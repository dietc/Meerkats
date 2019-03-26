package com.example.meerkats;

public class Rename {
    public String name;
    public String ext;

    public Rename(String name, String ext) {
        this.name = name;
        this.ext = ext;
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

}

