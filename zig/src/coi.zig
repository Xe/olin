const log = @import("./olin/olin.zig").log;

export fn cwa_main() i32 {
    const msg = "coi la munje";
    log.info(msg);
    log.warning(msg);
    log.err(msg);

    return 0;
}
