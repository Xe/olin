# zero

In this scheme, all written bytes are ignored, and all bytes read are the null byte `0`.

## Example URL

```
zero://
```

## Behaviors

### [Read](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#read)

Every byte available in the memory given by length will be written to with a `0`.

### [Write](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#write)

No data is written, but all data will report as successfully written.

### [Close](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#close)

No side effects.
