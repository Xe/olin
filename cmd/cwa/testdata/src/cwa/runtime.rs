use raw;

/// Returns the major part of the spec version implemented by host runtime.
pub fn spec_major() -> i32 {
    unsafe {
        raw::runtime_spec_major()
    }
}

/// Returns the minor part of the spec version implemented by host runtime.
pub fn spec_minor() -> i32 {
    unsafe {
        raw::runtime_spec_minor()
    }
}

/// Returns the name of host runtime.
pub fn name() -> &'static str {
    static mut NAME: Option<&'static str> = None;

    unsafe {
        if NAME.is_some() {
            NAME.unwrap()
        } else {
            let mut name: Vec<u8> = vec! [ 0; 32 ];
            let len: i32 = raw::runtime_name(&mut name[0], name.len());

            assert!(len >= 0);
            let len = len as usize;
            assert!(len <= name.len());

            name.set_len(len);
            let name = String::from_utf8(name).unwrap();

            let s: &'static str = ::std::mem::transmute::<&str, &'static str>(
                name.as_str()
            );
            ::std::mem::forget(name);

            NAME = Some(s);
            s
        }
    }
}
