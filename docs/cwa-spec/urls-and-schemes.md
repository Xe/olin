# URLs and Schemes

In CommonWA-compatible platforms, all resources like files and sockets should be identified by URLs.

### "Everything is a URL"

Similar to the Redox OS, CommonWA specifies a method to use the URL to uniquely identify a resource.

Schemes, as defined in [RFC1738](https://tools.ietf.org/html/rfc1738), indicates the provider that manages the resource, and can be seen as a "generalization" of file systems.

### Core Schemes

A small number of core schemes are defined here and should be implemented by all CommonWA-compatible platforms:

| Name | Description | Docs
| --- | --- | --- |
| null | The scheme to which all written bytes are discarded and from which nothing are read | [Docs](./scheme/null.md) |
| zero | The scheme to which all written bytes are discarded and from which a infinite stream of `0` are read | [Docs](./scheme/zero.md) |
| log | Logging | [Docs](./scheme/log.md) |

For platforms with specific features, the related feature-dependent schemes should be implemented:

Multi-task platforms:

| Name | Description | Docs
| --- | --- | --- |
| iac | Inter-application communication | (TODO) |

Platforms with access to Unix-like filesystems:

| Name | Description | Docs
| --- | --- | --- |
| file | The local filesystem | (TODO) |

Platforms with networking:

| Name | Description | Docs
| --- | --- | --- |
| tcp | TCP sockets | (TODO) |
| udp | UDP sockets | (TODO) |

### Application-provided Schemes

TODO
