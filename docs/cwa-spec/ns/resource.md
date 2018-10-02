# resource

Operations related to readable/writeable resources.

## Functions

### open

**Parameters:**

- `url`: `&[u8]`

**Returns:** `i32`

**Semantics:**

Opens the resource specified by `url`.

Returns the resource id on success, or the error code on failure.

The URL format is as defined in `urls-and-schemes.md`.

### read

**Parameters:**

- `id`: `i32`
- `data`: `&mut [u8]`

**Returns:** `i32`

**Semantics:**

Reads from the resource specified by `id`, into `data`.

Returns the actual bytes read on success, or the error on failure.

### write

**Parameters:**

- `id`: `i32`
- `data`: `&[u8]`

**Returns:** `i32`

**Semantics:**

Writes `data` to the resource specified by `id`.

Returns the actual bytes written on success, or the error on failure.

### close

**Parameters:**

- `id`: `i32`

**Returns:** `none`

**Semantics:**

Closes the resource specified by `id`.

Calling this on an invalid resource id is a fatal error.

### flush

**Parameters:**

- `id`: `i32`

**Returns:** `i32`

**Semantics:**

Instructs the resource specified by `id` to flush the output stream, ensuring that all intermediately buffered contents reach their destination.

Returns `0` on success, or the error on failure.
