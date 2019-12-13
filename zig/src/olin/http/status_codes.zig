/// StatusCode is a HTTP status code. This matches the list at https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
pub const StatusCode = enum(u32) {
    /// Continue indicates that the initial part of a request has been received and has not yet been rejected by the server.
    Continue = 100,

    /// SwitchingProtocols indicates that the server understands and is willing to comply with the client's request, via the Upgrade header field, for a change in the application protocol being used on this connection.
    SwitchingProtocols = 101,

    /// OK indicates that the request has succeeded.
    OK = 200,

    /// Created indicates that the request has been fulfilled and has resulted in one or more new resources being created.
    Created = 201,

    /// Accepted indicates that the request has been accepted for processing, but the processing has not been completed.
    Accepted = 202,

    /// NonAuthoritativeInformation indicates that the request was successful but the enclosed payload has been modified from that of the origin server's 200 (OK) response by a transforming proxy.
    NonAuthoritativeInformation = 203,

    /// NoContent indicates that the server has successfully fulfilled the request and that there is no additional content to send in the response payload body.
    NoContent = 204,

    /// ResetContent indicates that the server has fulfilled the request and desires that the user agent reset the "document view", which caused the request to be sent, to its original state as received from the origin server.
    ResetContent = 205,

    /// PartialContent indicates that the server is successfully fulfilling a range request for the target resource by transferring one or more parts of the selected representation that correspond to the satisfiable ranges found in the requests's Range header field.
    PartialContent = 206,

    /// MultipleChoices indicates that the target resource has more than one representation, each with its own more specific identifier, and information about the alternatives is being provided so that the user (or user agent) can select a preferred representation by redirecting its request to one or more of those identifiers.
    MultipleChoices = 300,

    /// MovedPermanently indicates that the target resource has been assigned a new permanent URI and any future references to this resource ought to use one of the enclosed URIs.
    MovedPermanently = 301,

    /// Found indicates that the target resource resides temporarily under a different URI.
    Found = 302,

    /// SeeOther indicates that the server is redirecting the user agent to a different resource, as indicated by a URI in the Location header field, that is intended to provide an indirect response to the original request.
    SeeOther = 303,

    /// NotModified indicates that a conditional GET request has been received and would have resulted in a 200 (OK) response if it were not for the fact that the condition has evaluated to false.
    NotModified = 304,

    /// UseProxy deprecated
    UseProxy = 305,

    /// TemporaryRedirect indicates that the target resource resides temporarily under a different URI and the user agent MUST NOT change the request method if it performs an automatic redirection to that URI.
    TemporaryRedirect = 307,

    /// BadRequest indicates that the server cannot or will not process the request because the received syntax is invalid, nonsensical, or exceeds some limitation on what the server is willing to process.
    BadRequest = 400,

    /// Unauthorized indicates that the request has not been applied because it lacks valid authentication credentials for the target resource.
    Unauthorized = 401,

    /// PaymentRequired reserved
    PaymentRequired = 402,

    /// Forbidden indicates that the server understood the request but refuses to authorize it.
    Forbidden = 403,

    /// NotFound indicates that the origin server did not find a current representation for the target resource or is not willing to disclose that one exists.
    NotFound = 404,

    /// MethodNotAllowed indicates that the method specified in the request-line is known by the origin server but not supported by the target resource.
    MethodNotAllowed = 405,

    /// NotAcceptable indicates that the target resource does not have a current representation that would be acceptable to the user agent, according to the proactive negotiation header fields received in the request, and the server is unwilling to supply a default representation.
    NotAcceptable = 406,

    /// ProxyAuthenticationRequired is similar to 401 (Unauthorized), but indicates that the client needs to authenticate itself in order to use a proxy.
    ProxyAuthenticationRequired = 407,

    /// RequestTimeout indicates that the server did not receive a complete request message within the time that it was prepared to wait.
    RequestTimeout = 408,

    /// Conflict indicates that the request could not be completed due to a conflict with the current state of the resource.
    Conflict = 409,

    /// Gone indicates that access to the target resource is no longer available at the origin server and that this condition is likely to be permanent.
    Gone = 410,

    /// LengthRequired indicates that the server refuses to accept the request without a defined Content-Length.
    LengthRequired = 411,

    /// PreconditionFailed indicates that one or more preconditions given in the request header fields evaluated to false when tested on the server.
    PreconditionFailed = 412,

    /// PayloadTooLarge indicates that the server is refusing to process a request because the request payload is larger than the server is willing or able to process.
    PayloadTooLarge = 413,

    /// URITooLong indicates that the server is refusing to service the request because the request-target is longer than the server is willing to interpret.
    URITooLong = 414,

    /// UnsupportedMediaType indicates that the origin server is refusing to service the request because the payload is in a format not supported by the target resource for this method.
    UnsupportedMediaType = 415,

    /// RangeNotSatisfiable indicates that none of the ranges in the request's Range header field overlap the current extent of the selected resource or that the set of ranges requested has been rejected due to invalid ranges or an excessive request of small or overlapping ranges.
    RangeNotSatisfiable = 416,

    /// ExpectationFailed indicates that the expectation given in the request's Expect header field could not be met by at least one of the inbound servers.
    ExpectationFailed = 417,

    /// Imateapot Any attempt to brew coffee with a teapot should result in the error code 418 I'm a teapot.
    Imateapot = 418,

    /// UpgradeRequired indicates that the server refuses to perform the request using the current protocol but might be willing to do so after the client upgrades to a different protocol.
    UpgradeRequired = 426,

    /// InternalServerError indicates that the server encountered an unexpected condition that prevented it from fulfilling the request.
    InternalServerError = 500,

    /// NotImplemented indicates that the server does not support the functionality required to fulfill the request.
    NotImplemented = 501,

    /// BadGateway indicates that the server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.
    BadGateway = 502,

    /// ServiceUnavailable indicates that the server is currently unable to handle the request due to a temporary overload or scheduled maintenance, which will likely be alleviated after some delay.
    ServiceUnavailable = 503,

    /// GatewayTimeout indicates that the server, while acting as a gateway or proxy, did not receive a timely response from an upstream server it needed to access in order to complete the request.
    GatewayTimeout = 504,

    /// HTTPVersionNotSupported indicates that the server does not support, or refuses to support, the protocol version that was used in the request message.
    HTTPVersionNotSupported = 505,

    /// Processing is an interim response used to inform the client that the server has accepted the complete request, but has not yet completed it.
    Processing = 102,

    /// MultiStatus provides status for multiple independent operations.
    MultiStatus = 207,

    /// IMUsed The server has fulfilled a GET request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance.
    IMUsed = 226,

    /// PermanentRedirect The target resource has been assigned a new permanent URI and any future references to this resource outght to use one of the enclosed URIs. [...] This status code is similar to 301 Moved Permanently (Section 7.3.2 of rfc7231), except that it does not allow rewriting the request method from POST to GET.
    PermanentRedirect = 308,

    /// UnprocessableEntity means the server understands the content type of the request entity (hence a 415(Unsupported Media Type) status code is inappropriate), and the syntax of the request entity is correct (thus a 400 (Bad Request) status code is inappropriate) but was unable to process the contained instructions.
    UnprocessableEntity = 422,

    /// Locked means the source or destination resource of a method is locked.
    Locked = 423,

    /// FailedDependency means that the method could not be performed on the resource because the requested action depended on another action and that action failed.
    FailedDependency = 424,

    /// PreconditionRequired indicates that the origin server requires the request to be conditional.
    PreconditionRequired = 428,

    /// TooManyRequests indicates that the user has sent too many requests in a given amount of time ("rate limiting").
    TooManyRequests = 429,

    /// RequestHeaderFieldsTooLarge indicates that the server is unwilling to process the request because its header fields are too large.
    RequestHeaderFieldsTooLarge = 431,

    /// UnavailableForLegalReasons This status code indicates that the server is denying access to the resource in response to a legal demand.
    UnavailableForLegalReasons = 451,

    /// VariantAlsoNegotiates indicates that the server has an internal configuration error: the chosen variant resource is configured to engage in transparent content negotiation itself, and is therefore not a proper end point in the negotiation process.
    VariantAlsoNegotiates = 506,

    /// InsufficientStorage means the method could not be performed on the resource because the server is unable to store the representation needed to successfully complete the request.
    InsufficientStorage = 507,

    /// NetworkAuthenticationRequired indicates that the client needs to authenticate to gain network access.
    NetworkAuthenticationRequired = 511,
};

