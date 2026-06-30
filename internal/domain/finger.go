package domain

// Finger identifies which finger is responsible for a key in touch typing.
type Finger int

const (
	LPinky Finger = iota
	LRing
	LMiddle
	LIndex
	RIndex
	RMiddle
	RRing
	RPinky
	Thumb
)

// FingerCount is the number of distinct fingers used for assignment and theming.
const FingerCount = 9

var fingerNames = [FingerCount]string{
	"left pinky", "left ring", "left middle", "left index",
	"right index", "right middle", "right ring", "right pinky", "thumb",
}

func (f Finger) String() string {
	if f < 0 || int(f) >= FingerCount {
		return "unknown"
	}
	return fingerNames[f]
}
