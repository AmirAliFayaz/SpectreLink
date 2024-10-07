#pragma once

#include <spectre_def.h>

void debug_printf(const char *format, ...);

void debug_error(const char *data);

void debug_bytes(const BYTES *data, size_t length);

