// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package disasm provides functions for disassembling WebAssembly bytecode.
package disasm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/leb128"
	ops "github.com/go-interpreter/wagon/wasm/operators"
)

// Instr describes an instruction, consisting of an operator, with its
// appropriate immediate value(s).
type Instr struct {
	Op ops.Op

	// Immediates are arguments to an operator in the bytecode stream itself.
	// Valid value types are:
	// - (u)(int/float)(32/64)
	// - wasm.BlockType
	Immediates  []interface{}
	NewStack    *StackInfo // non-nil if the instruction creates or unwinds a stack.
	Block       *BlockInfo // non-nil if the instruction starts or ends a new block.
	Unreachable bool       // whether the operator can be reached during execution
	// IsReturn is true if executing this instruction will result in the
	// function returning. This is true for branches (br, br_if) to
	// the depth <max_relative_depth> + 1, or the return operator itself.
	// If true, NewStack for this instruction is nil.
	IsReturn bool
	// If the operator is br_table (ops.BrTable), this is a list of StackInfo
	// fields for each of the blocks/branches referenced by the operator.
	Branches []StackInfo
}

// StackInfo stores details about a new stack created or unwinded by an instruction.
type StackInfo struct {
	StackTopDiff int64 // The difference between the stack depths at the end of the block
	PreserveTop  bool  // Whether the value on the top of the stack should be preserved while unwinding
	IsReturn     bool  // Whether the unwind is equivalent to a return
}

// BlockInfo stores details about a block created or ended by an instruction.
type BlockInfo struct {
	Start     bool           // If true, this instruction starts a block. Else this instruction ends it.
	Signature wasm.BlockType // The block signature

	// Indices to the accompanying control operator.
	// For 'if', this is the index to the 'else' operator.
	IfElseIndex int
	// For 'else', this is the index to the 'if' operator.
	ElseIfIndex int
	// The index to the `end' operator for if/else/loop/block.
	EndIndex int
	// For end, it is the index to the operator that starts the block.
	BlockStartIndex int
}

// Disassembly is the result of disassembling a WebAssembly function.
type Disassembly struct {
	Code     []Instr
	MaxDepth int // The maximum stack depth that can be reached while executing this function
}

func (d *Disassembly) checkMaxDepth(depth int) {
	if depth > d.MaxDepth {
		d.MaxDepth = depth
	}
}

func pushPolymorphicOp(indexStack [][]int, index int) {
	indexStack[len(indexStack)-1] = append(indexStack[len(indexStack)-1], index)
}

func isInstrReachable(indexStack [][]int) bool {
	return len(indexStack[len(indexStack)-1]) == 0
}

var ErrStackUnderflow = errors.New("disasm: stack underflow")

// Disassemble disassembles the given function. It also takes the function's
// parent module as an argument for locating any other functions referenced by
// fn.
func Disassemble(fn wasm.Function, module *wasm.Module) (*Disassembly, error) {
	code := fn.Body.Code
	reader := bytes.NewReader(code)
	disas := &Disassembly{}

	for {
		op, err := reader.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		opStr, err := ops.New(op)
		if err != nil {
			return nil, err
		}
		instr := Instr{
			Op: opStr,
		}

		switch op {
		case ops.Unreachable:
		case ops.Drop:
		case ops.Select:
		case ops.Return:
		case ops.End, ops.Else:

		case ops.Block, ops.Loop, ops.If:
			sig, err := leb128.ReadVarint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Block = &BlockInfo{
				Start:     true,
				Signature: wasm.BlockType(sig),
			}
			instr.Immediates = append(instr.Immediates, wasm.BlockType(sig))
		case ops.Br, ops.BrIf:
			depth, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, depth)

		case ops.BrTable:
			targetCount, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, targetCount)
			for i := uint32(0); i < targetCount; i++ {
				entry, err := leb128.ReadVarUint32(reader)
				if err != nil {
					return nil, err
				}
				instr.Immediates = append(instr.Immediates, entry)
			}

			defaultTarget, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, defaultTarget)
		case ops.Call, ops.CallIndirect:
			index, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, index)
			if op == ops.CallIndirect {
				reserved, err := leb128.ReadVarUint32(reader)
				if err != nil {
					return nil, err
				}
				instr.Immediates = append(instr.Immediates, reserved)
			}
		case ops.GetLocal, ops.SetLocal, ops.TeeLocal, ops.GetGlobal, ops.SetGlobal:
			index, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, index)
		case ops.I32Const:
			i, err := leb128.ReadVarint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, i)
		case ops.I64Const:
			i, err := leb128.ReadVarint64(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, i)
		case ops.F32Const:
			var b [4]byte
			if _, err := io.ReadFull(reader, b[:]); err != nil {
				return nil, err
			}
			i := binary.LittleEndian.Uint32(b[:])
			instr.Immediates = append(instr.Immediates, math.Float32frombits(i))
		case ops.F64Const:
			var b [8]byte
			if _, err := io.ReadFull(reader, b[:]); err != nil {
				return nil, err
			}
			i := binary.LittleEndian.Uint64(b[:])
			instr.Immediates = append(instr.Immediates, math.Float64frombits(i))
		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			// read memory_immediate
			flags, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, flags)

			offset, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, offset)
		case ops.CurrentMemory, ops.GrowMemory:
			res, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, uint8(res))
		}

		disas.Code = append(disas.Code, instr)
	}

	return disas, nil
}
