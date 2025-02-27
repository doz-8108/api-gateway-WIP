package utils

func (u *Utils) EncodeUserId(id uint64) string {
	seq := []uint64{u.SqIdInitNum, id}
	sid, _ := u.SqId.Encode(seq)
	return sid
}

func (u *Utils) DecodeUserId(id string) uint64 {
	seq := u.SqId.Decode(id)
	return seq[len(seq)-1]
}
