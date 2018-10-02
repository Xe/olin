# log

This scheme offers a simple write-only logger. The logging data passed through this scheme can be handled in an implementation-dependent manner.

If provided, `?prefix=` should automatically pretend every log line with a given prefix, followed by a colon, and a space. EG:

```
const char* furl = "log://?prefix=test";
int32 fd = resource_open(furl, 18);
const char* msg = "hello, world!"
int32 res = resource_write(fd, msg, 13);
```

Should produce output including:

```
test: hello, world!
```

With anything before or after the given text being implementation-dependent.

## Example URL

```
log://
log://?prefix=test
```

## Behaviors

### [Read](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#read)

Any calls to `read` should succeed, but result in 0 bytes being read. Implementations are suggested to log this happening.

### [Write](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#write)

Data passed to the `write` function will be supplied to an implementation-specific logger stack. If supplied in the URL, the prefix is prepended to every line as described in the top of the file.

### [Close](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#close)

The `close` function on this scheme should clean up any underlying resources and forcibly synchronize any asynchronous operations. 