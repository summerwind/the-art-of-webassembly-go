(module
  (import "go" "increment" (func $go_increment (result i32)))
  (import "go" "decrement" (func $go_decrement (result i32)))
  (table $tbl (export "tbl") 4 funcref)

  (global $i (mut i32) (i32.const 0))

  (func $increment (export "increment") (result i32)
    (global.set $i (i32.add (global.get $i) (i32.const 1)))
    global.get $i
  )

  (func $decrement (export "decrement") (result i32)
	  (global.set $i (i32.sub (global.get $i) (i32.const 1)))
	  global.get $i
  )

  (elem (i32.const 0) $go_increment $go_decrement $increment $decrement)
)