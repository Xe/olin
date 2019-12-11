// Ported from https://github.com/faithanalog/x/blob/master/zigtest/triangle.zig
const std = @import("std");
const log = @import("./olin/log.zig");
const stdout = @import("./olin/resource.zig").Resource.stdout;

// Dark -> Light color map
const colorMap = " ,:!=+%#$";

// Hardcoded terminal size
const fbWidth = 80;
const fbHeight = 19;

// Framebuffer
var fb: [fbWidth * fbHeight]f32 = undefined;

const Vec2 = struct {
    x: f32,
    y: f32,
};

const Tri = struct {
    v0: Vec2,
    v1: Vec2,
    v2: Vec2,
    c0: f32,
    c1: f32,
    c2: f32,
};

// Clockwise winding
fn edgeFunction(a: Vec2, b: Vec2, c: Vec2) f32 {
    return (c.x - a.x) * (b.y - a.y) - (c.y - a.y) * (b.x - a.x);
}

// Basic unoptimized barycentric rasterizer
fn renderTri(t: Tri) void {
    const area = edgeFunction(t.v0, t.v1, t.v2);

    var y: usize = 0;
    while (y < fbHeight) {
        var x: usize = 0;
        while (x < fbWidth) {
            const p = Vec2{ .x = @intToFloat(f32, x), .y = @intToFloat(f32, y) };
            const w0 = edgeFunction(t.v1, t.v2, p);
            const w1 = edgeFunction(t.v2, t.v0, p);
            const w2 = edgeFunction(t.v0, t.v1, p);
            if (w0 >= 0 and w1 >= 0 and w2 >= 0) {
                const c0 = t.c0 * (w0 / area);
                const c1 = t.c1 * (w1 / area);
                const c2 = t.c2 * (w2 / area);
                fb[(fbHeight - y - 1) * fbWidth + x] = c0 + c1 + c2;
            }
            x = x + 1;
        }
        y = y + 1;
    }
}

fn clearFramebuffer() !void {
    for (fb) |*x, i| {
        x.* = 0;
    }
}

export fn cwa_main() i32 {
    const tri = Tri{
        .v0 = Vec2{
            .x = 75,
            .y = 18,
        },
        .c0 = 1,

        .v1 = Vec2{
            .x = 60,
            .y = 2,
        },
        .c1 = 0.8,

        .v2 = Vec2{
            .x = 4,
            .y = 14,
        },
        .c2 = 0.2,
    };

    renderTri(tri);

    var fbChars: [fbWidth * fbHeight * 10]u8 = undefined;
    var fbNum: usize = 0;

    var buf: [8]u8 = undefined;

    var bright = true;

    // Reset and then set to color 15
    const init = "\x1B[0m\x1B[1;37m";

    // Reset
    const reset = "\x1B[0m";

    // Bold off (dark text)
    const makeDark = "\x1B[21m";

    // Bold on (light text)
    const makeBright = "\x1B[1m";

    for (init) |b, i| fbChars[fbNum + i] = b;
    fbNum = fbNum + init.len;

    for (fb) |*x, pixi| {
        const colzz = std.math.floor(x.* * @intToFloat(f32, colorMap.len * 2));
        const colus = @floatToInt(usize, colzz);
        const colClamped = std.math.min(std.math.max(0, colus), (colorMap.len * 2) - 1);

        if (colClamped < colorMap.len) {
            if (bright) {
                for (makeDark) |b, i| fbChars[fbNum + i] = b;
                fbNum = fbNum + makeDark.len;
                bright = false;
            }
        } else {
            if (!bright) {
                for (makeBright) |b, i| fbChars[fbNum + i] = b;
                fbNum = fbNum + makeBright.len;
                bright = true;
            }
        }

        fbChars[fbNum] = colorMap[colClamped >> 1];
        fbNum = fbNum + 1;

        if (pixi % fbWidth == 0) {
            fbChars[fbNum] = '\n';
            fbNum = fbNum + 1;
        }
    }

    fbChars[fbNum] = '\n';
    fbNum = fbNum + 1;

    for (reset) |b, i| fbChars[fbNum + i] = b;
    fbNum = fbNum + reset.len;

    if(stdout()) |fout| {
        if (fout.write(fbChars[0..])) |n| {
            return 0;
        } else |err| {
            log.err(@errorName(err));
            return 1;
        }
    } else |err| {
        log.err(@errorName(err));
        return 1;
    }

    return 0;
}
