# null

All bytes written to this scheme are discarded. No bytes are read from this scheme.

## Example URL

```
null://
```

## Behaviors

### [Read](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#read)

No data will be read. The call will succeed, but no bytes will be read.

### [Write](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#write)

No data will be written. The call will succeed, but no bytes will be written.

### [Close](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#close)

No side effects.
