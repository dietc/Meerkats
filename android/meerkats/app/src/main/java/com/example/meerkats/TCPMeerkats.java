package com.example.meerkats;



import android.content.Context;
import android.icu.text.SymbolTable;

import com.google.gson.Gson;


import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.File;
import java.io.FileOutputStream;
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
import java.util.Arrays;
import java.util.List;




public class TCPMeerkats extends Thread{

    static String ip = "178.128.45.7";

    private static int port = 4356;

    private static Socket socketClient;

    private static InputStream is;

    private static InputStreamReader isr;

    private static BufferedReader br;

    private static OutputStream os;

    private static String PATH = "/data/data/com.example.meerkats/files";

    Context context;

    //Create a socket client instance

    public void createInstance(){

        socketClient = new Socket();
    }


    //Connect the socket

    public void connectSocket(){

        try{

            SocketAddress remoteAddress = new InetSocketAddress(ip, port);
            socketClient.connect(remoteAddress);
            System.out.println(socketClient.isConnected());

        } catch (IOException e){
            System.out.println("ERROR!TRY AGAIN!");

        }

    }


    ///Receive Message
    public byte[] receiveMessage() {


        int tcpHeaderLength = 10;
        int tcpBodyLength = 0;

        byte[] tcpHeader = new byte[tcpHeaderLength];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(tcpHeader,0,10);


        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }



        tcpBodyLength = ((int)tcpHeader[8] << 8) + (int)(tcpHeader[9]);
        byte[] recvBytes = new byte[tcpBodyLength];
        try {
            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(recvBytes, 0, tcpBodyLength);
        }catch(IOException e){
            System.out.println("ERROR! TRY AGAIN!");
        }


