#ifndef TEST_FUNCTIONS_H
#define TEST_FUNCTIONS_H

// Simple test functions for demonstration
int add(int a, int b);
int subtract(int a, int b);
void print_message(const char* message);
double multiply(double x, double y);

// Function with no parameters
int get_value(void);

// Function with pointer parameter
void process_array(int* array, int length);

#endif // TEST_FUNCTIONS_H