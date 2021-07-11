package promptx

////////////////////////////////////////////////////////////////////////////////
// Mode Interface

// ConsoleBlocks Mode Interface
type ConsoleBlocks interface {
	InitBlocks()
	Active() bool
	SetActive(active bool)
	IsDraw(status int) bool
	Render(ctx PrintContext, preCursor int) (cursor int)
	OnEvent(ctx PressContext, key Key, in []byte) (exit bool)
	ResetBuffer()
	GetBuffer() *Buffer
}

// EmptyBlocks empty for basic operation
type EmptyBlocks struct {
	active  bool
	keyBind map[Key]KeyBindFunc
	// is draw
	isDraw func(status int) bool
}

// InitBlocks init blocks
func (m *EmptyBlocks) InitBlocks() {
	// default is active status
	m.active = true
	return
}

// Active is blocks active
func (m *EmptyBlocks) Active() bool {
	return m.active
}

// SetActive set active status
func (m *EmptyBlocks) SetActive(active bool) {
	m.active = active
	return
}

func (m *EmptyBlocks) IsDraw(status int) bool {
	if m.isDraw == nil {
		return true
	}
	return m.isDraw(status)
}

func (m *EmptyBlocks) SetIsDraw(f func(status int) (draw bool)) {
	m.isDraw = f
}

// Render rendering blocks.
func (m *EmptyBlocks) Render(ctx PrintContext, preCursor int) (cursor int) {
	return
}

// OnEvent deal console key press
func (m *EmptyBlocks) OnEvent(ctx PressContext, key Key, in []byte) (exit bool) {
	if m.keyBind == nil {
		return
	}
	if !m.Active() {
		return
	}
	if bind, ok := m.keyBind[key]; ok {
		if bind(ctx) {
			exit = true
		}
	}
	if key == NotDefined && len(in) == 1 {
		// ascii code event
		key = Key(in[0]) + NotDefined + 1
		if bind, ok := m.keyBind[key]; ok {
			if bind(ctx) {
				exit = true
			}
		}
	}

	return
}

// ResetBuffer reset buffer
func (m EmptyBlocks) ResetBuffer() {}

// GetBuffer get input buffer
func (m *EmptyBlocks) GetBuffer() *Buffer {
	return nil
}

// BindKey bind key funcs
func (m *EmptyBlocks) BindKey(bind KeyBindFunc, keys ...Key) {
	if m.keyBind == nil {
		m.keyBind = map[Key]KeyBindFunc{}
	}
	for _, key := range keys {
		m.keyBind[key] = bind
	}
}

// BindASCII bind ascii code func
func (m *EmptyBlocks) BindASCII(bind KeyBindFunc, ins ...byte) {
	if m.keyBind == nil {
		m.keyBind = map[Key]KeyBindFunc{}
	}
	for _, in := range ins {
		m.keyBind[Key(in)+NotDefined+1] = bind
	}
}
