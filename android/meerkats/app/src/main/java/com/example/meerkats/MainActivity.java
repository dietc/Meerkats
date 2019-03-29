package com.example.meerkats;

import android.Manifest;
import android.os.AsyncTask;
import android.os.Bundle;
import android.os.Environment;
import android.os.Looper;
import android.support.annotation.NonNull;
import android.support.design.widget.FloatingActionButton;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.LinearLayoutManager;
import android.support.v7.widget.RecyclerView;
import android.view.KeyEvent;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.Toast;

import com.google.gson.Gson;
import com.yanzhenjie.permission.AndPermission;
import com.yanzhenjie.permission.PermissionListener;
import com.example.meerkats.adapter.FileHolder;
import com.example.meerkats.adapter.FileAdapter;
import com.example.meerkats.adapter.TitleAdapter;
import com.example.meerkats.adapter.base.RecyclerViewAdapter;
import com.example.meerkats.bean.FileBean;
import com.example.meerkats.bean.TitlePath;
import com.example.meerkats.bean.FileType;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;


public class MainActivity extends AppCompatActivity {
    private RecyclerView title_recycler_view;
    private RecyclerView recyclerView;
    private FileAdapter fileAdapter;
    private List<FileBean> beanList = new ArrayList<>();
    private File rootFile;
    private LinearLayout empty_rel;
    private int PERMISSION_CODE_WRITE_EXTERNAL_STORAGE = 100;
    private String rootPath;
    private TitleAdapter titleAdapter;
    private FloatingActionButton ButtonSync;
    private static Gson gson = new Gson();
    private ExecutorService threadPool;
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


        /*while (a<100) {
            Toast.makeText(MainActivity.this, "Successful Synchronization", Toast.LENGTH_SHORT).show();
            Timer timer = new Timer();
            timer.schedule(new TimerTask() {

                public void run() {
                    a++;
                    System.out.print(a);
                }
            }, 0, 1000);

        }*/


        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        File file = new File(getExternalFilesDir("Meerkats").getPath());
        file.mkdir();
        String rootPath = file.getPath();
        final byte[] messageBodyByte = func(file, tcpMeerkats);


        FloatingActionButton ButtonSync = (FloatingActionButton) findViewById(R.id.buttonsync);

        ButtonSync.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                threadPool = Executors.newCachedThreadPool();
                threadPool.execute(new Runnable() {
                    @Override
                    public void run() {
                        File file = new File(getExternalFilesDir("Meerkats").getPath());
                        file.mkdir();
                        String rootPath = file.getPath();
                        byte[] messageBodyByte = func(file, tcpMeerkats);

                        tcpMeerkats.createInstance();
                        tcpMeerkats.connectSocket();

                        tcpMeerkats.sendMessage(messageBodyByte);


                        tcpMeerkats.checkCMDFlag(tcpMeerkats.receiveMessage());

                    }
                });


