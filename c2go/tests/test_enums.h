#ifndef TEST_ENUMS_H
#define TEST_ENUMS_H

// Simple enum
typedef enum {
    COLOR_RED,
    COLOR_GREEN,
    COLOR_BLUE
} Color;

// Enum with explicit values
typedef enum {
    STATUS_UNKNOWN = -1,
    STATUS_OFF = 0,
    STATUS_ON = 1,
    STATUS_ERROR = 2
} Status;

// Flag enum
typedef enum {
    FLAG_NONE = 0,
    FLAG_READ = 1,
    FLAG_WRITE = 2,
    FLAG_EXECUTE = 4,
    FLAG_ALL = FLAG_READ | FLAG_WRITE | FLAG_EXECUTE
} Permissions;

#endif // TEST_ENUMS_H