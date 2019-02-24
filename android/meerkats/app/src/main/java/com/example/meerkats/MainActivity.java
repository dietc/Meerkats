package com.example.meerkats;

import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import com.example.meerkats.TCPMeerkats;

public class MainActivity extends AppCompatActivity {


    private ExecutorService threadPool;

    private Handler handler;

    private TCPMeerkats tcpMeerkats = new TCPMeerkats();

    private Button connect, send, receive;

    private TextView result;

    private byte[] messageBody = {0x68, 0x65, 0x6c, 0x6c, 0x6f};


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        connect = (Button) findViewById(R.id.connect);
        send = (Button) findViewById(R.id.send);
        receive = (Button) findViewById(R.id.receive);
        result = (TextView) findViewById(R.id.result);


        threadPool = Executors.newCachedThreadPool();


        connect.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                threadPool.execute(new Runnable() {
                    @Override
                    public void run() {


                        tcpMeerkats.createInstance();
                        tcpMeerkats.connectSocket();

                    }
                });
            }
        });


        send.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                threadPool.execute(new Runnable() {
                    @Override
                    public void run() {

                        byte[] sendMessage = tcpMeerkats.buildDataPackage(messageBody, (byte) 0x01, (byte) 0x01);
                        tcpMeerkats.sendMessage(sendMessage);
                        System.out.println(sendMessage);


                    }
                });
            }
        });


    };

}