                Toast.makeText(MainActivity.this, "Successful Synchronization", Toast.LENGTH_SHORT).show();
            }
        });


        //Set Title
        title_recycler_view = (RecyclerView) findViewById(R.id.title_recycler_view);
        title_recycler_view.setLayoutManager(new LinearLayoutManager(this, LinearLayoutManager.HORIZONTAL, false));
        titleAdapter = new TitleAdapter(MainActivity.this, new ArrayList<TitlePath>());
        title_recycler_view.setAdapter(titleAdapter);

        recyclerView = (RecyclerView) findViewById(R.id.recycler_view);

        fileAdapter = new FileAdapter(this, beanList);
        recyclerView.setLayoutManager(new LinearLayoutManager(this));
        recyclerView.setAdapter(fileAdapter);

        empty_rel = (LinearLayout) findViewById(R.id.empty_rel);

        fileAdapter.setOnItemClickListener(new RecyclerViewAdapter.OnItemClickListener() {
            @Override
            public void onItemClick(View view, RecyclerView.ViewHolder viewHolder, int position) {
                if (viewHolder instanceof FileHolder) {
                    FileBean file = beanList.get(position);
                    FileType fileType = file.getFileType();
                    if (fileType == FileType.directory) {
                        getFile(file.getPath());

                        refreshTitleState(file.getName(), file.getPath());
                    } else if (fileType == FileType.apk) {
                        //Install app
                        com.example.meerkats.FileUtil.openAppIntent(MainActivity.this, new File(file.getPath()));
                    } else if (fileType == FileType.image) {
                        com.example.meerkats.FileUtil.openImageIntent(MainActivity.this, new File(file.getPath()));
                    } else if (fileType == FileType.txt) {
                        com.example.meerkats.FileUtil.openTextIntent(MainActivity.this, new File(file.getPath()));
                    } else if (fileType == FileType.music) {
                        com.example.meerkats.FileUtil.openMusicIntent(MainActivity.this, new File(file.getPath()));
                    } else if (fileType == FileType.video) {
                        com.example.meerkats.FileUtil.openVideoIntent(MainActivity.this, new File(file.getPath()));
                    } else {
                        com.example.meerkats.FileUtil.openApplicationIntent(MainActivity.this, new File(file.getPath()));
                    }
                }
            }
        });

        fileAdapter.setOnItemLongClickListener(new RecyclerViewAdapter.OnItemLongClickListener() {
            @Override
            public boolean onItemLongClick(View view, RecyclerView.ViewHolder viewHolder, int position) {
                if (viewHolder instanceof FileHolder) {
                    FileBean fileBean = (FileBean) fileAdapter.getItem(position);
                    FileType fileType = fileBean.getFileType();
                    if (fileType != null && fileType != FileType.directory) {
                        com.example.meerkats.FileUtil.sendFile(MainActivity.this, new File(fileBean.getPath()));
                    }
                }
                return false;
            }
        });

        titleAdapter.setOnItemClickListener(new RecyclerViewAdapter.OnItemClickListener() {
            @Override
            public void onItemClick(View view, RecyclerView.ViewHolder viewHolder, int position) {
                TitlePath titlePath = (TitlePath) titleAdapter.getItem(position);
                getFile(titlePath.getPath());

                int count = titleAdapter.getItemCount();
                int removeCount = count - position - 1;
                for (int i = 0; i < removeCount; i++) {
                    titleAdapter.removeLast();
                }
            }
        });



        refreshTitleState("Your Meerkat", "/storage/emulated/0/Android/data/com.example.meerkats/files/Meerkats/");

        // Judge whether there is permission first
        if (AndPermission.hasPermission(this, Manifest.permission.WRITE_EXTERNAL_STORAGE)) {
            // Right to do anything directly
            getFile(rootPath);
        } else {
            //Application authority
            AndPermission.with(this)
                    .requestCode(PERMISSION_CODE_WRITE_EXTERNAL_STORAGE)
                    .permission(Manifest.permission.WRITE_EXTERNAL_STORAGE)
                    .send();
        }
    }

    public void getFile(String path) {
        rootFile = new File(path + File.separator);
        new MyTask(rootFile).executeOnExecutor(AsyncTask.THREAD_POOL_EXECUTOR, "");
    }

    class MyTask extends AsyncTask {
        File file;

        MyTask(File file) {
            this.file = file;
        }

        @Override
        protected Object doInBackground(Object[] params) {
            List<FileBean> fileBeenList = new ArrayList<>();
            if (file.isDirectory()) {
                File[] filesArray = file.listFiles();
                if (filesArray != null) {
                    List<File> fileList = new ArrayList<>();
                    Collections.addAll(fileList, filesArray);  //Convert an array to a list
                    Collections.sort(fileList, com.example.meerkats.FileUtil.comparator);  //Sort by name

                    for (File f : fileList) {
                        if (f.isHidden()) continue;

                        FileBean fileBean = new FileBean();
                        fileBean.setName(f.getName());
                        fileBean.setPath(f.getAbsolutePath());
                        fileBean.setFileType(com.example.meerkats.FileUtil.getFileType(f));
                        fileBean.setChildCount(com.example.meerkats.FileUtil.getFileChildCount(f));
                        fileBean.setSize(f.length());
                        fileBean.setHolderType(0);

                        fileBeenList.add(fileBean);

                        FileBean lineBean = new FileBean();
                        lineBean.setHolderType(1);
                        fileBeenList.add(lineBean);

                    }
                }
            }

            beanList = fileBeenList;
            return fileBeenList;
        }

        @Override
        protected void onPostExecute(Object o) {
            if (beanList.size() > 0) {
                empty_rel.setVisibility(View.GONE);
            } else {
                empty_rel.setVisibility(View.VISIBLE);
            }
            fileAdapter.refresh(beanList);
        }
    }

    void refreshTitleState(String title, String path) {
        TitlePath filePath = new TitlePath();
        filePath.setNameState(title + " > ");
        filePath.setPath(path);
        titleAdapter.addItem(filePath);
        title_recycler_view.smoothScrollToPosition(titleAdapter.getItemCount());
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK
                && event.getRepeatCount() == 0) {

            List<TitlePath> titlePathList = (List<TitlePath>) titleAdapter.getAdapterData();
            if (titlePathList.size() == 1) {
                finish();
            } else {
                titleAdapter.removeItem(titlePathList.size() - 1);
                getFile(titlePathList.get(titlePathList.size() - 1).getPath());
            }
            return true;
        }
        return super.onKeyDown(keyCode, event);
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        // Just call this sentence and pass the rest to AndPermission. The last parameter is PermissionListener
        AndPermission.onRequestPermissionsResult(requestCode, permissions, grantResults, listener);
    }

    private PermissionListener listener = new PermissionListener() {
        @Override
        public void onSucceed(int requestCode, List<String> grantedPermissions) {
            // Successful callback of privilege application
            if (requestCode == PERMISSION_CODE_WRITE_EXTERNAL_STORAGE) {
                getFile(rootPath);
            }
        }

        @Override
        public void onFailed(int requestCode, List<String> deniedPermissions) {
            // Callback for Failure of Privilege Application
            AndPermission.defaultSettingDialog(MainActivity.this, PERMISSION_CODE_WRITE_EXTERNAL_STORAGE)
                    .setTitle("Failure of permission application")
                    .setMessage("Some of the permissions we need have been rejected by you or the system failed to apply for errors. Please go to the settings page to authorize manually, otherwise the function will not work properly!")
                    .setPositiveButton("Okay, go ahead and set it up")
                    .show();
        }
    };
}
