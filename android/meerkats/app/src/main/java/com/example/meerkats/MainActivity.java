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

public class MainActivity extends AppCompatActivity {


    private ExecutorService threadPool;

    private Handler handler;

    private Socket socket;

    InputStream is;

    InputStreamReader isr;

    BufferedReader br;

    String response;

    OutputStream os;

    Byte sendByte;

    private Button connect, send, receive;

    private TextView result;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        connect = (Button) findViewById(R.id.connect);
        send = (Button) findViewById(R.id.send);
        receive = (Button) findViewById(R.id.receive);
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

        connect.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                threadPool.execute(new Runnable() {
                    @Override
                    public void run() {

                        try{

                            socket = new Socket("178.128.45.7", 4356);
                            System.out.println(socket.isConnected());
                            System.out.println("Connected!!!!!!!!!!!!");

                        } catch (IOException e){
                            e.printStackTrace();

                        }
                    }
                });
            }
        });


        receive.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                    threadPool.execute(new Runnable() {
                        @Override
                        public void run() {
                            try {
                                is = socket.getInputStream();
                                isr = new InputStreamReader(is);
                                br = new BufferedReader(isr);

                                response = br.readLine();

                                Message msg = Message.obtain();
                                msg.what = 0;
                                handler.sendMessage(msg);

                            }
                            catch (IOException e){
                                e.printStackTrace();
                            }
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
                        try {
                            os = socket.getOutputStream();

                            os.write(sendByte);

                            os.flush();
                        } catch (IOException e) {

                            e.printStackTrace();

                        }
                    }
                });
            }
        });

    };

};

