package com.example.meerkats;

import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.widget.Button;
import android.widget.TextView;
import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class MainActivity extends AppCompatActivity {


    private ExecutorService threadPool;

    private Handler handler;

    private Socket socket;

    InputStream is;

    InputStreamReader isr;

    BufferedReader br;

    String response;

    OutputStream os;

    private Button connect, send;

    private TextView result;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        connect = (Button) findViewById(R.id.connect);
        send = (Button) findViewById(R.id.send);
        result = (TextView) findViewById(R.id.result);


        threadPool = Executors.newCachedThreadPool();

        handler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case 0:
                        result.setText(response);
                        break;
                }
            }
        };

    };

};

