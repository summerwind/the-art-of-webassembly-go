(module
  (import "env" "log" (func $log (param i32 i32)))

  (func $loop_test (export "loop_test") (param $n i32) (result i32)
    (local $i i32)
    (local $factorial i32)

    (local.set $factorial (i32.const 1))

    (loop $continue (block $break
      (local.set $i
        (i32.add (local.get $i) (i32.const 1))
      )

      (local.set $factorial
        (i32.mul (local.get $i) (local.get $factorial))
      )

      (call $log (local.get $i) (local.get $factorial))

      (br_if $break
        (i32.eq (local.get $i) (local.get $n))
      )

      br $continue
    ))

    local.get $factorial
  )
)