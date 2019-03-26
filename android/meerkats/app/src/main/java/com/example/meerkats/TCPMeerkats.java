package com.example.meerkats;



import android.content.Context;
import android.icu.text.SymbolTable;
import android.support.v7.app.AppCompatActivity;
import android.widget.Switch;

import com.example.meerkats.bean.FileType;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;



import java.io.FileInputStream;
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
import java.util.jar.Attributes;

import static java.lang.System.in;


public class TCPMeerkats extends Thread {

    static String ip = "178.128.45.7";

    private static int port = 4356;

    private static Socket socketClient;
    private static byte deviceID = 0x03;


    private static String PATH = "/data/data/com.example.meerkats/files";
    private static String NEWPATH = "/data/data/com.example.meerkats/files/backup";


    //Create a socket client instance

    public void createInstance() {

        socketClient = new Socket();
    }


    //Connect the socket

    public void connectSocket() {

        try {

            SocketAddress remoteAddress = new InetSocketAddress(ip, port);
            socketClient.connect(remoteAddress);
            System.out.println(socketClient.isConnected());

        } catch (IOException e) {
            System.out.println("ERROR!TRY AGAIN!");

        }

    }


    ///Receive Message
    public byte[] receiveMessage() {


        int tcpHeaderLength = 10;

        byte[] tcpHeader = new byte[tcpHeaderLength];

        try {

            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(tcpHeader, 0, 10);


        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }


        int tcpBodyLength = ((tcpHeader[8] & 0xff) << 8) + (tcpHeader[9] & 0xff);
        byte[] recvBytes = new byte[tcpBodyLength];


        try {
            InputStream is = socketClient.getInputStream();
            DataInputStream dis = new DataInputStream(is);
            dis.readFully(recvBytes);
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

        } catch (IOException e) {

            System.out.println("ERROR! TRY AGAIN!");
        }


        if (endFlag[0] == (byte) 0xff && endFlag[1] == (byte) 0xee) {

            return unpackData(recvBytes);

        } else {

            return null;
        }
    }

