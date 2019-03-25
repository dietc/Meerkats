package com.example.meerkats.adapter;

import android.view.View;
import android.widget.TextView;

import com.example.meerkats.R;
import com.example.meerkats.adapter.base.RecyclerViewAdapter;
import com.example.meerkats.adapter.base.RecyclerViewHolder;
import com.example.meerkats.bean.TitlePath;


public class TitleHolder extends RecyclerViewHolder<TitleHolder> {

    TextView textView ;

    public TitleHolder(View itemView) {
        super(itemView);

        textView = (TextView) itemView.findViewById(R.id.title_Name );
    }

    @Override
    public void onBindViewHolder(TitleHolder lineHolder, RecyclerViewAdapter adapter, int position) {
        TitlePath titlePath = (TitlePath) adapter.getItem( position );
        lineHolder.textView.setText( titlePath.getNameState() );
    }
}
