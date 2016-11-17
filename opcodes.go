// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

// Opcode is a wasm opcode.
type Opcode byte

// Language types opcodes as defined by:
// http://webassembly_org/docs/binary-encoding/#language-types
const (
	Op_i32     Opcode = 0x7f
	Op_i64            = 0x7e
	Op_f32            = 0x7d
	Op_f64            = 0x7c
	Op_anyfunc        = 0x70
	Op_func           = 0x60
	Op_empty          = 0x40
)

// Control flow operators
const (
	Op_unreachable Opcode = 0x00
	Op_nop                = 0x01
	Op_block              = 0x02
	Op_loop               = 0x03
	Op_if                 = 0x04
	Op_else               = 0x05
	Op_end                = 0x0b
	Op_br                 = 0x0c
	Op_br_if              = 0x0d
	Op_br_table           = 0x0e
	Op_return             = 0x0f
)

// Call operators
const (
	Op_call          Opcode = 0x10
	Op_call_indirect        = 0x11
)

// Parametric operators
const (
	Op_drop   Opcode = 0x1a
	Op_select        = 0x1b
)

// Variable access
const (
	Op_get_local  Opcode = 0x20
	Op_set_local         = 0x21
	Op_tee_local         = 0x22
	Op_get_global        = 0x23
	Op_set_global        = 0x24
)

// Memory-related operators
const (
	Op_i32_load       Opcode = 0x28
	Op_i64_load              = 0x29
	Op_f32_load              = 0x2a
	Op_f64_load              = 0x2b
	Op_i32_load8_s           = 0x2c
	Op_i32_load8_u           = 0x2d
	Op_i32_load16_s          = 0x2e
	Op_i32_load16_u          = 0x2f
	Op_i64_load8_s           = 0x30
	Op_i64_load8_u           = 0x31
	Op_i64_load16_s          = 0x32
	Op_i64_load16_u          = 0x33
	Op_i64_load32_s          = 0x34
	Op_i64_load32_u          = 0x35
	Op_i32_store             = 0x36
	Op_i64_store             = 0x37
	Op_f32_store             = 0x38
	Op_f64_store             = 0x39
	Op_i32_store8            = 0x3a
	Op_i32_store16           = 0x3b
	Op_i64_store8            = 0x3c
	Op_i64_store16           = 0x3d
	Op_i64_store32           = 0x3e
	Op_current_memory        = 0x3f
	Op_grow_memory           = 0x40
)

// Constants opcodes
const (
	Op_i32_const Opcode = 0x41
	Op_i64_const        = 0x42
	Op_f32_const        = 0x43
	Op_f64_const        = 0x44
)

// Comparison operators
const (
	Op_i32_eqz  Opcode = 0x45
	Op_i32_eq          = 0x46
	Op_i32_ne          = 0x47
	Op_i32_lt_s        = 0x48
	Op_i32_lt_u        = 0x49
	Op_i32_gt_s        = 0x4a
	Op_i32_gt_u        = 0x4b
	Op_i32_le_s        = 0x4c
	Op_i32_le_u        = 0x4d
	Op_i32_ge_s        = 0x4e
	Op_i32_ge_u        = 0x4f
	Op_i64_eqz         = 0x50
	Op_i64_eq          = 0x51
	Op_i64_ne          = 0x52
	Op_i64_lt_s        = 0x53
	Op_i64_lt_u        = 0x54
	Op_i64_gt_s        = 0x55
	Op_i64_gt_u        = 0x56
	Op_i64_le_s        = 0x57
	Op_i64_le_u        = 0x58
	Op_i64_ge_s        = 0x59
	Op_i64_ge_u        = 0x5a
	Op_f32_eq          = 0x5b
	Op_f32_ne          = 0x5c
	Op_f32_lt          = 0x5d
	Op_f32_gt          = 0x5e
	Op_f32_le          = 0x5f
	Op_f32_ge          = 0x60
	Op_f64_eq          = 0x61
	Op_f64_ne          = 0x62
	Op_f64_lt          = 0x63
	Op_f64_gt          = 0x64
	Op_f64_le          = 0x65
	Op_f64_ge          = 0x66
)

