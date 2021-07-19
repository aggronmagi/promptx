package input

// WinSize represents the width and height of terminal.
type WinSize struct {
	Row int
	Col int
}

// key 映射
var keyMap map[uint64]Key

func init() {
	if keyMap == nil {
		keyMap = make(map[uint64]Key)
	}
	for _, k := range asciiSequences {
		if v, ok := convertUint64(k.ASCIICode); ok {
			keyMap[v] = k.Key
		}
	}
}

func convertUint64(b []byte) (flag uint64, ok bool) {
	if len(b) > 8 {
		return
	}
	for _, v := range b {
		flag = flag << 8
		flag |= uint64(v)
	}
	ok = true
	return
}

// InstallASCIICodeParse install ascii code parse.
// NOTE: it must call in init func
func InstallASCIICodeParse(codes ...*ASCIICode) {
	if keyMap == nil {
		keyMap = make(map[uint64]Key)
	}
	for _, k := range codes {
		if v, ok := convertUint64(k.ASCIICode); ok {
			keyMap[v] = k.Key
		}
	}
}

// GetKey returns Key correspond to input byte codes.
func GetKey(b []byte) Key {
	flag, ok := convertUint64(b)
	if !ok {
		return NotDefined
	}
	if k, ok := keyMap[flag]; ok {
		return k
	}
	// ignore other input begin by '0x1b'. it is control ascii code
	if b[0] == 0x1b {
		return Ignore
	}
	return NotDefined
}

// ConsoleParser is an interface to abstract input layer.
type ConsoleParser interface {
	// Setup should be called before starting input
	Setup() error
	// TearDown should be called after stopping input
	TearDown() error
	// GetWinSize returns WinSize object to represent width and height of terminal.
	GetWinSize() *WinSize
	// Read returns byte array.
	Read() ([]byte, error)
}

