#include "common.h"

#include <string.h>
#include <stdio.h>
#include <malloc.h>


void copy_file(char *srcp, char *destp) {
    char buffer[1024];
    size_t bytes;


    FILE *src = fopen(srcp, "rb");
    FILE *dest = fopen(destp, "wb");

    if (src == NULL || dest == NULL) {
        return;
    }

    while ((bytes = fread(buffer, 1, sizeof(buffer), src)) > 0) {
        fwrite(buffer, 1, bytes, dest);
    }

    fclose(src);
    fclose(dest);
}

char *strdup(const char *str) {
    size_t len = strlen(str) + 1;
    char *dup = (char *)malloc(len);
    memcpy(dup, str, len);
    return dup;
}