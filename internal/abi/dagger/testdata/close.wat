(module
 ;; import functions from env
 (func $open (import "dagger" "open") (param i32 i32) (result i32))
 (func $close (import "dagger" "close") (param i32) (result i32))

 ;; memory
 (memory $mem 1)

 ;; constants
 (data (i32.const 200) "fd://1")

 (func $main (result i32)
       (local $fd i32)
       (set_local $fd (call $open
             ;; pointer to the file name
             (i32.const 200)
             ;; flags, 0 because they don't matter here
             (i32.const 0)))
       (call $close (get_local $fd)))

 (export "main" (func $main)))
