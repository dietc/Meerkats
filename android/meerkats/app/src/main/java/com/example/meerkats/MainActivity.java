package com.example.meerkats;


import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;

import android.app.ListActivity;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import com.google.gson.Gson;


public class MainActivity extends AppCompatActivity {

    private static Gson gson = new Gson();
    private ExecutorService threadPool;

    ///private Handler handler;

    private TCPMeerkats tcpMeerkats = new TCPMeerkats();

    public static byte[] func(File file, TCPMeerkats tcpMeerkats) {
        String PATH = "/storage/emulated/0/Android/data/com.example.meerkats/files/Meerkats/";
        byte[] md5;
        byte[] fileData = new byte[10000];
        ArrayList<FileInfoJson> fileInfoJsons = new ArrayList<>(100);
        File[] fs = file.listFiles();
        for (File f : fs) {
            if (f.isDirectory()) {
                func(f, tcpMeerkats);
            }
            if (f.isFile()) {
                try {
                    FileInputStream fis = new FileInputStream(f);
                    fis.read(fileData, 0, fileData.length);
                } catch (IOException e) {
                    e.printStackTrace();
                }
                md5 = tcpMeerkats.getMd5Hash(fileData);
                ArrayList<Byte> md5List = new ArrayList<>(16);
                for (int i = 0; i < 16; i++) {
                    md5List.add(md5[i]);
                }
                FileInfoJson fileInfoJson = new FileInfoJson(f.getPath().replace(PATH, ""), (byte) 0x01, md5List);
                fileInfoJsons.add(fileInfoJson);
            }
        }

        String messageBodyString = gson.toJson(fileInfoJsons);
        byte[] messageBody = messageBodyString.getBytes();
        byte[] messageBodyByte = tcpMeerkats.buildDataPackageForPull(messageBody, (byte) 0x02, (byte) 0x03);
        return messageBodyByte;
    }
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        threadPool = Executors.newCachedThreadPool();
        threadPool.execute(new Runnable() {
            @Override
            public void run() {
                File file = new File(getExternalFilesDir("Meerkats").getPath());
                file.mkdir();
                String rootPath = file.getPath();
                byte[] messageBodyByte = func(file,tcpMeerkats);

                tcpMeerkats.createInstance();
                tcpMeerkats.connectSocket();

                tcpMeerkats.sendMessage(messageBodyByte);


                tcpMeerkats.checkCMDFlag(tcpMeerkats.receiveMessage());

            }
        });
    }
}
