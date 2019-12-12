pub const OlinError = error {
    Unknown,
    InvalidArgument,
    PermissionDenied,
    NotFound,
    EOF,
};

pub fn parse(inp: i32) OlinError!i32 {
    switch (inp) {
        -1 => {
            return error.Unknown;
        },
        -2 => {
            return error.InvalidArgument;
        },
        -3 => {
            return error.PermissionDenied;
        },
        -4 => {
            return error.NotFound;
        },
        -5 => {
            return error.EOF;
        },
        else => {
            return inp;
        },
    }
}
