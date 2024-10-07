#include "spectre_time.h"

#include <time.h>
#include <malloc.h>
#include <sys/time.h>

int64_t get_sys_milliseconds(void) {
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return tv.tv_sec * 1000 + tv.tv_usec / 1000;
}


SpectreTimer *start_timer(long seconds) {
    SpectreTimer *timer = malloc(sizeof(SpectreTimer));
    timer->start = time(NULL);
    timer->end = timer->start + seconds;
    return timer;
}

void stop_timer(SpectreTimer *timer) {
    timer->end = 0;
}


bool is_timer_done(SpectreTimer *timer) {
    return time(NULL) > timer->end;
}