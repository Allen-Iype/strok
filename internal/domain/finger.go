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

var fingerShort = [FingerCount]string{
	"L-pinky", "L-ring", "L-mid", "L-index",
	"R-index", "R-mid", "R-ring", "R-pinky", "thumb",
}

func (f Finger) String() string {
	if f < 0 || int(f) >= FingerCount {
		return "unknown"
	}
	return fingerNames[f]
}

// ShortName returns a compact finger label for legends, e.g. "L-pinky".
func (f Finger) ShortName() string {
	if f < 0 || int(f) >= FingerCount {
		return "?"
	}
	return fingerShort[f]
}
