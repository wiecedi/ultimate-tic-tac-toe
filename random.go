package main

import (
    "time"
)


type Random struct {
    x uint32
}


func GetSeed() uint32 {
    x := time.Now().UnixNano()
    return uint32((x >> 32) ^ x)
}

func (random *Random) Get(max int) int {
    if max == 0 {
        return 0
    }
    for random.x == 0 {
        random.x = GetSeed()
    }

    random.x ^= random.x << 13
    random.x ^= random.x >> 17
    random.x ^= random.x << 5
    
    return int(random.x) % max
}
