# Errors

This document defines all the error codes that may be returned by API functions.

### UnknownError: -1

Represents an error that cannot be represented by any other error code.

Implementations should provide details about this error in logs.

### InvalidArgumentError: -2

One or more arguments passed to the function is invalid.

Any state of the application must not be changed if this error is returned.

### PermissionDeniedError: -3

The application is trying to perform an operation that is not allowed by the security policy.

### NotFoundError: -4

The requested resource is not found.
