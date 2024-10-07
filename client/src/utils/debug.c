#include "debug.h"

#if defined(DEBUG) || defined(DEBUG_LOG)

#include <stdarg.h>
#include <stdio.h>
#include <time.h>
#include <ctype.h>

#endif


void debug_time() {
    char timeText[32];
    time_t t = time(NULL);
    strftime(timeText, sizeof(timeText), "%Y-%m-%d %H:%M:%S", localtime(&t));
    printf("[%s] ", timeText);
}

void debug_printf(const char *format, ...) {
#if defined(DEBUG) || defined(DEBUG_LOG)
    debug_time();
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
    printf("\r\n");
#endif
}


void debug_bytes(const BYTES *data, const size_t length) {
#if defined(DEBUG) || defined(DEBUG_LOG)
    printf("b\"");
    for (size_t i = 0; i < length; i++) {
        if (data[i] == '\n') {
            printf("\\n");
        } else if (data[i] == '\r') {
            printf("\\r");
        } else if (data[i] == '\t') {
            printf("\\t");
        } else if (isprint(data[i])) {
            printf("%c", data[i]);
        } else {
            printf("\\x%02x", data[i]);
        }
    }
    printf("\"");

#endif
}

void debug_error(const char *data) {
#if defined(DEBUG) || defined(DEBUG_LOG)
    perror(data);
#endif
}
