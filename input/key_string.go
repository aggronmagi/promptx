// Code generated by "stringer -type=Key"; DO NOT EDIT.

package input

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Escape-0]
	_ = x[ControlA-1]
	_ = x[ControlB-2]
	_ = x[ControlC-3]
	_ = x[ControlD-4]
	_ = x[ControlE-5]
	_ = x[ControlF-6]
	_ = x[ControlG-7]
	_ = x[ControlH-8]
	_ = x[controlI-9]
	_ = x[controlJ-10]
	_ = x[ControlK-11]
	_ = x[ControlL-12]
	_ = x[controlM-13]
	_ = x[ControlN-14]
	_ = x[ControlO-15]
	_ = x[ControlP-16]
	_ = x[ControlQ-17]
	_ = x[ControlR-18]
	_ = x[ControlS-19]
	_ = x[ControlT-20]
	_ = x[ControlU-21]
	_ = x[ControlV-22]
	_ = x[ControlW-23]
	_ = x[ControlX-24]
	_ = x[ControlY-25]
	_ = x[ControlZ-26]
	_ = x[MetaA-27]
	_ = x[MetaB-28]
	_ = x[MetaC-29]
	_ = x[MetaD-30]
	_ = x[MetaE-31]
	_ = x[MetaF-32]
	_ = x[MetaG-33]
	_ = x[MetaH-34]
	_ = x[MetaI-35]
	_ = x[MetaJ-36]
	_ = x[MetaK-37]
	_ = x[MetaL-38]
	_ = x[MetaM-39]
	_ = x[MetaN-40]
	_ = x[MetaO-41]
	_ = x[MetaP-42]
	_ = x[MetaQ-43]
	_ = x[MetaR-44]
	_ = x[MetaS-45]
	_ = x[MetaT-46]
	_ = x[MetaU-47]
	_ = x[MetaV-48]
	_ = x[MetaW-49]
	_ = x[MetaX-50]
	_ = x[MetaY-51]
	_ = x[MetaZ-52]
	_ = x[MetaShiftA-53]
	_ = x[MetaShiftB-54]
	_ = x[MetaShiftC-55]
	_ = x[MetaShiftD-56]
	_ = x[MetaShiftE-57]
	_ = x[MetaShiftF-58]
	_ = x[MetaShiftG-59]
	_ = x[MetaShiftH-60]
	_ = x[MetaShiftI-61]
	_ = x[MetaShiftJ-62]
	_ = x[MetaShiftK-63]
	_ = x[MetaShiftL-64]
	_ = x[MetaShiftM-65]
	_ = x[MetaShiftN-66]
	_ = x[MetaShiftO-67]
	_ = x[MetaShiftP-68]
	_ = x[MetaShiftQ-69]
	_ = x[MetaShiftR-70]
	_ = x[MetaShiftS-71]
	_ = x[MetaShiftT-72]
	_ = x[MetaShiftU-73]
	_ = x[MetaShiftV-74]
	_ = x[MetaShiftW-75]
	_ = x[MetaShiftX-76]
	_ = x[MetaShiftY-77]
	_ = x[MetaShiftZ-78]
	_ = x[ControlSpace-79]
	_ = x[ControlBackslash-80]
	_ = x[ControlSquareClose-81]
	_ = x[ControlCircumflex-82]
	_ = x[ControlUnderscore-83]
	_ = x[ControlLeft-84]
	_ = x[ControlRight-85]
	_ = x[ControlUp-86]
	_ = x[ControlDown-87]
	_ = x[Up-88]
	_ = x[Down-89]
	_ = x[Right-90]
	_ = x[Left-91]
	_ = x[ShiftLeft-92]
	_ = x[ShiftUp-93]
	_ = x[ShiftDown-94]
	_ = x[ShiftRight-95]
	_ = x[Home-96]
	_ = x[End-97]
	_ = x[Delete-98]
	_ = x[ShiftDelete-99]
	_ = x[ControlDelete-100]
	_ = x[PageUp-101]
	_ = x[PageDown-102]
	_ = x[BackTab-103]
	_ = x[Insert-104]
	_ = x[Backspace-105]
	_ = x[Tab-106]
	_ = x[Enter-107]
	_ = x[F1-108]
	_ = x[F2-109]
	_ = x[F3-110]
	_ = x[F4-111]
	_ = x[F5-112]
	_ = x[F6-113]
	_ = x[F7-114]
	_ = x[F8-115]
	_ = x[F9-116]
	_ = x[F10-117]
	_ = x[F11-118]
	_ = x[F12-119]
	_ = x[F13-120]
	_ = x[F14-121]
	_ = x[F15-122]
	_ = x[F16-123]
	_ = x[F17-124]
	_ = x[F18-125]
	_ = x[F19-126]
	_ = x[F20-127]
	_ = x[F21-128]
	_ = x[F22-129]
	_ = x[F23-130]
	_ = x[F24-131]
	_ = x[Any-132]
	_ = x[CPRResponse-133]
	_ = x[Vt100MouseEvent-134]
	_ = x[WindowsMouseEvent-135]
	_ = x[BracketedPaste-136]
	_ = x[Ignore-137]
	_ = x[NotDefined-138]
}

