use raw;

/// Returns the environment variable associated with `key`.
/// If there is no environment variable with the specified key or the value is not valid
/// UTF-8, `None` is returned.
pub fn get(key: &str) -> Option<String> {
    let key = key.as_bytes();
    let mut current_len: usize = 32;
    loop {
        let mut val: Vec<u8> = Vec::with_capacity(current_len);
        let real_len: i32;

        unsafe {
            val.set_len(current_len);
            real_len = raw::env_get(
                slice_raw_ptr_or_null!(key),
                key.len(),
                slice_raw_ptr_or_null_mut!(&mut val),
                val.len()
            );
        }

        if real_len < 0 {
            return None;
        }

        let real_len = real_len as usize;
        if real_len > current_len {
            current_len = real_len;
            continue;
        }

        unsafe {
            val.set_len(real_len);
        }

        return match String::from_utf8(val) {
            Ok(v) => Some(v),
            Err(_) => None
        };
    }
}
