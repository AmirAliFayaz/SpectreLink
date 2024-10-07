#pragma once

#include <spectre_def.h>
#include <pthread.h>

typedef struct {
    CONNECTION sock;
    pthread_t keepalive_thread;
    pthread_mutex_t *mx;
    bool connected;
} cnc_conn;

cnc_conn *get_cnc_conn();

void init_cnc_connection();

void cnc_conn_close();

bool cnc_conn_open();

bool is_cnc_connected();