        byte[] md5 = new byte[16];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(md5,0,16);

        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }



        byte[] endFlag = new byte[2];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(endFlag,0,2);

        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }


        if (endFlag[0] == (byte)0xff && endFlag[1] == (byte)0xee) {

            return unpackData(recvBytes);

        } else {

            return null;
       }
    }

    public String receiveMessageForDownload() {

        byte[] recvBytes;

        byte packetNum;

        int tcpHeaderLength = 10;
        int tcpBodyLength;

        byte[] tcpHeader = new byte[tcpHeaderLength];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(tcpHeader, 0, 10);

        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }



        tcpBodyLength = ((tcpHeader[8] & 0xff) << 8) + (tcpHeader[9] & 0xff);


        byte[] packetType = new byte[1];
        int fileDataLength;
        byte[] fileJson = new byte[300];
        fileDataLength = tcpBodyLength - 300 - 1;
        recvBytes = new byte[fileDataLength];
        try {


            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(packetType, 0, 1);


        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(fileJson, 0, 300);


        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }


        int endFileJsonFlag = 0;
        int i = 0;
        for (byte flag : fileJson) {
            if (flag == 0x00) {
                endFileJsonFlag = i;
                break;
            }
            i++;
        }


        byte[] fileJsonn = new byte[endFileJsonFlag];

        System.arraycopy(fileJson, 0, fileJsonn, 0, endFileJsonFlag);


        String fileJsonString = new String(fileJsonn);
        FileIndexJson fJson = new Gson().fromJson(fileJsonString, FileIndexJson.class);
        String fileName = fJson.getName();
        packetNum = fJson.getNum();



        try {
            DataInputStream dis = new DataInputStream(socketClient.getInputStream());
            dis.readFully(recvBytes,0,fileDataLength);

        } catch (IOException e) {
            System.out.println("ERROR! TRY AGAIN!");
        }


        byte[] md5 = new byte[16];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(md5, 0, 16);
        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }

        byte[] endFlag = new byte[2];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(endFlag, 0, 2);

            for (byte edn : endFlag) {
                System.out.printf("%x\n", edn);
            }
        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }


        try {
        if (endFlag[0] == (byte) 0xff && endFlag[1] == (byte) 0xee) {

            System.out.println("8");
            FileOutputStream fos = context.getApplicationContext().openFileOutput(fileName, Context.MODE_PRIVATE);
            fos.write(recvBytes, 0, recvBytes.length);

            System.out.println("8");


            if (packetNum == 1) {
                fos.flush();
                fos.close();
                System.out.println("9");
            } else {
                while (packetNum > 1) {

                    byte[] fileData = receiveMessage();
                    fos.write(fileData, 0, fileDataLength);
                    System.out.println("1");

                    packetNum--;

                }
                fos.flush();
                fos.close();
                System.out.println("10");

            }
        } else {
            return null;
        }

    } catch (IOException e) {
        System.out.println("ERROR!TRY AGAIN!");
    }

        return fileName +  "DOWNLOAD FINISHED!";

}

    ///Send Message

    public void sendMessage (byte[] sendBytes){

        ///Check if connected
        if (socketClient.isConnected()) {
            try {
                os = socketClient.getOutputStream();
                os.write(sendBytes);
                os.flush();
            } catch (IOException e) {
                System.out.println("ERROR! TRY AGAIN!");

            }
        }

    }


   /* public void uploadFile (String[] fileName) {

        byte packetType = 0x20;
        byte deviceID = 0x03;
        byte[] fileData;

        int messageBodyLengthMax = 65233;

        for (int i = 0; i < fileName.length; i++){

            File f = new File(PATH + fileName[i]);

            long fileLength = f.length();

            int messageBodyLength = 0;

            if (fileLength <= messageBodyLengthMax){

                messageBodyLength = (int)fileLength;

                byte[] messageBodyByte = new byte[messageBodyLength + 30 + 300];

                byte[] initiator = {0x11, (byte)0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e};

                System.arraycopy(initiator, 0, messageBodyByte,0, initiator.length);

                int contextLength = messageBodyLength + 2 + 300;

                byte[] length = { (byte)(contextLength >> 8), (byte)(contextLength & 0xff)};

                System.arraycopy(length,0, messageBodyByte,8, length.length);

                messageBodyByte[10] = packetType;

                messageBodyByte[11] = deviceID;

                byte[] fileJson = new byte[300];

                FileIndexJson fjs = new FileIndexJson(fileName[i], (byte)0x1, (byte)0);
                Gson gson = new Gson();


                String fjsString = gson.toJson(fjs);

                byte[] fjsByte = fjsString.getBytes();

                System.arraycopy(fjsByte,0,fileJson,0,fjsString.length());

                System.arraycopy(fileJson,0,messageBodyByte,12, fileJson.length);

                fileData = new byte[messageBodyLength];

                System.arraycopy(fileData,0,messageBodyByte,12+300,fileData.length);

                byte[] md5 = getCheckSum(packetType, deviceID, fileData,true);

                System.arraycopy(md5,0, messageBodyByte,12 + 300 + messageBodyLength, md5.length);

                messageBodyByte[messageBodyLength + 28 + 300] = (byte) 0xff;
                messageBodyByte[messageBodyLength + 29 + 300] = (byte) 0xee;

                sendMessage(messageBodyByte);

            } else {

                double packetNum =  Math.ceil((double)fileLength / (double)messageBodyLengthMax);
                byte index = 0;
                fileData = null;
                byte[] messageBodyByte;
                while (packetNum > 0) {
                    if (packetNum != 1) {
                        messageBodyLength = messageBodyLengthMax;
                    } else {
                        messageBodyLength = (int) (fileLength - messageBodyLengthMax * index);
                    }

                    messageBodyByte = new byte[messageBodyLength + 30 + 300];
                    byte[] initiator = {0x11, (byte) 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e};
                    System.arraycopy(initiator, 0, messageBodyByte, 0, initiator.length);
                    int contextLength = messageBodyLength + 2 + 300;
                    byte[] length = {(byte) (contextLength >> 8), (byte) (contextLength & (byte) 0xff)};
                    System.arraycopy(length, 0, messageBodyByte, 8, length.length);

                    messageBodyByte[10] = packetType;
                    messageBodyByte[11] = deviceID;

                    byte[] fileJson = new byte[300];

                    FileIndexJson fjs = new FileIndexJson(fileName[i], (byte) packetNum, index);

                    byte[] jsonByte = new byte[1];
                    System.arraycopy(jsonByte, 0, fileJson, 0, jsonByte.length);

                    System.arraycopy(fileJson, 0, messageBodyByte, 12, fileJson.length);

                    fileData = new byte[messageBodyLength];

                    System.arraycopy(fileData, 0, messageBodyByte, 12 + 300, fileData.length);

                    byte[] md5 = getCheckSum(packetType, deviceID, fileData, true);

                    System.arraycopy(md5, 0, messageBodyByte, 12 + 300 + messageBodyLength, md5.length);

                    messageBodyByte[messageBodyLength + 28 + 300] = (byte) 0xff;
                    messageBodyByte[messageBodyLength + 29 + 300] = (byte) 0xee;

                    packetNum--;
                    index++;

                    sendMessage(messageBodyByte);

                }





            }

            }



    }



*/






    ///Build data package
    public byte[] buildDataPackageForPull(byte[] messageBody, byte packetType, byte deviceID ) {

        int messageBodyLength = 0;
        messageBodyLength = messageBody.length;

        byte[] messageBodyByte = new byte[messageBodyLength + 30];

        byte[] Initiator = {0x11, (byte)0xff, 0x6c,0x6f, 0x6e, 0x64, 0x6f, 0x6e};
        System.arraycopy(Initiator, 0, messageBodyByte, 0, Initiator.length);


        byte contextLength = (byte)(messageBodyLength + 2);

        byte[] length = { (byte)( contextLength >> 8), (byte)(contextLength & 0xff) };

        System.arraycopy(length, 0, messageBodyByte, 8, length.length);

        messageBodyByte[10] = packetType;

        messageBodyByte[11] = deviceID;

        System.arraycopy(messageBody, 0, messageBodyByte, 12, messageBody.length);


        byte[] md5 = getCheckSum(packetType, deviceID, messageBody, true);

        System.arraycopy(md5, 0, messageBodyByte, 12 + messageBody.length, md5.length);

        messageBodyByte[messageBodyLength + 28] = (byte)0xff;

        messageBodyByte[messageBodyLength + 29] = (byte)0xee;

        return messageBodyByte;


    }



    public byte[] unpackData(byte[] recvMsg) {

        if (recvMsg.length != 0) {

            int tcpBodyLength = recvMsg.length;
            int contextLength = tcpBodyLength - 1;

            int packetType = recvMsg[0];


            byte[] msg = new byte[contextLength];

            System.arraycopy(recvMsg, 1, msg, 0, contextLength);

            return msg;

        } else {

            return null;
        }

    }

    //public List<String> checkUpload(byte[] recvMsg) {

   //     String resultStr = new String(recvMsg);

    //    String cmd = null;
     //

   // }


    private byte[] getCheckSum(byte packetType, byte deviceId, byte[] msg, boolean sendOrReceive ) {

        int index = 1;

        if (sendOrReceive) {
            index++;
        }

        byte[] privateKey = {0x61, 0x61, 0x61, 0x61, 0x61};
        byte[] checkSum = new byte[index + msg.length + privateKey.length];
        checkSum[0] = packetType;

        if (sendOrReceive) {
            checkSum[1] = deviceId;
        }

        System.arraycopy(msg, 0, checkSum, index, msg.length);

        System.arraycopy(privateKey, 0, checkSum, msg.length + index, privateKey.length);

        byte[] md5 = getMd5Hash(checkSum);

        return md5;


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


