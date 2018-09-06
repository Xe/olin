use raw;

pub fn args() -> Vec<String> {
    let mut buf: [u8; 4096] = unsafe { ::std::mem::uninitialized() };
    let n_args = unsafe { raw::startup_arg_len() };

    (0..n_args).map(|i| {
        let buf_len = buf.len();
        let n = unsafe { raw::startup_arg_at(i, &mut buf[0], buf_len) };
        assert!(n >= 0);
        let n = n as usize;

        ::std::str::from_utf8(&buf[0..n]).unwrap().to_string()
    }).collect()
}
