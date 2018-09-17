# startup

Obtain startup information about the current application.

## Functions

### arg_len

**Parameters:**

(None)

**Returns:** `i32`

**Semantics:**

Returns the number of arguments passed to the current application, including the path to application binary itself.

### arg_at

**Parameters:**

- `id`: `i32`
- `out`: `&mut [u8]`

**Returns:** `i32`

**Semantics:**

Writes the argument at position `id` to `out`.

The result is truncated if the argument length is greater than that of `out`.

Returns the number of bytes written, or `InvalidArgumentError` if `id` is greater than or equal to the total number of arguments.
