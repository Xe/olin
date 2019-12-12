# Run Twiggy

To build the library and run Twiggy, in this directory:

```
rustup install nightly
rustup update nightly
rustup target add wasm32-unknown-unknown --toolchain nightly

chmod +x build.sh
./build.sh

twiggy <subcommand> cwagi.wasm
```