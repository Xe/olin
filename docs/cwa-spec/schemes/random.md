# random

Crpytographic randomness. The output of this function should be considered safe
to use as the input to cryptographic operations.

## Example URL

```
random://
```

## Behaviors

### [Read](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#read)

Data is read from the system cryptographic entropy source. How this is implemented is up to the runtime.

### [Write](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#write)

Write calls will return `UnknownError`.

### [Close](https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#close)

There is no implication to the close operation.