pub fn reasonPhrase(sc: StatusCode) []const u8 {
    return switch (sc) {
        StatusCode.Continue => "Continue"[0..],
        StatusCode.SwitchingProtocols => "Switching Protocols"[0..],
        StatusCode.OK => "OK"[0..],
        StatusCode.Created => "Created"[0..],
        StatusCode.Accepted => "Accepted"[0..],
        StatusCode.NonAuthoritativeInformation => "Non-Authoritative Information"[0..],
        StatusCode.NoContent => "No Content"[0..],
        StatusCode.ResetContent => "Reset Content"[0..],
        StatusCode.PartialContent => "Partial Content"[0..],
        StatusCode.MultipleChoices => "Multiple Choices"[0..],
        StatusCode.MovedPermanently => "Moved Permanently"[0..],
        StatusCode.Found => "Found"[0..],
        StatusCode.SeeOther => "See Other"[0..],
        StatusCode.NotModified => "Not Modified"[0..],
        StatusCode.UseProxy => "Use Proxy"[0..],
        StatusCode.TemporaryRedirect => "Temporary Redirect"[0..],
        StatusCode.BadRequest => "Bad Request"[0..],
        StatusCode.Unauthorized => "Unauthorized"[0..],
        StatusCode.PaymentRequired => "Payment Required"[0..],
        StatusCode.Forbidden => "Forbidden"[0..],
        StatusCode.NotFound => "Not Found"[0..],
        StatusCode.MethodNotAllowed => "Method Not Allowed"[0..],
        StatusCode.NotAcceptable => "Not Acceptable"[0..],
        StatusCode.ProxyAuthenticationRequired => "Proxy Authentication Required"[0..],
        StatusCode.RequestTimeout => "Request Timeout"[0..],
        StatusCode.Conflict => "Conflict"[0..],
        StatusCode.Gone => "Gone"[0..],
        StatusCode.LengthRequired => "Length Required"[0..],
        StatusCode.PreconditionFailed => "Precondition Failed"[0..],
        StatusCode.PayloadTooLarge => "Payload Too Large"[0..],
        StatusCode.URITooLong => "URI Too Long"[0..],
        StatusCode.UnsupportedMediaType => "Unsupported Media Type"[0..],
        StatusCode.RangeNotSatisfiable => "Range Not Satisfiable"[0..],
        StatusCode.ExpectationFailed => "Expectation Failed"[0..],
        StatusCode.Imateapot => "I'm a teapot"[0..],
        StatusCode.UpgradeRequired => "Upgrade Required"[0..],
        StatusCode.InternalServerError => "Internal Server Error"[0..],
        StatusCode.NotImplemented => "Not Implemented"[0..],
        StatusCode.BadGateway => "Bad Gateway"[0..],
        StatusCode.ServiceUnavailable => "Service Unavailable"[0..],
        StatusCode.GatewayTimeout => "Gateway Time-out"[0..],
        StatusCode.HTTPVersionNotSupported => "HTTP Version Not Supported"[0..],
        StatusCode.Processing => "Processing"[0..],
        StatusCode.MultiStatus => "Multi-Status"[0..],
        StatusCode.IMUsed => "IM Used"[0..],
        StatusCode.PermanentRedirect => "Permanent Redirect"[0..],
        StatusCode.UnprocessableEntity => "Unprocessable Entity"[0..],
        StatusCode.Locked => "Locked"[0..],
        StatusCode.FailedDependency => "Failed Dependency"[0..],
        StatusCode.PreconditionRequired => "Precondition Required"[0..],
        StatusCode.TooManyRequests => "Too Many Requests"[0..],
        StatusCode.RequestHeaderFieldsTooLarge => "Request Header Fields Too Large"[0..],
        StatusCode.UnavailableForLegalReasons => "Unavailable For Legal Reasons"[0..],
        StatusCode.VariantAlsoNegotiates => "Variant Also Negotiates"[0..],
        StatusCode.InsufficientStorage => "Insufficient Storage"[0..],
        StatusCode.NetworkAuthenticationRequired => "Network Authentication Required"[0..],
        else => "Unknown"[0..],
    };
}
