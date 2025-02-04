package text

import (
	"fmt"

	"github.com/tdewolff/canvas/font"
)

type ScriptItem struct {
	Script
	Text string
}

// ScriptItemizer divides the string in parts for each different script.
func ScriptItemizer(runes []rune, embeddingLevels []int) []ScriptItem {
	i := 0
	var curLevel int
	var curScript Script
	items := []ScriptItem{}
	for j, r := range runes {
		script, level := LookupScript(r), embeddingLevels[j]
		if j == 0 {
			curLevel = level
			curScript = script
		} else if level != curLevel || script != curScript && script != ScriptInherited && script != ScriptCommon && curScript != ScriptInherited && curScript != ScriptCommon {
			items = append(items, ScriptItem{
				Script: curScript,
				Text:   string(runes[i:j]),
			})
			curLevel = level
			curScript = script
			i = j
		} else if curScript == ScriptInherited || curScript == ScriptCommon {
			curScript = script
		}
	}
	items = append(items, ScriptItem{
		Script: curScript,
		Text:   string(runes[i:]),
	})
	return items
}

// Glyph is a shaped glyph for the given font and font size. It specified the glyph ID, the cluster ID, its X and Y advance and offset in font units, and its representation as text.
type Glyph struct {
	SFNT *font.SFNT
	Size float64
	Script
	Vertical bool // is false for Latin/Mongolian/etc in a vertical layout

	ID       uint16
	Cluster  uint32
	XAdvance int32
	YAdvance int32
	XOffset  int32
	YOffset  int32
	Text     rune
}

func (g Glyph) String() string {
	return fmt.Sprintf("%s GID=%v Cluster=%v Adv=(%v,%v) Off=(%v,%v)", string(g.Text), g.ID, g.Cluster, g.XAdvance, g.YAdvance, g.XOffset, g.YOffset)
}

func (g Glyph) Rotation() Rotation {
	rot := NoRotation
	if !g.Vertical {
		rot = ScriptRotation(g.Script)
		if rot == NoRotation {
			rot = CW
		}
	}
	return rot
}

// TODO: implement Liang's (soft) hyphenation algorithm? Add \u00AD at opportunities, unless \u2060 or \uFEFF is present

// IsParagraphSeparator returns true for paragraph separator runes.
func IsParagraphSeparator(r rune) bool {
	// line feed, vertical tab, form feed, carriage return, next line, line separator, paragraph separator
	return 0x0A <= r && r <= 0x0D || r == 0x85 || r == '\u2028' || r == '\u2029'
}

func IsSpacelessScript(script Script) bool {
	// missing: S'gaw Karen
	return script == Han || script == Hangul || script == Katakana || script == Khmer || script == Lao || script == PhagsPa || script == Brahmi || script == TaiTham || script == NewTaiLue || script == TaiLe || script == TaiViet || script == Thai || script == Tibetan || script == Myanmar
}

func IsVerticalScript(script Script) bool {
	return script == Bopomofo || script == EgyptianHieroglyphs || script == Hiragana || script == Katakana || script == Han || script == Hangul || script == MeroiticCursive || script == MeroiticHieroglyphs || script == Mongolian || script == Ogham || script == OldTurkic || script == PhagsPa || script == Yi
}

type Rotation float64

const (
	NoRotation Rotation = 0.0
	CW         Rotation = -90.0
	CCW        Rotation = 90.0
)

func ScriptRotation(script Script) Rotation {
	if script == Mongolian || script == PhagsPa {
		return CW
	} else if script == Ogham || script == OldTurkic {
		return CCW
	}
	return NoRotation
}
