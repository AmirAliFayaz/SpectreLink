#pragma once

#include <spectre_def.h>

typedef struct {
    char *Username;
    char *OS;
    char *Kernel;
    char *Arch;
    char *Version;
    char *InfectionMethod;
    int Processors;
    int UpTime;
    float TotalMemory;
    float FreeMemory;
    double TimeZoneDiff;
    int64_t SystemTime;
    bool IsRoot;
    bool LittleEndian;
    bool Is64Bit;
    bool IsDebugMode;
    bool IsIPv6Supported;
    bool HasAnySSLLib;
    bool FirewallStatus;
} SpectreInfo;

SpectreInfo *get_bot_info();

void init_bot_info();

void set_infection_method(char *string);

void print_info();