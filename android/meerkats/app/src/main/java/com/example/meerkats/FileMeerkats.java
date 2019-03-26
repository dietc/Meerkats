package com.example.meerkats;

import android.content.Context;
import android.content.Intent;
import android.net.Uri;

import com.example.meerkats.bean.FileType;

import org.json.JSONArray;
import org.json.JSONObject;

import java.io.File;
import java.util.Comparator;

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
                    || fileName.endsWith(".xml") || fileName.endsWith(".XML") ||
                    fileName.endsWith("log") || fileName.endsWith(".LOG")) {

                return FileType.text;
            }

            if (fileName.endsWith(".zip") || fileName.endsWith(".ZIP")
                    || fileName.endsWith(".rar") || fileName.endsWith(".rar")) {
                return FileType.zip;
            }

            return FileType.other;
        }


    }

    //A-Z display files

    public static Comparator comparator = new Comparator<File>() {
        @Override
        public int compare(File f1, File f2) {
            if (f1.isDirectory() && f2.isFile()) {
                return -1;
            } else if (f1.isFile() && f2.isDirectory()){
                return 1;
            } else {
                return f1.getName().compareTo(f2.getName());
            }
        }
    };

    //count child files number

    public static int countChildFiles(File file){
        int i =0;
        if (file.isDirectory()){
            File[] files = file.listFiles();
            for (File f : files){
                i++;
            }
        }
        return i;
    }

    public static JSONArray getAllFiles(String dirPath, String _type) {
        File f = new File(dirPath);
        if (!f.exists()) {//判断路径是否存在
            return null;
        }

        File[] files = f.listFiles();

        if(files==null){//判断权限
            return null;
        }

        JSONArray fileList = new JSONArray();
        for (File _file : files) {//遍历目录
            if(_file.isFile() && _file.getName().endsWith(_type)){
                String _name=_file.getName();
                String filePath = _file.getAbsolutePath();//获取文件路径
                String fileName = _file.getName().substring(0,_name.length()-4);//获取文件名
//                Log.d("LOGCAT","fileName:"+fileName);
//                Log.d("LOGCAT","filePath:"+filePath);
                try {
                    JSONObject _fInfo = new JSONObject();
                    _fInfo.put("name", fileName);
                    _fInfo.put("path", filePath);
                    fileList.put(_fInfo);
                }catch (Exception e){
                }
            } else if(_file.isDirectory()){//查询子目录
                getAllFiles(_file.getAbsolutePath(), _type);
            } else{
            }
        }
        return fileList;
    }

    //open music type files

    public static void openMusicIntent( Context context , File file ){
        Uri path = Uri.fromFile(file);
        Intent intent = new Intent(Intent.ACTION_VIEW);
        intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
        intent.setDataAndType(path, "audio/*");
        context.startActivity(intent);
    }

    }




