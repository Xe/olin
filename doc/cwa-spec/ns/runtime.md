# runtime

Control and obtain information about the runtime environment.

## Functions

### spec_major

**Parameters:**

(None)

**Returns:** `i32`

**Semantics:**

Returns the major version of the CommonWA spec implemented by the runtime environment.

This value must be **equal** to which the application targets to ensure all features working correctly.

### spec_minor

**Parameters:**

(None)

**Returns:** `i32`

**Semantics:**

Returns the minor version of the CommonWA spec implemented by the runtime environment.

This value must be **greater or equal** to which the application targets to ensure all features working correctly.

### name

**Parameters:**

- `out`: `&mut [u8]`

**Returns:** `i32`

**Semantics:**

Writes the name of the current runtime environment to `out` and returns the number of bytes written.

If the length of `out` is less than the length of the name to write, InvalidArgumentError should be returned.

The name must not be longer than 32 bytes and must be valid UTF-8.

### msleep

**Parameters:**

- `ms`: `i32`

**Returns:** `none`

**Semantics:**

Sleeps for the given `ms` milliseconds.

If the environment does not support sleeping, calling this results in a fatal error.
