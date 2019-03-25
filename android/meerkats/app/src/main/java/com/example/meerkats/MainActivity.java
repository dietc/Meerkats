package com.example.meerkats;

import android.content.Context;
import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.Socket;
import java.util.Arrays;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import com.example.meerkats.TCPMeerkats;
import com.google.gson.Gson;

public class MainActivity extends AppCompatActivity {


    private ExecutorService threadPool;

    ///private Handler handler;

    private static Context context;
    private TCPMeerkats tcpMeerkats = new TCPMeerkats();

    private Button connect, send, receive;

    private TextView result;

    private byte[] messageBody = {};

    private byte deviceID = (byte) 0x03;

    private byte packetType = (byte) 0x21;

    private static Socket socketClient;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);


        tcpMeerkats.createInstance();
        tcpMeerkats.connectSocket();

        tcpMeerkats.buildDataPackageForPull(messageBody,packetType,deviceID);

        String message = tcpMeerkats.receiveMessageForDownload(3);

        System.out.println(message);
    }
}
