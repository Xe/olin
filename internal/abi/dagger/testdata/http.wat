(module
 ;; import functions from env
 (func $open  (import "dagger" "open")  (param i32 i32)     (result i32))
 (func $write (import "dagger" "write") (param i32 i32 i32) (result i32))
 (func $flush (import "dagger" "flush") (param i32)         (result i32))

 ;; memory
 (memory $mem 1)

 ;; constants
 (data (i32.const 200) "http://127.0.0.1:30405/")
 (data (i32.const 250) "GET / HTTP/1.1\nHost: 127.0.0.1:30405\nUser-Agent: Go-http-client/1.1\n\n")

 (func $main (result i32)
       (local $fd i32)
       (set_local $fd (call $open (i32.const 200) (i32.const 0)))

       (call $write (get_local $fd) (i32.const 250) (i32.const 69))
       (drop)

       (call $flush (get_local $fd)))
 (export "main" (func $main)))