const _Key_name = "EscapeControlAControlBControlCControlDControlEControlFControlGControlHcontrolIcontrolJControlKControlLcontrolMControlNControlOControlPControlQControlRControlSControlTControlUControlVControlWControlXControlYControlZMetaAMetaBMetaCMetaDMetaEMetaFMetaGMetaHMetaIMetaJMetaKMetaLMetaMMetaNMetaOMetaPMetaQMetaRMetaSMetaTMetaUMetaVMetaWMetaXMetaYMetaZMetaShiftAMetaShiftBMetaShiftCMetaShiftDMetaShiftEMetaShiftFMetaShiftGMetaShiftHMetaShiftIMetaShiftJMetaShiftKMetaShiftLMetaShiftMMetaShiftNMetaShiftOMetaShiftPMetaShiftQMetaShiftRMetaShiftSMetaShiftTMetaShiftUMetaShiftVMetaShiftWMetaShiftXMetaShiftYMetaShiftZControlSpaceControlBackslashControlSquareCloseControlCircumflexControlUnderscoreControlLeftControlRightControlUpControlDownUpDownRightLeftShiftLeftShiftUpShiftDownShiftRightHomeEndDeleteShiftDeleteControlDeletePageUpPageDownBackTabInsertBackspaceTabEnterF1F2F3F4F5F6F7F8F9F10F11F12F13F14F15F16F17F18F19F20F21F22F23F24AnyCPRResponseVt100MouseEventWindowsMouseEventBracketedPasteIgnoreNotDefined"

var _Key_index = [...]uint16{0, 6, 14, 22, 30, 38, 46, 54, 62, 70, 78, 86, 94, 102, 110, 118, 126, 134, 142, 150, 158, 166, 174, 182, 190, 198, 206, 214, 219, 224, 229, 234, 239, 244, 249, 254, 259, 264, 269, 274, 279, 284, 289, 294, 299, 304, 309, 314, 319, 324, 329, 334, 339, 344, 354, 364, 374, 384, 394, 404, 414, 424, 434, 444, 454, 464, 474, 484, 494, 504, 514, 524, 534, 544, 554, 564, 574, 584, 594, 604, 616, 632, 650, 667, 684, 695, 707, 716, 727, 729, 733, 738, 742, 751, 758, 767, 777, 781, 784, 790, 801, 814, 820, 828, 835, 841, 850, 853, 858, 860, 862, 864, 866, 868, 870, 872, 874, 876, 879, 882, 885, 888, 891, 894, 897, 900, 903, 906, 909, 912, 915, 918, 921, 924, 935, 950, 967, 981, 987, 997}

func (i Key) String() string {
	if i >= Key(len(_Key_index)-1) {
		return "Key(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Key_name[_Key_index[i]:_Key_index[i+1]]
}