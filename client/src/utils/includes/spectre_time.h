#pragma once


#include <stdint.h>
#include <stdbool.h>
#include <time.h>

typedef struct {
    time_t start;
    time_t end;
} SpectreTimer;


int64_t get_sys_milliseconds();

SpectreTimer *start_timer(long seconds);

void stop_timer(SpectreTimer *timer);

bool is_timer_done(SpectreTimer *timer);