#include "common.h"

#include <spectre_def.h>

#include <debug.h>
#include <string.h>

#ifndef _WIN32

#include <ifaddrs.h>
#include <arpa/inet.h>
#include <netinet/tcp.h>
#include <bits/types/struct_timeval.h>

#else

#include <windows.h>
#include <winsock2.h>
#include <ws2tcpip.h>

#endif


char *addr_to_string(void *src) {
    if (src == NULL) return NULL;

    char ipv6_buf[INET6_ADDRSTRLEN + 1];
    if (inet_ntop(AF_INET6, src, ipv6_buf, sizeof(ipv6_buf)) != NULL) {
        return strdup(ipv6_buf);
    } else {
        char ipv4_buf[INET_ADDRSTRLEN + 1];
        if (inet_ntop(AF_INET, src, ipv4_buf, sizeof(ipv4_buf)) != NULL) {
            return strdup(ipv4_buf);
        }
    }

    return NULL;
}


bool string_to_addr(int af, const char *src, void *dst) {
#ifdef _WIN32
    struct sockaddr_storage ss;
    int len = sizeof(ss);
    char buf[INET6_ADDRSTRLEN + 1];
    bool ret = 0;
    memset(&ss, 0, sizeof(ss));
    memset(buf, 0, sizeof(buf));

    if (WSAStringToAddressA((char *) src, af, NULL, (struct sockaddr *) &ss, &len) == 0) {
        if (af == AF_INET) {
            *(struct in_addr *) dst = ((struct sockaddr_in *) &ss)->sin_addr;
            ret = 1;
        } else if (af == AF_INET6) {
            *(struct in6_addr *) dst = ((struct sockaddr_in6 *) &ss)->sin6_addr;
            ret = 1;
        }
    }
    return ret ;
#else
    return inet_pton(af, src, dst) > 0;
#endif
}

const int enable = 1;

CONNECTION open_raw_socket(const int protocol, bool is_v6) {
#ifndef _WIN32
    const CONNECTION sock = socket(is_v6 ? AF_INET6 : AF_INET, SOCK_RAW, protocol);
    if (sock == INVALID_SOCKET_T) {
        debug_error("socket");
        return INVALID_SOCKET_T;
    }

    const int level = is_v6 ? IPPROTO_IPV6 : IPPROTO_IP;
    const int optname = is_v6 ? IPV6_HDRINCL : IP_HDRINCL;

    if (setsockopt(sock, level, optname, &enable, sizeof(enable)) < 0) {
        CLOSE_SOCKET(sock);
        return INVALID_SOCKET_T;
    }

    if (setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, &enable, sizeof(enable)) < 0) {
        CLOSE_SOCKET(sock);
        return INVALID_SOCKET_T;
    }

    return sock;
#else
    return INVALID_SOCKET_T;
#endif
}

CONNECTION open_socket(int af, int protocol) {
    CONNECTION sock = INVALID_SOCKET_T;

    if ((sock = socket(af, protocol, 0)) == INVALID_SOCKET_T) {
        debug_printf("cannot open socket");
        return INVALID_SOCKET_T;
    }

    return sock;
}


CONNECTION dial_with_timeout(struct sockaddr *addr, int addr_size, const int timeoutSec) {
    CONNECTION sock = INVALID_SOCKET_T;

    if ((sock = open_socket(addr->sa_family, SOCK_STREAM)) == INVALID_SOCKET_T) {
        debug_printf("cannot open socket");
        return INVALID_SOCKET_T;
    }

    struct timeval timeout = {
            .tv_sec = timeoutSec,
            .tv_usec = 0,
    };

    if (setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, (const char *) &timeout, sizeof(timeout)) != 0) {
        debug_printf("cannot set socket SO_RCVTIMEO");
        return INVALID_SOCKET_T;
    }

    if (setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, (const char *) &timeout, sizeof(timeout)) != 0) {
        debug_printf("cannot set socket SO_SNDTIMEO");
        return INVALID_SOCKET_T;
    }

    setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, (const char *) &enable, sizeof(enable));
    setsockopt(sock, IPPROTO_TCP, TCP_NODELAY, (const char *) &enable, sizeof(enable));

    if (connect(sock, addr, addr_size) == SOCKET_ERROR_T) {
        debug_printf("cannot connect socket");
        return INVALID_SOCKET_T;
    }

    return sock;
}