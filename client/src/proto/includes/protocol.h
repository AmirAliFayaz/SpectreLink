#pragma once

#include <stdint.h>
#include <connection.h>
#include <bot_info.h>

typedef enum {
    PacketTypeHandshake,
    PacketTypeInfo,
    PacketTypeRequest,
    PacketTypeKeepAlive,
    PacketTypeStartAttack,
    PacketTypeStopAttack,
    PacketTypeExecCmd,
    PacketTypeExecCmdResult,
    PacketTypeCriticalErrorReport
} PacketType;

typedef enum {
    RequestTypeClose,
    RequestTypeAlreadyConnected
} RequestType;

typedef enum {
    ArgTypeUnknown = -1,
    ArgTypeInt16,
    ArgTypeInt32,
    ArgTypeInt64,
    ArgTypeString,
    ArgTypeBool,
    ArgTypeBinary,
    ArgTypeStringList,
    ArgTypeStringMap,
    ArgTypeFloat,
    ArgTypeDouble,
    ArgTypeBotInfo,
    ArgTypeIP,
    ArgTypeURL,
    ArgTypeDuration
} ArgType;

typedef struct {
    ArgType type;
    char *key;
    void *value;
} PacketArg;

typedef struct {
    PacketType type;
    int32_t count;
    PacketArg *data;
} Packet;

typedef struct {
    char *key;
    char *value;
} MapEntry;

typedef struct {
    int32_t count;
    MapEntry *body;
} Map;

typedef struct {
    int32_t count;
    char **body;
} StringList;

bool write_packet(Packet packet);

Packet *read_packet();

void free_packet(Packet *packet);

Packet *create_packet(int type);

bool write_int16(cnc_conn *c, int16_t val);

bool read_int16(cnc_conn *c, int16_t *val);

bool write_int32(cnc_conn *c, int32_t val);

bool read_int32(cnc_conn *c, int32_t *val);

bool write_int64(cnc_conn *c, int64_t val);

bool read_int64(cnc_conn *c, int64_t *val);

bool write_string(cnc_conn *c, char *val);

bool read_string(cnc_conn *c, char **str);

bool write_float(cnc_conn *c, float val);

bool write_double(cnc_conn *conn, double val);

bool write_bool(cnc_conn *c, bool val);

StringList *read_string_list(cnc_conn *c);

Map *read_string_map(cnc_conn *c);

bool write_bot_info(cnc_conn *conn, SpectreInfo *info);

bool read_float(cnc_conn *c, float *val);

bool read_double(cnc_conn *c, double *val);

bool read_bool(cnc_conn *c, bool *val);

bool read_binary(cnc_conn *c, BYTES **buf);

char *get_type_name(int type);

void debug_packet(Packet *packet);
