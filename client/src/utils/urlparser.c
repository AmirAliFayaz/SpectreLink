#include "urlparser.h"

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <common.h>
#include <strings.h>

char *read_until(char **ptr, char end);

URLComponents *parse_url(char *url) {
    URLComponents *components = (URLComponents *) malloc(sizeof(URLComponents));
    if (url == NULL) return NULL;

    components->scheme = NULL;
    components->port = 80;
    components->path = NULL;
    components->authority = NULL;
    components->is_ipv6 = 0;

    char **p = &url;

    if ((components->scheme = read_until(p, ':')) == NULL) {
        free_url(components);
        return NULL;
    }

    if (strcasecmp(components->scheme, "https") == 0) components->port = 443;

    *p += 3; // skip '://'

    if (**p == '[') { // IPv6 address
        (*p)++; // skip '['
        if ((components->authority = read_until(p, ']')) == NULL) {
            free_url(components);
            return NULL;
        }
        (*p)++; // skip ']'

        if (**p == ':') {
            char *port = read_until(p, '/');

            port++;

            if (strcmp(port, "") != 0) {
                components->port = (int16_t) strtol(port, NULL, 10);
            }
        }

        components->is_ipv6 = 1;
    } else { // IPv4 or hostname
        if ((components->authority = read_until(p, '/')) == NULL) {
            free_url(components);
            return NULL;
        }

        if (strchr(components->authority, ':')) {
            char *port = strchr(components->authority, ':') + 1;

            *(strchr(components->authority, ':')) = '\0';

            if (strcmp(port, "") != 0) {
                components->port = (int16_t) strtol(port, NULL, 10);
            }
        }
    }

    if (components->authority == NULL) {
        free_url(components);
        return NULL;
    }

    if (**p == '/') {
        components->path = read_until(p, '\0');
    } else {
        components->path = strdup("/");
    }

    return components;
}

void free_url(URLComponents *components) {
    if (components == NULL) return;
    if (components->scheme != NULL) free(components->scheme);
    if (components->authority != NULL) free(components->authority);
    if (components->path != NULL) free(components->path);
    free(components);
}


char *read_until(char **ptr, char end) {
    char *result = malloc(1);
    if (result == NULL) return NULL;

    size_t length = 0;
    while (**ptr != '\0' && **ptr != end) {
        char *temp;
        if ((temp = realloc(result, length + 2)) == NULL) {
            free(result);
            return NULL;
        }

        result = temp;
        result[length++] = *(*ptr)++;
    }

    result[length] = '\0';
    return result;
}


char *url_to_string(URLComponents *components) {
    size_t scheme_len = strlen(components->scheme);
    size_t authority_len = strlen(components->authority);
    size_t path_len = strlen(components->path);
    size_t port_len = 0;

    if (components->port != 80 && components->port != 443) {
        port_len = snprintf(NULL, 0, "%d", components->port) + 1;
    }

    if (components->is_ipv6) {
        authority_len += 2;
    }

    size_t total_len = scheme_len + 3 + authority_len + path_len + port_len + 1;

    char *url = (char *) malloc(total_len);
    if (url == NULL) return NULL;

    if (components->is_ipv6) {
        if (components->port == 80 || components->port == 443) {
            total_len = snprintf(url, total_len, "%s://[%s]%s", components->scheme, components->authority,
                                 components->path);
        } else {
            total_len = snprintf(url, total_len, "%s://[%s]:%d%s", components->scheme, components->authority,
                                 components->port, components->path);
        }
    } else {
        if (components->port == 80 || components->port == 443) {
            total_len = snprintf(url, total_len, "%s://%s%s", components->scheme, components->authority,
                                 components->path);
        } else {
            total_len = snprintf(url, total_len, "%s://%s:%d%s", components->scheme, components->authority,
                                 components->port, components->path);
        }
    }

    url[total_len] = '\0';

    return url;
}
