(module
  (import "env" "mem" (memory 1)) 

  (global $obj_base_addr (import "env" "obj_base_addr") i32)
  (global $obj_count (import "env" "obj_count") i32)
  (global $obj_stride (import "env" "obj_stride") i32)

  (global $x_offset (import "env" "x_offset") i32)
  (global $y_offset (import "env" "y_offset") i32)
  (global $radius_offset (import "env" "radius_offset") i32)
  (global $collision_offset (import "env" "collision_offset") i32)

  (func $collision_check
	(param $x1 i32) (param $y1 i32) (param $r1 i32)
	(param $x2 i32) (param $y2 i32) (param $r2 i32)
	(result i32)
	
    (local $x_diff_sq i32)
	(local $y_diff_sq i32)
	(local $r_sum_sq i32)
	
    local.get $x1
	local.get $x2
	i32.sub
	local.tee $x_diff_sq
	local.get $x_diff_sq
	i32.mul
	local.set $x_diff_sq
	
	local.get $y1
	local.get $y2
	i32.sub
	local.tee $y_diff_sq
	local.get $y_diff_sq
	i32.mul
	local.set $y_diff_sq

	local.get $r1
	local.get $r2
	i32.add
	local.tee $r_sum_sq
	local.get $r_sum_sq
	i32.mul
	local.tee $r_sum_sq

	local.get $x_diff_sq
	local.get $y_diff_sq
	i32.add

	i32.gt_u
  )

  (func $get_attr
    (param $obj_base i32) (param $attr_offset i32)
	(result i32)

	local.get $obj_base
	local.get $attr_offset
	i32.add
	i32.load
  )

  (func $set_collision
	(param $obj_base_1 i32) (param $obj_base_2 i32)

	local.get $obj_base_1
	global.get $collision_offset
	i32.add
	i32.const 1
	i32.store

	local.get $obj_base_2
	global.get $collision_offset
	i32.add
	i32.const 1
	i32.store
  )

  (func $init
	(local $i i32)
	(local $i_obj i32)
	(local $xi i32)(local $yi i32)(local $ri i32)

	(local $j i32)
	(local $j_obj i32)
	(local $xj i32)(local $yj i32)(local $rj i32)
	
    (loop $outer_loop
	  (local.set $j (i32.const 0))
	
	  (loop $inner_loop
        (block $inner_continue
		  (br_if $inner_continue (i32.eq (local.get $i) (local.get $j) ) )

		  (i32.add (global.get $obj_base_addr) (i32.mul (local.get $i) (global.get $obj_stride)))

		  (call $get_attr (local.tee $i_obj) (global.get $x_offset) )
		  local.set $xi 

		  (call $get_attr (local.get $i_obj) (global.get $y_offset) )
	      local.set $yi 

		  (call $get_attr (local.get $i_obj) (global.get $radius_offset) )
		  local.set $ri 

		  (i32.add (global.get $obj_base_addr) (i32.mul (local.get $j)(global.get $obj_stride)))

		  (call $get_attr (local.tee $j_obj) (global.get $x_offset) )
		  local.set $xj 

		  (call $get_attr (local.get $j_obj) (global.get $y_offset) )
		  local.set $yj 

		  (call $get_attr (local.get $j_obj) (global.get $radius_offset) )
		  local.set $rj 

		  (call $collision_check
		    (local.get $xi)(local.get $yi)(local.get $ri)
		    (local.get $xj)(local.get $yj)(local.get $rj)
		  )
		  if
		    (call $set_collision (local.get $i_obj) (local.get $j_obj))
		  end
	    )
	  
        (i32.add (local.get $j) (i32.const 1))

		(br_if $inner_loop (i32.lt_u (local.tee $j) (global.get $obj_count)))
	  )

	  (i32.add (local.get $i) (i32.const 1))

	  (br_if $outer_loop (i32.lt_u (local.tee $i) (global.get $obj_count)))
	)
  )

  (start $init)
)