package com.example.meerkats;

import android.os.Message;

import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.math.BigInteger;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.nio.Buffer;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.List;


import java.util.concurrent.ThreadPoolExecutor;

public class TCPMeerkats extends Thread{

    static String ip = "178.128.45.7";

    private static int port = 4356;

    private static Socket socketClient;

    private static InputStream is;

    private static InputStreamReader isr;

    private static BufferedReader br;

    private static OutputStream os;


    
    public  static List<String> listMessage = new ArrayList<String>();

    ///Create a socket client instance

    public void createInstance(){

        socketClient = new Socket();
    }


    ///Connect the socket

    public void connectSocket(){

        try{

            SocketAddress remoteAddr = new InetSocketAddress(ip, port);
            socketClient.connect(remoteAddr);
            System.out.println(socketClient.isConnected());

        } catch (IOException e){
            System.out.println("ERROR!TRY AGAIN!");

        }

    }

    ///Receive Message
    public byte[] receiveMessage() {
        byte[] message = new byte[2014];
        try {
            OutputStream socketOutputStream = socketClient.getOutputStream();
            socketOutputStream.write(message);
        } catch (IOException e) {
            System.out.println("ERROR! TRY AGAIN!");
        }
        return message;
    }

    ///Receive Message
    ///public void receiveMessage() {


    ///}



    ///Send Message

    public void sendMessage (byte[] sendBytes){

        ///Check if connected
        if (socketClient.isConnected()){

            try{
                os = socketClient.getOutputStream();
                os.write(sendBytes);
                os.flush();
            }catch (IOException e){
                System.out.println("ERROR! TRY AGAIN!");



            }

            }

        }

    ///Build data package
    public byte[] buildDataPackage(byte[] messageBody, byte packageType, byte deviceID ) {

        int messageBodyLength = 0;
        messageBodyLength = messageBody.length;

        byte[] messageBodyByte = new byte[messageBodyLength + 30];

        byte[] Initiator = {0x11, (byte)0xff, 0x6c,0x6f, 0x6e, 0x64, 0x6f, 0x6e};
        System.arraycopy(Initiator, 0, messageBodyByte, 0, Initiator.length);


        byte contextLength = (byte)(messageBodyLength + 2);

        byte[] length = { (byte)( contextLength >> 8), (byte)(contextLength & 0xff) };

        System.arraycopy(length, 0, messageBodyByte, 8, length.length);

        messageBodyByte[10] = packageType;

        messageBodyByte[11] = deviceID;

        System.arraycopy(messageBody, 0, messageBodyByte, 12, messageBody.length);

        byte[] privateKey = {0x61, 0x61, 0x61, 0x61, 0x61};

        byte[] checkSum = new byte[messageBodyLength + privateKey.length + 2];

        checkSum[0] = packageType;
        checkSum[1] = deviceID;

        System.arraycopy(messageBody, 0, checkSum, 2, messageBody.length);

        System.arraycopy(privateKey, 0, checkSum, messageBody.length +2, privateKey.length);

        byte[] md5 = getMd5Hash(checkSum);

        System.arraycopy(md5, 0, messageBodyByte, 12 + messageBody.length, md5.length);

        messageBodyByte[messageBodyLength + 28] = (byte)0xff;
        messageBodyByte[messageBodyLength + 29] = (byte)0xee;

        return messageBodyByte;


    }

    public byte[] unpackData(byte[] recvMsg){

        int contextLength = 11;

        int messageBodyLength = (int)recvMsg[9] - 1;

        byte[] msg = new byte[messageBodyLength];

        System.arraycopy(recvMsg, contextLength, msg, 0, messageBodyLength);

        return msg;

    }


    private static byte[] getMd5Hash(byte[] byteData){

        byte[] md5 = new byte[16];

        try {

            MessageDigest messageDigest = MessageDigest.getInstance("MD5");

            md5 = messageDigest.digest(byteData);


        } catch (NoSuchAlgorithmException e){
            System.out.println("ERROR! TRY AGAIN!");
        }

        return md5;

    }





}