    public String receiveMessageForDownload(int fileType) {

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
            dis.readFully(recvBytes, 0, fileDataLength);

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
                File file = new File(PATH, fileName);
                FileOutputStream fos = new FileOutputStream(file, false);

                switch (fileType) {
                    case 0:
                        fos.write(recvBytes);
                        break;
                    case 1:
                        fos.write(recvBytes);
                        break;
                }

                if (packetNum == 1) {
                    fos.flush();
                    fos.close();
                    System.out.println("9");
                    return fileName + "DOWNLOADED!";
                } else {

                    FileOutputStream fos1 = new FileOutputStream(file, true);
                    while (packetNum > 1) {

                        byte[] fileData = unpackData(receiveMessage());
                        switch (fileType) {
                            case 0:
                                fos1.write(fileData);
                                break;
                            case 1:
                                fos1.write(fileData);
                                break;
                        }
                        packetNum--;
                    }
                    fos1.flush();
                    fos1.close();
                    System.out.println("10");

                }
            } else {
                return null;
            }

        } catch (IOException e) {
            System.out.println("ERROR!TRY AGAIN!");
        }

        return fileName + "DOWNLOAD FINISHED!";
    }

    ///Send Message

    public void sendMessage(byte[] sendBytes) {

        ///Check if connected
        if (socketClient.isConnected()) {
            try {
                OutputStream os = socketClient.getOutputStream();
                os.write(sendBytes);
                os.flush();
            } catch (IOException e) {
                System.out.println("ERROR! TRY AGAIN!");

            }
        }

    }

    public String downloadFile(String fileName, int downloadType) {

        byte[] messageBody = fileName.getBytes();
        byte[] messageBodyForDownload = new byte[30 + messageBody.length];
        messageBodyForDownload = buildDataPackageForPull(messageBody, (byte) 0x21, deviceID);
        sendMessage(messageBodyForDownload);
        receiveMessageForDownload(downloadType);
        return "DOWNLOAD SUCCESSED!";
    }

    public String renameFile(String fileName, String newFileName) {


        File f = new File(PATH + fileName);
        if (f.exists()) {
            f.renameTo(new File(PATH + newFileName));
        } else {
            System.out.println("FILE NOT EXIST");
        }
        return "RENAME SUCCESSFULLY!";
    }

    public String backupFile(String fileName) {

        File f = new File(PATH + fileName);
        if (f.exists() && !f.isDirectory()) {
            f.renameTo(new File(NEWPATH + fileName + ".bak"));
        } else {
            System.out.println("FILE NOT EXIST");
        }
        return "BACKUP SUCCESSFULLY!";
    }


    public String deleteFile(String fileName) {
        File f = new File(PATH + fileName);
        if (f.exists() && !f.isDirectory()) {
            f.delete();
        } else {
            System.out.println("FILE NOT EXIST");
        }
        return "DELETE SUCCESSFULLY!";
    }


    public String uploadFile(String fileName, int uploadType) {

        byte packetType = 0x20;
        byte deviceID = 0x03;
        byte[] fileData;

        int messageBodyLengthMax = 65233;


        File f = new File(PATH + "/" + fileName);

        long fileLength = f.length();

        int messageBodyLength = 0;

        if (fileLength <= messageBodyLengthMax) {

            messageBodyLength = (int) fileLength;

            byte[] messageBodyByte = new byte[messageBodyLength + 30 + 300];

            byte[] initiator = {0x11, (byte) 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e};

            System.arraycopy(initiator, 0, messageBodyByte, 0, initiator.length);

            int contextLength = messageBodyLength + 2 + 300;

            byte[] length = {(byte) (contextLength >> 8), (byte) (contextLength & 0xff)};

            System.arraycopy(length, 0, messageBodyByte, 8, length.length);

            messageBodyByte[10] = packetType;

            messageBodyByte[11] = deviceID;

            byte[] fileJson = new byte[300];

            FileIndexJson fjs = new FileIndexJson(fileName, (byte) 0x1, (byte) 0);
            Gson gson = new Gson();
            String fjsString = gson.toJson(fjs);

            byte[] fjsByte = fjsString.getBytes();

            System.arraycopy(fjsByte, 0, fileJson, 0, fjsString.length());

            System.arraycopy(fileJson, 0, messageBodyByte, 12, fileJson.length);
            fileData = new byte[messageBodyLength];
            switch (uploadType) {
                case 0:
                    try {
                        //这里可能有问题 fileData的值不确定能不能传过去
                        FileInputStream fis = new FileInputStream(f);
                        fis.read(fileData, 0, messageBodyLength);
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                    System.arraycopy(fileData, 0, messageBodyByte, 12 + 300, messageBodyLength);
                    break;
                case 1:
                    try {
                        FileInputStream fis = new FileInputStream(f);
                        fis.read(fileData, 0, messageBodyLength);
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                    System.arraycopy(fileData, 0, messageBodyByte, 12 + 300, messageBodyLength);
                    break;
            }

            byte[] md5 = getCheckSum(packetType, deviceID, fileData, true);

            System.arraycopy(md5, 0, messageBodyByte, 12 + 300 + messageBodyLength, md5.length);

            messageBodyByte[messageBodyLength + 28 + 300] = (byte) 0xff;
            messageBodyByte[messageBodyLength + 29 + 300] = (byte) 0xee;

            sendMessage(messageBodyByte);

            byte[] resflag = receiveMessage();
            if (resflag[0] != 0x6f || resflag[1] != 0x6b) {
                return "UPLOAD FAILED!:(" + fileName;
            }
        } else {

            double packetNumD = Math.ceil((double) fileLength / (double) messageBodyLengthMax);
            byte packetNumB = (byte) packetNumD;
            byte index = 0;
            byte[] messageBodyByte;
            while (packetNumD > 0) {
                if (packetNumD != 1) {
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

                FileIndexJson fjs = new FileIndexJson(fileName, packetNumB, index);
                Gson gson = new Gson();
                String fjsString = gson.toJson(fjs);

                byte[] fjsByte = fjsString.getBytes();

                System.arraycopy(fjsByte, 0, fileJson, 0, fjsString.length());

                System.arraycopy(fileJson, 0, messageBodyByte, 12, fileJson.length);


                fileData = new byte[messageBodyLength];

                try {
                    FileInputStream fis = new FileInputStream(f);
                    fis.read(fileData, 0, (int) messageBodyLength);

                } catch (IOException e) {
                    e.printStackTrace();
                }

                System.arraycopy(fileData, 0, messageBodyByte, 12 + 300, messageBodyLength);

                byte[] md5 = getCheckSum(packetType, deviceID, fileData, true);

                System.arraycopy(md5, 0, messageBodyByte, 12 + 300 + messageBodyLength, md5.length);

                messageBodyByte[messageBodyLength + 28 + 300] = (byte) 0xff;
                messageBodyByte[messageBodyLength + 29 + 300] = (byte) 0xee;

                packetNumD--;
                index++;

                sendMessage(messageBodyByte);

            }
            byte[] resflag = receiveMessage();
            if (resflag[0] != 0x6f || resflag[1] != 0x6b) {
                return "UPLOAD FAILED!:(" + fileName;
            }
        }
        return "UPLOAD SUCCESSED!";
    }








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



    public byte[] unpackData(byte[] recvMsg)
    {

        if (recvMsg != null) {

            byte packetType = recvMsg[0];
            int tcpBodyLength = recvMsg.length;
            int contextLength = tcpBodyLength - 1;

            byte[] msg = new byte[contextLength];

            System.arraycopy(recvMsg, 1, msg, 0, contextLength);

            return msg;

        } else {

            return null;
        }

    }

    public void checkCMDFlag(byte[] recvMsg) {
        String resultStr = new String(recvMsg);
        int cmdFlag = 0;
        try {
            JSONArray jArray = new JSONArray(resultStr);
            for (int i = 0; i < jArray.length(); i ++){
                JSONObject jObject = jArray.getJSONObject(i);
                cmdFlag = Integer.parseInt(jObject.getString("Cmd"));
                switch (cmdFlag){
                    case 1:
                        uploadFile(jObject.getString("Name"),0);
                        break;
                    case 2:
                        downloadFile(jObject.getString("Name"),0);
                        break;
                    case 3:
                        renameFile(jObject.getString("Name"),jObject.getString("Ext"));
                        break;
                    case 4:
                        //differ upload
                        uploadFile(jObject.getString("Name"),1);
                        break;
                    case 5:
                        //differ download
                        downloadFile(jObject.getString("Name"),1);
                        break;
                    case 6:
                        deleteFile(jObject.getString("Name"));
                        break;
                    case 7:
                        backupFile(jObject.getString("Name"));
                        break;
                    }
                }

        }catch (JSONException e){
            System.out.println("JSON CONVERT FAILED!");
        }


    }


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


