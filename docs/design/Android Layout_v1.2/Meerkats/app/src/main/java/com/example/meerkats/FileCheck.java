package com.example.meerkats;


import java.util.List;

public class FileCheck {

    public List<FileTypeTCP> upload;
    public List<FileTypeTCP> download;
    public List<Rename> rename;
    public List<Delete> delete;
    public List<Backup> backup;

    public List<FileTypeTCP> getUpload() {
        return upload;
    }

    public void setUpload(List<FileTypeTCP> upload) {
        this.upload = upload;
    }

    public List<FileTypeTCP> getDownload() {
        return download;
    }

    public void setDownload(List<FileTypeTCP> download) {
        this.download = download;
    }

    public List<Rename> getRename() {
        return rename;
    }

    public void setRename(List<Rename> rename) {
        this.rename = rename;
    }

    public List<Delete> getDelete() {
        return delete;
    }

    public void setDelete(List<Delete> delete) {
        this.delete = delete;
    }

    public List<Backup> getBackup() {
        return backup;
    }

    public void setBackup(List<Backup> backup) {
        this.backup = backup;
    }
}

