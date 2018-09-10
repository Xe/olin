extern crate log;

use self::log::*;

struct OlinLogger;

impl self::log::Log for OlinLogger {
    fn enabled(&self, metadata: &Metadata) -> bool {
        metadata.level() <= Level::Info
    }

    fn log(&self, record: &self::log::Record) {
        if self.enabled(record.metadata()) {
            match record.level() {
                Level::Info => ::log::info(&format!("{}", record.args())),
                Level::Warn => ::log::warning(&format!("{}", record.args())),
                Level::Error => ::log::error(&format!("{}", record.args())),
                _ => ::log::error(&format!(
                    "unknown level {:?} - {}",
                    record.level(),
                    record.args()
                )),
            }
            ::log::info(&format!("{} - {}", record.level(), record.args()));
        }
    }

    fn flush(&self) {}
}

static LOGGER: OlinLogger = OlinLogger;

pub fn init() -> Result<(), SetLoggerError> {
    log::set_logger(&LOGGER).map(|()| log::set_max_level(LevelFilter::Info))
}
