#pragma once

#include <stdbool.h>
#include <stdint.h>

typedef struct {
    char *scheme;
    char *authority;
    int16_t port;
    char *path;
    bool is_ipv6;
} URLComponents;


URLComponents *parse_url(char *url);

char *url_to_string(URLComponents *components);

void free_url(URLComponents *components);