(module
  (import "env" "print_string" (func $print_string (param i32 i32)))
  (import "env" "buffer" (memory 1))

  (data (i32.const 128) "0123456789ABCDEF")
  (data (i32.const 256) " 0")
  (data (i32.const 384) " 0x0")
  (data (i32.const 512) " 0000 0000 0000 0000 0000 0000 0000 0000")

  (global $dec_string_len i32 (i32.const 16))
  (global $hex_string_len i32 (i32.const 16))
  (global $bin_string_len i32 (i32.const 40))

  (func $set_dec_string (param $num i32) (param $string_len i32)
	(local $index i32)
	(local $digit_char i32)
	(local $digit_val i32)

	local.get $string_len
	local.set $index

	local.get $num
	i32.eqz
    if
	  local.get $index
	  i32.const 1
	  i32.sub
	  local.set $index
	  (i32.store8 offset=256 (local.get $index) (i32.const 48))
	end

	(loop $digit_loop (block $break
	  local.get $index
	  i32.eqz
	  br_if $break

	  local.get $num
	  i32.const 10
	  i32.rem_u

	  local.set $digit_val
	  local.get $num
	  i32.eqz
	  if
	    i32.const 32
		local.set $digit_char
	  else
		(i32.load8_u offset=128 (local.get $digit_val))
		local.set $digit_char
	  end
 
	  local.get $index
	  i32.const 1
	  i32.sub
	  local.set $index
	  (i32.store8 offset=256 (local.get $index) (local.get $digit_char)) 

	  local.get $num
	  i32.const 10
	  i32.div_u
	  local.set $num
	  br $digit_loop
	))
  )

  (func $set_hex_string (param $num i32) (param $string_len i32)
    (local $index i32)
	(local $digit_char i32)
	(local $digit_val i32)
	(local $x_pos i32)
	
    global.get $hex_string_len
	local.set $index
	
    (loop $digit_loop (block $break
	  local.get $index
	  i32.eqz 
	  br_if $break

	  local.get $num
	  i32.const 0xf
	  i32.and

	  local.set $digit_val
	  local.get $num
	  i32.eqz
	  if
	    local.get $x_pos
		i32.eqz
		if
		  local.get $index
		  local.set $x_pos
		else
		  i32.const 32
		  local.set $digit_char
		end
	  else
		(i32.load8_u offset=128 (local.get $digit_val))
		local.set $digit_char
	  end 
	  
      local.get $index
	  i32.const 1
	  i32.sub
	  local.tee $index
	  local.get $digit_char

	  i32.store8 offset=384
	  local.get $num
	  i32.const 4
	  i32.shr_u
	  local.set $num

	  br $digit_loop
	))
	
    local.get $x_pos
	i32.const 1
	i32.sub
	i32.const 120
	i32.store8 offset=384

	local.get $x_pos
	i32.const 2
	i32.sub

	i32.const 48
	i32.store8 offset=384
  )

  (func $set_bin_string (param $num i32) (param $string_len i32)
    (local $index i32)
	(local $loops_remaining i32)
	(local $nibble_bits i32)
	
    global.get $bin_string_len
	local.set $index

	i32.const 8
	local.set $loops_remaining
	
    (loop $bin_loop (block $outer_break
	  local.get $index 
	  i32.eqz
	  br_if $outer_break
	  i32.const 4
	  local.set $nibble_bits
	
      (loop $nibble_loop (block $nibble_break
	    local.get $index 
		i32.const 1
		i32.sub
		local.set $index

		local.get $num
		i32.const 1
		i32.and
        if
		  local.get $index
		  i32.const 49
		  i32.store8 offset=512
		else
		  local.get $index
		  i32.const 48
		  i32.store8 offset=512
		end
		
        local.get $num
		i32.const 1
		i32.shr_u
		local.set $num

		local.get $nibble_bits
		i32.const 1
		i32.sub
		local.tee $nibble_bits
		i32.eqz
		br_if $nibble_break

		br $nibble_loop
	  ))
		
      local.get $index 
	  i32.const 1
	  i32.sub
	  local.tee $index
	  i32.const 32
	  i32.store8 offset=512

	  br $bin_loop
	))
  )

  (func (export "to_string") (param $num i32)
    (call $set_dec_string (local.get $num) (global.get $dec_string_len))
	(call $print_string (i32.const 256) (global.get $dec_string_len))		
	
    (call $set_hex_string (local.get $num) (global.get $hex_string_len))
	(call $print_string (i32.const 384) (global.get $hex_string_len))	

	(call $set_bin_string (local.get $num) (global.get $bin_string_len))
	(call $print_string (i32.const 512) (global.get $bin_string_len))
  )
)