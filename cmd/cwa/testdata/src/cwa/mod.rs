//! Rust API library for the [CommonWA specification](https://github.com/CommonWA/cwa-spec).

#![feature(wasm_import_module)]

#[macro_use]
mod utils;

#[macro_use]
pub mod macros;

pub mod raw;
pub mod log;
pub mod env;
pub mod runtime;
pub mod startup;
pub mod resource;
pub mod io;
