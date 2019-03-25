package com.example.meerkats.adapter.base;

import android.support.v7.widget.RecyclerView;
import android.view.View;



public abstract class RecyclerViewHolder<T> extends RecyclerView.ViewHolder {

    public RecyclerViewHolder(View itemView) {
        super(itemView);
    }

    public abstract void onBindViewHolder (T t , com.example.meerkats.adapter.base.RecyclerViewAdapter adapter , int position) ;
}
