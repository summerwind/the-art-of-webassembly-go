(module
  (memory 1)
  (global $pointer i32 (i32.const 128))

  (func $init
    (i32.store
      (global.get $pointer)
      (i32.const 99)
    )
  )

  (func (export "get_ptr") (result i32)
	  (i32.load (global.get $pointer))
  )
  
  (start $init)
)