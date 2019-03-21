# io

File descriptors for standard input and output.

## Functions

### io\_get\_stdin

**Parameters:**

None

**Returns:** `i32` file descriptor

**Semantics:**

Returns a read-only file descriptor pointing to the semantic standard input of
the WebAssembly VM.

### io\_get\_stdout

**Parameters:**

None

**Returns:** `i32` file descriptor

**Semantics:**

Returns a write-only file descriptor pointing to the semantic standard output of
the WebAssembly VM.

### io\_get\_stderr

**Parameters:**

None

**Returns:** `i32` file descriptor

**Semantics:**

Returns a write-only file descriptor pointing to the semantic standard error output of
the WebAssembly VM.
