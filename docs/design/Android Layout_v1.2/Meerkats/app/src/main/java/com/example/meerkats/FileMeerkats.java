package com.example.meerkats;

import com.example.meerkats.bean.FileType;

import java.io.File;

public class FileMeerkats {

    public FileType getFileType(File f) {
        if (f.isDirectory()) {
            return FileType.directory;
        } else {
            String fileName = f.getName();

            if (fileName.endsWith(".png") || fileName.endsWith(".PNG")
                    || fileName.endsWith(".jpg") || fileName.endsWith(".JPG")
                    || fileName.endsWith(".jpeg") || fileName.endsWith(".JPEG")
                    || fileName.endsWith(".gif") || fileName.endsWith(".GIF")) {
                return FileType.image;
            }

            if (fileName.endsWith(".mp3") || fileName.endsWith(".MP3")
                    || fileName.endsWith(".AAC") || fileName.endsWith(".aac")
                    || fileName.endsWith(".WAV") || fileName.endsWith(".wav")) {
                return FileType.music;
            }


            if (fileName.endsWith(".mp4") || fileName.endsWith(".MP4")
                    || fileName.endsWith(".RMVB") || fileName.endsWith(".rmvb")
                    || fileName.endsWith(".avi") || fileName.endsWith(".AVI")
                    || fileName.endsWith(".mov") || fileName.endsWith(".MOV")) {
                return FileType.video;
            }

            if (fileName.endsWith(".txt") || fileName.endsWith(".TXT")

            ) {
                return FileType.text;
            }


        }
        return FileType.other;
    }
}

