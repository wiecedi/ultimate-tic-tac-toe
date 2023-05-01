package main

import (
    "math/bits"
    "fmt"
)


type BitArray81 struct {
    v1 uint32
    v2 uint64
}

type BitArray9 struct {
    v1 uint16
}


func (ba *BitArray81) getNthPos(n int) int {
    if n >= ba.OnesCount() {
        panic("error in getNthPos")
    }

    if bits.OnesCount64(ba.v2) > n {
        v2 := ba.v2
        
        for i := 1; i <= n; i++ {
            z := bits.TrailingZeros64(v2)
            v2 ^= 1 << z
        }
        return bits.TrailingZeros64(v2)
    }
    n -= bits.OnesCount64(ba.v2)
    v1 := ba.v1
    
    for i := 1; i <= n; i++ {
        z := bits.TrailingZeros32(v1)
        v1 ^= 1 << z
    }
    return 64 + bits.TrailingZeros32(v1)
}

func (ba *BitArray81) OnesCount() int {
    return bits.OnesCount32(ba.v1 & 0x1FFFF) + bits.OnesCount64(ba.v2)
}

func (ba *BitArray81) And(other BitArray81) BitArray81 {
    return BitArray81 {
        v1: ba.v1 & other.v1 & 0x1FFFF,
        v2: ba.v2 & other.v2,
    }
}

func (ba *BitArray81) AndNot(other BitArray81) BitArray81 {
    return BitArray81 {
        v1: ba.v1 &^ other.v1 & 0x1FFFF,
        v2: ba.v2 &^ other.v2,
    }
}

func (ba *BitArray81) Or(other BitArray81) BitArray81 {
    return BitArray81 {
        v1: ba.v1 | other.v1 & 0x1FFFF,
        v2: ba.v2 | other.v2,
    }
}

func (ba *BitArray81) Not() BitArray81 {
    return BitArray81 {
        v1: ^ba.v1 & 0x1FFFF,
        v2: ^ba.v2,
    }
}

func (ba BitArray81) Equals(other BitArray81) bool {
    return ba.v1 & 0x1FFFF == other.v1 & 0x1FFFF && ba.v2 == other.v2
}

func (ba BitArray81) ToList() []Action {
    list := make([]Action, 0, 0)
    for i := 0; i < ba.OnesCount(); i++ {
        list = append(list, Action(ba.getNthPos(i)))
    }
    return list
}

func (ba *BitArray81) SetBit(k uint64) {
    if k < 64 {
        ba.v2 |= 1 << k
    } else {
        if k >= 81 {
            panic("value to high")
        }
        ba.v1 |= 1 << (k - 64)
    }
}

func (ba *BitArray81) ClearBit(k uint64) {
    if k < 64 {
        ba.v2 ^= 1 << k
    } else {
        if k >= 81 {
            panic("value to high")
        }
        ba.v1 ^= 1 << (k - 64)
    }
}

func (ba *BitArray81) GetBit(k uint64) bool {
    if k < 64 {
        return ba.v2 & (1 << k) != 0
    } else {
        if k >= 81 {
            panic("value to high")
        }
        return ba.v1 & (1 << (k - 64)) != 0
    }
}

func (ba *BitArray81) Print() {
    s := "[\n"
    for r := 0; r < 9; r++ {
        for c := 0; c < 9; c++ {
            if ba.GetBit(uint64(9*r+c)) {
                s += "X"
            } else {
                s += " "
            }
        }
        s += "\n"
    }
    fmt.Println(s + "]")
}

func (ba *BitArray9) OnesCount() int {
    return bits.OnesCount16(ba.v1 & 0x1FF)
}

func (ba *BitArray9) And(other BitArray9) BitArray9 {
    return BitArray9 {
        v1: ba.v1 & other.v1 & 0x1FF,
    }
}

func (ba *BitArray9) AndNot(other BitArray9) BitArray9 {
    return BitArray9 {
        v1: ba.v1 &^ other.v1 & 0x1FF,
    }
}

func (ba *BitArray9) Or(other BitArray9) BitArray9 {
    return BitArray9 {
        v1: ba.v1 | other.v1 & 0x1FF,
    }
}

func (ba *BitArray9) Not() BitArray9 {
    return BitArray9 {
        v1: ^ba.v1 & 0x1FF,
    }
}

func (ba BitArray9) Equals(other BitArray9) bool {
    return ba.v1 & 0x1FF == other.v1 & 0x1FF
}

func (ba *BitArray9) SetBit(k uint64) {
    if k >= 9 {
        panic("value to high")
    }
    ba.v1 |= 1 << k
}

func (ba *BitArray9) ClearBit(k uint64) {
    if k >= 9 {
        panic("value to high")
    }
    ba.v1 ^= 1 << k
}

func (ba *BitArray9) GetBit(k uint64) bool {
    if k >= 9 {
        panic("value to high")
    }
    return ba.v1 & (1 << k) != 0
}

func (ba *BitArray9) IsFull() bool {
    return ba.v1 >= 0x1FF
}

func (ba *BitArray9) Print() {
    s := "[\n"
    for r := 0; r < 3; r++ {
        for c := 0; c < 3; c++ {
            if ba.GetBit(uint64(3*r+c)) {
                s += "X"
            } else {
                s += " "
            }
        }
        s += "\n"
    }
    fmt.Println(s + "]")
}
