# env

Get configuration information (like database credentials) from the host.

## Functions

### get

**Parameters:**

- `key`: `&[u8]`
- `value`: `&mut [u8]`

**Returns:** `i32`

**Semantics:**

Requests a value from the environment, like environment variables for
common applications. `key` indicate the buffer
containing the key for the information we want to request.
`value` indicate the buffer in which the value
will be written.

Returns:
- `NotFoundError` if the key does not exist
- the number of bytes written to the value buffer if the key exists and the buffer is big enough
- the needed size for the value buffer if it was not big enough. No data is written, and the caller must check if that return value is larger than `value.len` to reallocate and retry.
