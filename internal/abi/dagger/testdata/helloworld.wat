(module
 ;; import functions from env
 (func $open  (import "dagger" "open")  (param i32 i32)     (result i32))
 (func $write (import "dagger" "write") (param i32 i32 i32) (result i32))

 ;; memory
 (memory $mem 1)

 ;; constants
 (data (i32.const 200) "fd://1")
 (data (i32.const 230) "Hello, world!\n")

 (func $main (result i32)
       ;; $fd is the file descriptor of the file we're gonna open
       (local $fd i32)

       ;; $fd = $open("fd://1", 0);
       (set_local $fd
                  (call $open
                        ;; pointer to the file name
                        (i32.const 200)
                        ;; flags, 0 because they don't matter here
                        (i32.const 0)))

       ;; $write($fd, "Hello, World!\n", 14);
       (call $write
             (get_local $fd)
             (i32.const 230)
             (i32.const 14))
       (drop)

       (i32.const 0))
 (export "main" (func $main)))
