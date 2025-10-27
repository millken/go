#ifndef TEST_COMPLEX_H
#define TEST_COMPLEX_H

#include "test_functions.h"
#include "test_structs.h"
#include "test_enums.h"

// Typedef for basic types
typedef int Counter;
typedef char* String;
typedef void* Handle;

// Complex function using multiple types
typedef struct {
    Point location;
    Color color;
    Status status;
    char name[256];
} GameObject;

// Functions that use the complex types
GameObject* create_game_object(const char* name, Point location, Color color);
void update_game_object(GameObject* obj, Point new_location);
void destroy_game_object(GameObject* obj);

#endif // TEST_COMPLEX_H