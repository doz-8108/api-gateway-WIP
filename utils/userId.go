package utils

func (ctx *Utils) EncodeUserId(id uint64) string {
	seq := []uint64{ctx.SqIdInitNum, id}
	sid, _ := ctx.SqId.Encode(seq)
	return sid
}

func (ctx *Utils) DecodeUserId(id string) uint64 {
	seq := ctx.SqId.Decode(id)
	return seq[len(seq)-1]
}
