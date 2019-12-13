package orm

const (
	typeString  = "string"
	typeTime    = "Time"
	typeInt     = "int"
	typeUint    = "uint"
	typeUint8   = "uint8"
	typeUint16  = "uint16"
	typeUint32  = "uint32"
	typeUint64  = "uint64"
	typeInt8    = "int8"
	typeInt16   = "int16"
	typeInt32   = "int32"
	typeInt64   = "int64"
	typeFloat32 = "float32"
	typeFloat64 = "float64"
	typeBool    = "bool"
)

//TimeLayout ..
const timeLayout = "2006-01-02 15:04:05"

//LikeMatch ..
type LikeMatch uint8

const (
	// BothMatch like %xx%
	BothMatch LikeMatch = 0
	// LeftMatch like %xx
	LeftMatch LikeMatch = 1
	// RightMatch like xx%
	RightMatch LikeMatch = 2
)
