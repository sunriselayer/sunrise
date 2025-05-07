package consts

func CalculateTotalSize(rows uint32, cols uint32) uint64 {
	return uint64(rows) * uint64(cols) * ElementSize
}

func CalculateShardCountPerRow(cols uint32) uint32 {
	return cols / ElementsLenPerShard
}

func CalculateShardCount(rows uint32, cols uint32) uint32 {
	return rows * CalculateShardCountPerRow(cols)
}
