# Types

Apart from the WebAssembly primitive types (i32/i64/f32/f64), we define the following types to be used in API function parameters, which should expand to primitive types in implementations.

### i8/u8

8-bit signed/unsigned integers.

Only used in indirect access to memory (through pointers).

### &T / &mut T

A pointer to one value with type `T`.

Lowers to `i32`.

### &[T] / &mut [T]

A "fat" pointer to a sequence of values of type `T`.

Lowers to `(ptr: i32, len: i32)` where `ptr` is a pointer to the beginning of the sequence and `len` is the number of elements in the sequence (NOT the size in bytes).
