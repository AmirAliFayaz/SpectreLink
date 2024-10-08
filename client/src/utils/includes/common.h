#pragma once

#include <spectre_def.h>
#include <stdbool.h>

#ifndef _WIN32

#include <arpa/inet.h>

#else
#include <windows.h>
#include <winsock2.h>
#include <ws2tcpip.h>
#endif


void copy_file(char *srcp, char *destp);

char *addr_to_string(void *src);

bool string_to_addr(int af, const char *src, void *dst);

CONNECTION open_raw_socket(int protocol, bool is_v6);

CONNECTION dial_with_timeout(struct sockaddr *addr, int addr_size, int timeoutSec);

CONNECTION open_socket(int af, int protocol);

char *strdup(const char *str);
