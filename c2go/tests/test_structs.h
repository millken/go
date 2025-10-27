#ifndef TEST_STRUCTS_H
#define TEST_STRUCTS_H

// Simple struct
typedef struct {
    int x;
    int y;
} Point;

// Struct with different types
typedef struct {
    char* name;
    int age;
    double salary;
    void* data;
} Person;

// Struct with array
typedef struct {
    int values[10];
    int count;
} DataBuffer;

#endif // TEST_STRUCTS_H