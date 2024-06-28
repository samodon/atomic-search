---
Language:
  - C
date: 2024-06-18
Topics:
  - Memory Management
  - C
tags:
  - languages/C
---
The calloc function in C allows for [[Dynamic Memory Management]]. It allocates an array of elements and initializes all of the bytes to zero.
```C
void *calloc(size_t num, size_t size);
```
**Example**:
```C
int *ptr = (int *)calloc(10, sizeof(int));  // Allocates memory for 10 integers and initializes them to zero
```