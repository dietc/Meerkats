package com.example.meerkats.bean;

public class FileBean {
    private String fName;
    private String fPath;
    private FileType fType;
    private long fSize;

    public String getfName() {
        return fName;
    }

    public void setfName(String fName) {
        this.fName = fName;
    }

    public String getfPath() {
        return fPath;
    }

    public void setfPath(String fPath) {
        this.fPath = fPath;
    }

    public FileType getfType() {
        return fType;
    }

    public void setfType(FileType fType) {
        this.fType = fType;
    }

    public long getfSize() {
        return fSize;
    }

    public void setfSize(long fSize) {
        this.fSize = fSize;
    }
}
