# log

Logging facilities for applications.

## Functions

### write

**Parameters:**

- `level`: `i32`
- `text`: `&[u8]`

**Returns:** `none`

**Semantics:**

Writes `text` to the environment-provided logger.

The text must be valid UTF-8 or otherwise the behavior is implementation-defined.

`level` can be one of:

- 1: Error
- 3: Warning
- 6: Info