// asciiSequences holds mappings of the key and byte array.
var asciiSequences = []*ASCIICode{
	{Key: Escape, ASCIICode: []byte{0x1b}},

	{Key: ControlSpace, ASCIICode: []byte{0x00}},
	{Key: ControlA, ASCIICode: []byte{0x1}},
	{Key: ControlB, ASCIICode: []byte{0x2}},
	{Key: ControlC, ASCIICode: []byte{0x3}},
	{Key: ControlD, ASCIICode: []byte{0x4}},
	{Key: ControlE, ASCIICode: []byte{0x5}},
	{Key: ControlF, ASCIICode: []byte{0x6}},
	{Key: ControlG, ASCIICode: []byte{0x7}},
	{Key: ControlH, ASCIICode: []byte{0x8}},
	//{Key: controlI, ASCIICode: []byte{0x9}},
	//{Key: controlJ, ASCIICode: []byte{0xa}},
	{Key: ControlK, ASCIICode: []byte{0xb}},
	{Key: ControlL, ASCIICode: []byte{0xc}},
	//{Key: controlM, ASCIICode: []byte{0xd}},
	{Key: ControlN, ASCIICode: []byte{0xe}},
	{Key: ControlO, ASCIICode: []byte{0xf}},
	{Key: ControlP, ASCIICode: []byte{0x10}},
	{Key: ControlQ, ASCIICode: []byte{0x11}},
	{Key: ControlR, ASCIICode: []byte{0x12}},
	{Key: ControlS, ASCIICode: []byte{0x13}},
	{Key: ControlT, ASCIICode: []byte{0x14}},
	{Key: ControlU, ASCIICode: []byte{0x15}},
	{Key: ControlV, ASCIICode: []byte{0x16}},
	{Key: ControlW, ASCIICode: []byte{0x17}},
	{Key: ControlX, ASCIICode: []byte{0x18}},
	{Key: ControlY, ASCIICode: []byte{0x19}},
	{Key: ControlZ, ASCIICode: []byte{0x1a}},

	{Key: MetaA, ASCIICode: []byte{0x1b, 0x60 + 0x1}},
	{Key: MetaB, ASCIICode: []byte{0x1b, 0x60 + 0x2}},
	{Key: MetaC, ASCIICode: []byte{0x1b, 0x60 + 0x3}},
	{Key: MetaD, ASCIICode: []byte{0x1b, 0x60 + 0x4}},
	{Key: MetaE, ASCIICode: []byte{0x1b, 0x60 + 0x5}},
	{Key: MetaF, ASCIICode: []byte{0x1b, 0x60 + 0x6}},
	{Key: MetaG, ASCIICode: []byte{0x1b, 0x60 + 0x7}},
	{Key: MetaH, ASCIICode: []byte{0x1b, 0x60 + 0x8}},
	{Key: MetaI, ASCIICode: []byte{0x1b, 0x60 + 0x9}},
	{Key: MetaJ, ASCIICode: []byte{0x1b, 0x60 + 0xa}},
	{Key: MetaK, ASCIICode: []byte{0x1b, 0x60 + 0xb}},
	{Key: MetaL, ASCIICode: []byte{0x1b, 0x60 + 0xc}},
	{Key: MetaM, ASCIICode: []byte{0x1b, 0x60 + 0xd}},
	{Key: MetaN, ASCIICode: []byte{0x1b, 0x60 + 0xe}},
	{Key: MetaO, ASCIICode: []byte{0x1b, 0x60 + 0xf}},
	{Key: MetaP, ASCIICode: []byte{0x1b, 0x60 + 0x10}},
	{Key: MetaQ, ASCIICode: []byte{0x1b, 0x60 + 0x11}},
	{Key: MetaR, ASCIICode: []byte{0x1b, 0x60 + 0x12}},
	{Key: MetaS, ASCIICode: []byte{0x1b, 0x60 + 0x13}},
	{Key: MetaT, ASCIICode: []byte{0x1b, 0x60 + 0x14}},
	{Key: MetaU, ASCIICode: []byte{0x1b, 0x60 + 0x15}},
	{Key: MetaV, ASCIICode: []byte{0x1b, 0x60 + 0x16}},
	{Key: MetaW, ASCIICode: []byte{0x1b, 0x60 + 0x17}},
	{Key: MetaX, ASCIICode: []byte{0x1b, 0x60 + 0x18}},
	{Key: MetaY, ASCIICode: []byte{0x1b, 0x60 + 0x19}},
	{Key: MetaZ, ASCIICode: []byte{0x1b, 0x60 + 0x1a}},

	{Key: MetaShiftA, ASCIICode: []byte{0x1b, 0x40 + 0x1}},
	{Key: MetaShiftB, ASCIICode: []byte{0x1b, 0x40 + 0x2}},
	{Key: MetaShiftC, ASCIICode: []byte{0x1b, 0x40 + 0x3}},
	{Key: MetaShiftD, ASCIICode: []byte{0x1b, 0x40 + 0x4}},
	{Key: MetaShiftE, ASCIICode: []byte{0x1b, 0x40 + 0x5}},
	{Key: MetaShiftF, ASCIICode: []byte{0x1b, 0x40 + 0x6}},
	{Key: MetaShiftG, ASCIICode: []byte{0x1b, 0x40 + 0x7}},
	{Key: MetaShiftH, ASCIICode: []byte{0x1b, 0x40 + 0x8}},
	{Key: MetaShiftI, ASCIICode: []byte{0x1b, 0x40 + 0x9}},
	{Key: MetaShiftJ, ASCIICode: []byte{0x1b, 0x40 + 0xa}},
	{Key: MetaShiftK, ASCIICode: []byte{0x1b, 0x40 + 0xb}},
	{Key: MetaShiftL, ASCIICode: []byte{0x1b, 0x40 + 0xc}},
	{Key: MetaShiftM, ASCIICode: []byte{0x1b, 0x40 + 0xd}},
	{Key: MetaShiftN, ASCIICode: []byte{0x1b, 0x40 + 0xe}},
	{Key: MetaShiftO, ASCIICode: []byte{0x1b, 0x40 + 0xf}},
	{Key: MetaShiftP, ASCIICode: []byte{0x1b, 0x40 + 0x10}},
	{Key: MetaShiftQ, ASCIICode: []byte{0x1b, 0x40 + 0x11}},
	{Key: MetaShiftR, ASCIICode: []byte{0x1b, 0x40 + 0x12}},
	{Key: MetaShiftS, ASCIICode: []byte{0x1b, 0x40 + 0x13}},
	{Key: MetaShiftT, ASCIICode: []byte{0x1b, 0x40 + 0x14}},
	{Key: MetaShiftU, ASCIICode: []byte{0x1b, 0x40 + 0x15}},
	{Key: MetaShiftV, ASCIICode: []byte{0x1b, 0x40 + 0x16}},
	{Key: MetaShiftW, ASCIICode: []byte{0x1b, 0x40 + 0x17}},
	{Key: MetaShiftX, ASCIICode: []byte{0x1b, 0x40 + 0x18}},
	{Key: MetaShiftY, ASCIICode: []byte{0x1b, 0x40 + 0x19}},
	{Key: MetaShiftZ, ASCIICode: []byte{0x1b, 0x40 + 0x1a}},

	{Key: ControlBackslash, ASCIICode: []byte{0x1c}},
	{Key: ControlSquareClose, ASCIICode: []byte{0x1d}},
	{Key: ControlCircumflex, ASCIICode: []byte{0x1e}},
	{Key: ControlUnderscore, ASCIICode: []byte{0x1f}},
	{Key: Backspace, ASCIICode: []byte{0x7f}},

	{Key: Up, ASCIICode: []byte{0x1b, 0x5b, 0x41}},
	{Key: Down, ASCIICode: []byte{0x1b, 0x5b, 0x42}},
	{Key: Right, ASCIICode: []byte{0x1b, 0x5b, 0x43}},
	{Key: Left, ASCIICode: []byte{0x1b, 0x5b, 0x44}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x48}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x30, 0x48}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x46}},
	{Key: End, ASCIICode: []byte{0x1b, 0x30, 0x46}},

	{Key: Enter, ASCIICode: []byte{0xd}}, // Alias controlM
	{Key: Enter, ASCIICode: []byte{0xa}}, // Alias controlJ
	{Key: Delete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: ShiftDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x32, 0x7e}},
	{Key: ControlDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x35, 0x7e}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x34, 0x7e}},
	{Key: PageUp, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x7e}},
	{Key: PageDown, ASCIICode: []byte{0x1b, 0x5b, 0x36, 0x7e}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x37, 0x7e}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x38, 0x7e}},
	{Key: Tab, ASCIICode: []byte{0x9}}, // Alias controlI
	{Key: BackTab, ASCIICode: []byte{0x1b, 0x5b, 0x5a}},
	{Key: Insert, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x7e}},

	{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50}},
	{Key: F2, ASCIICode: []byte{0x1b, 0x4f, 0x51}},
	{Key: F3, ASCIICode: []byte{0x1b, 0x4f, 0x52}},
	{Key: F4, ASCIICode: []byte{0x1b, 0x4f, 0x53}},

	{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50, 0x41}}, // Linux console
	{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x42}}, // Linux console
	{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x43}}, // Linux console
	{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x44}}, // Linux console
	{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x45}}, // Linux console

	{Key: F1, ASCIICode: []byte{0x1b, 0x5b, 0x11, 0x7e}}, // rxvt-unicode
	{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x12, 0x7e}}, // rxvt-unicode
	{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x13, 0x7e}}, // rxvt-unicode
	{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x14, 0x7e}}, // rxvt-unicode

	{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}},
	{Key: F6, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}},
	{Key: F7, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}},
	{Key: F8, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}},
	{Key: F9, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}},
	{Key: F10, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}},
	{Key: F11, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x32, 0x7e}},
	{Key: F12, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x34, 0x7e, 0x8}},
	{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x25, 0x7e}},
	{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x26, 0x7e}},
	{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x28, 0x7e}},
	{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x29, 0x7e}},
	{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x7e}},
	{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x34, 0x7e}},

	// Xterm
	{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x50}},
	{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x51}},
	// &ASCIICode{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},  // Conflicts with CPR response
	{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},
	{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x15, 0x3b, 0x32, 0x7e}},
	{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x17, 0x3b, 0x32, 0x7e}},
	{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x18, 0x3b, 0x32, 0x7e}},
	{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x19, 0x3b, 0x32, 0x7e}},
	{Key: F21, ASCIICode: []byte{0x1b, 0x5b, 0x20, 0x3b, 0x32, 0x7e}},
	{Key: F22, ASCIICode: []byte{0x1b, 0x5b, 0x21, 0x3b, 0x32, 0x7e}},
	{Key: F23, ASCIICode: []byte{0x1b, 0x5b, 0x23, 0x3b, 0x32, 0x7e}},
	{Key: F24, ASCIICode: []byte{0x1b, 0x5b, 0x24, 0x3b, 0x32, 0x7e}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}},

	{Key: ShiftUp, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x41}},
	{Key: ShiftDown, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x42}},
	{Key: ShiftRight, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x43}},
	{Key: ShiftLeft, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x44}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: Up, ASCIICode: []byte{0x1b, 0x4f, 0x41}},
	{Key: Down, ASCIICode: []byte{0x1b, 0x4f, 0x42}},
	{Key: Right, ASCIICode: []byte{0x1b, 0x4f, 0x43}},
	{Key: Left, ASCIICode: []byte{0x1b, 0x4f, 0x44}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x41}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x42}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x43}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x44}},

	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x63}}, // rxvt
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x64}},  // rxvt

	{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x45}}, // Xterm
	{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x46}}, // Linux console
}