// Numeric operators
const (
	Op_i32_clz      Opcode = 0x67
	Op_i32_ctz             = 0x68
	Op_i32_popcnt          = 0x69
	Op_i32_add             = 0x6a
	Op_i32_sub             = 0x6b
	Op_i32_mul             = 0x6c
	Op_i32_div_s           = 0x6d
	Op_i32_div_u           = 0x6e
	Op_i32_rem_s           = 0x6f
	Op_i32_rem_u           = 0x70
	Op_i32_and             = 0x71
	Op_i32_or              = 0x72
	Op_i32_xor             = 0x73
	Op_i32_shl             = 0x74
	Op_i32_shr_s           = 0x75
	Op_i32_shr_u           = 0x76
	Op_i32_rotl            = 0x77
	Op_i32_rotr            = 0x78
	Op_i64_clz             = 0x79
	Op_i64_ctz             = 0x7a
	Op_i64_popcnt          = 0x7b
	Op_i64_add             = 0x7c
	Op_i64_sub             = 0x7d
	Op_i64_mul             = 0x7e
	Op_i64_div_s           = 0x7f
	Op_i64_div_u           = 0x80
	Op_i64_rem_s           = 0x81
	Op_i64_rem_u           = 0x82
	Op_i64_and             = 0x83
	Op_i64_or              = 0x84
	Op_i64_xor             = 0x85
	Op_i64_shl             = 0x86
	Op_i64_shr_s           = 0x87
	Op_i64_shr_u           = 0x88
	Op_i64_rotl            = 0x89
	Op_i64_rotr            = 0x8a
	Op_f32_abs             = 0x8b
	Op_f32_neg             = 0x8c
	Op_f32_ceil            = 0x8d
	Op_f32_floor           = 0x8e
	Op_f32_trunc           = 0x8f
	Op_f32_nearest         = 0x90
	Op_f32_sqrt            = 0x91
	Op_f32_add             = 0x92
	Op_f32_sub             = 0x93
	Op_f32_mul             = 0x94
	Op_f32_div             = 0x95
	Op_f32_min             = 0x96
	Op_f32_max             = 0x97
	Op_f32_copysign        = 0x98
	Op_f64_abs             = 0x99
	Op_f64_neg             = 0x9a
	Op_f64_ceil            = 0x9b
	Op_f64_floor           = 0x9c
	Op_f64_trunc           = 0x9d
	Op_f64_nearest         = 0x9e
	Op_f64_sqrt            = 0x9f
	Op_f64_add             = 0xa0
	Op_f64_sub             = 0xa1
	Op_f64_mul             = 0xa2
	Op_f64_div             = 0xa3
	Op_f64_min             = 0xa4
	Op_f64_max             = 0xa5
	Op_f64_copysign        = 0xa6
)

// Conversions
const (
	Op_i32_wrap_i64      Opcode = 0xa7
	Op_i32_trunc_s_f32          = 0xa8
	Op_i32_trunc_u_f32          = 0xa9
	Op_i32_trunc_s_f64          = 0xaa
	Op_i32_trunc_u_f64          = 0xab
	Op_i64_extend_s_i32         = 0xac
	Op_i64_extend_u_i32         = 0xad
	Op_i64_trunc_s_f32          = 0xae
	Op_i64_trunc_u_f32          = 0xaf
	Op_i64_trunc_s_f64          = 0xb0
	Op_i64_trunc_u_f64          = 0xb1
	Op_f32_convert_s_i32        = 0xb2
	Op_f32_convert_u_i32        = 0xb3
	Op_f32_convert_s_i64        = 0xb4
	Op_f32_convert_u_i64        = 0xb5
	Op_f32_demote_f64           = 0xb6
	Op_f64_convert_s_i32        = 0xb7
	Op_f64_convert_u_i32        = 0xb8
	Op_f64_convert_s_i64        = 0xb9
	Op_f64_convert_u_i64        = 0xba
	Op_f64_promote_f32          = 0xbb
)

// Reinterpretations
const (
	Op_i32_reinterpret_f32 Opcode = 0xbc
	Op_i64_reinterpret_f64        = 0xbd
	Op_f32_reinterpret_i32        = 0xbe
	Op_f64_reinterpret_i64        = 0xbf
)
