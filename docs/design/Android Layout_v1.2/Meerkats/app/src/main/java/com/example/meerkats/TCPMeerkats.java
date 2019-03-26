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

    public String downloadFile(List<FileTypeTCP> fileName) {

        for (FileTypeTCP ja : fileName) {
            byte[] msgBody = ja.ext.getBytes();
            byte[] messageBodyByteForDownload = new byte[30 + msgBody.length];
            messageBodyByteForDownload = buildDataPackageForPull(msgBody, (byte) 0x21, deviceID);
            sendMessage(messageBodyByteForDownload);
            receiveMessageForDownload(ja.type);
        }
        return "DOWNLOAD SUCCESS!";
    }

    public String renameFile(List<Rename> fileInfo) {
        for (Rename ja : fileInfo) {

            File f = new File(PATH + ja.name);
            if (f.exists() && !f.isDirectory()) {
                f.renameTo(new File(PATH + ja.name));

            } else {
                System.out.println("FILE NOT EXIST");
            }
        }
        return "RENAME SUCCESSFULLY!";
    }

    public String backupFile(List<Backup> fileName) {
        for (Backup ja : fileName) {

            File f = new File(PATH + ja.name);
            if (f.exists() && !f.isDirectory()) {
                f.renameTo(new File(NEWPATH + ja.name + ".bak"));
            } else {
                System.out.println("FILE NOT EXIST");
            }
        }
        return "BACKUP SUCCESSFULLY!";
    }


    public String deleteFile(List<Delete> fileName) {
        for (Delete ja : fileName) {
            File f = new File(PATH + ja.name);
            if (f.exists() && !f.isDirectory()) {
                f.delete();
            } else {
                System.out.println("FILE NOT EXIST");
            }
        }
        return "DELETE SUCCESSFULLY!";

    }


    public void uploadFile(String[] fileName) {

        byte packetType = 0x20;
        byte deviceID = 0x03;
        byte[] fileData;

        int messageBodyLengthMax = 65233;

        for (int i = 0; i < fileName.length; i++) {


            File f = new File(PATH + "/" + fileName[i]);

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

                FileIndexJson fjs = new FileIndexJson(fileName[i], (byte) 0x1, (byte) 0);
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

                sendMessage(messageBodyByte);

            } else {

                double packetNum = Math.ceil((double) fileLength / (double) messageBodyLengthMax);
                byte index = 0;
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

                    packetNum--;
                    index++;

                    sendMessage(messageBodyByte);

                }


            }

        }


    }


    ///Build data package
    public byte[] buildDataPackageForPull(byte[] messageBody, byte packetType, byte deviceID) {

        int messageBodyLength = 0;
        messageBodyLength = messageBody.length;

        byte[] messageBodyByte = new byte[messageBodyLength + 30];

        byte[] Initiator = {0x11, (byte) 0xff, 0x6c, 0x6f, 0x6e, 0x64, 0x6f, 0x6e};
        System.arraycopy(Initiator, 0, messageBodyByte, 0, Initiator.length);


        byte contextLength = (byte) (messageBodyLength + 2);

        byte[] length = {(byte) (contextLength >> 8), (byte) (contextLength & 0xff)};

        System.arraycopy(length, 0, messageBodyByte, 8, length.length);

        messageBodyByte[10] = packetType;

        messageBodyByte[11] = deviceID;

        System.arraycopy(messageBody, 0, messageBodyByte, 12, messageBody.length);


        byte[] md5 = getCheckSum(packetType, deviceID, messageBody, true);

        System.arraycopy(md5, 0, messageBodyByte, 12 + messageBody.length, md5.length);

        messageBodyByte[messageBodyLength + 28] = (byte) 0xff;

        messageBodyByte[messageBodyLength + 29] = (byte) 0xee;

        return messageBodyByte;


    }


    public byte[] unpackData(byte[] recvMsg) {

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

    public void checkCMDFlag(byte[] recvMsg, FileCheck fCheck) {
        String resultStr = new String(recvMsg);
        int cmdFlag = 0;
        List<FileTypeTCP> uploadList = new ArrayList<>();
        List<FileTypeTCP> downloadList = new ArrayList<>();
        List<Rename> renameList = new ArrayList<>();
        List<Backup> backupList = new ArrayList<>();
        fCheck.upload = uploadList;
        fCheck.download = downloadList;
        fCheck.rename = renameList;
        fCheck.backup = backupList;
        try {
            JSONArray jArray = new JSONArray(resultStr);
            for (int i = 0; i < jArray.length(); i++) {
                JSONObject jObject = jArray.getJSONObject(i);
                cmdFlag = Integer.parseInt(jObject.getString("Cmd"));
                switch (cmdFlag) {
                    case 1:
                        FileTypeTCP uploadAll = new FileTypeTCP(jObject.getString("Name"), jObject.getString("Ext"), 0);
                        fCheck.upload.add(uploadAll);
                        break;
                    case 2:
                        FileTypeTCP downloadAll = new FileTypeTCP(jObject.getString("Name"), jObject.getString("Ext"), 1);
                        fCheck.download.add(downloadAll);
                        break;
                    case 3:
                        Rename rename = new Rename(jObject.getString("Name"), jObject.getString("Ext"));
                        fCheck.rename.add(rename);
                        break;
                    case 4:
                        //differ upload
                        break;
                    case 5:
                        //differ download
                        break;
                    case 6:
                        Delete delete = new Delete(jObject.getString("Name"));
                        fCheck.delete.add(delete);
                        break;
                    case 7:
                        Backup backup = new Backup(jObject.getString("Name"));
                        fCheck.backup.add(backup);
                        break;
                }
            }

        } catch (JSONException e) {
            System.out.println("JSON CONVERT FAILED!");
        }


    }


    private byte[] getCheckSum(byte packetType, byte deviceId, byte[] msg, boolean sendOrReceive) {

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


    private static byte[] getMd5Hash(byte[] byteData) {

        byte[] md5 = new byte[16];

        try {

            MessageDigest messageDigest = MessageDigest.getInstance("MD5");

            md5 = messageDigest.digest(byteData);


        } catch (NoSuchAlgorithmException e) {
            System.out.println("ERROR! TRY AGAIN!");
        }

        return md5;

    }


    public static JSONArray getAllFiles(String dirPath, String _type) {
        File f = new File(dirPath);
        if (!f.exists()) {//判断路径是否存在
            return null;
        }

        File[] files = f.listFiles();

        if (files == null) {//判断权限
            return null;
        }

        JSONArray fileList = new JSONArray();
        for (File _file : files) {//遍历目录
            if (_file.isFile() && _file.getName().endsWith(_type)) {
                String _name = _file.getName();
                String filePath = _file.getAbsolutePath();//获取文件路径
                String fileName = _file.getName().substring(0, _name.length() - 4);//获取文件名
//                Log.d("LOGCAT","fileName:"+fileName);
//                Log.d("LOGCAT","filePath:"+filePath);
                try {
                    JSONObject _fInfo = new JSONObject();
                    _fInfo.put("name", fileName);
                    _fInfo.put("path", filePath);
                    fileList.put(_fInfo);
                } catch (Exception e) {
                }
            } else if (_file.isDirectory()) {//查询子目录
                getAllFiles(_file.getAbsolutePath(), _type);
            } else {
            }
        }
        return fileList;
    }


}


