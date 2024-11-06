// Package utils provides utilities for interacting with the terminal and formatting output.
package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func Table(style string, title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)

	t.SetStyle(getStyle(style))

	// TO DO: Custom
	//t.Style().Title.Colors = text.Colors{text.FgHiBlack, text.BgHiGreen}
	//t.Style().Color.Header = text.Colors{text.FgHiCyan, text.BgHiBlack}
	//t.Style().Color.IndexColumn = text.Colors{text.FgHiCyan, text.BgHiBlack}
	//t.Style().Color.Row = text.Colors{text.FgHiWhite, text.BgBlack}
	//t.Style().Color.RowAlternate = text.Colors{text.FgWhite, text.BgBlack}
	//t.Style().Color.Footer = text.Colors{text.FgCyan, text.BgHiBlack}
	//t.Style().Options.DrawBorder = false
	//t.Style().Options.SeparateHeader = false
	//t.Style().Options.SeparateRows = false
	//t.Style().Options.SeparateColumns = false
	//t.Style().Options.SeparateFooter = false

	return t
}

func getStyle(styleName string) table.Style {
	styles := map[string]table.Style{
		"DarkSimple":  DarkSimple,
		"LightSimple": LightSimple,
	}

	return styles[styleName]
}

var DarkSimple = table.Style{
	Name:    "DarkSimple",
	Box:     StyleBoxDefault,
	Color:   ColorOptionsCyanWhiteOnBlack,
	Format:  FormatOptionsDefault,
	Options: OptionsNoBordersAndSeparators,
	Title:   TitleOptionsBlackOnCyan,
}

var LightSimple = table.Style{
	Name:    "LightSimple",
	Box:     StyleBoxDefault,
	Color:   ColorOptionsBlackOnCyanWhite,
	Format:  FormatOptionsDefault,
	Options: OptionsNoBordersAndSeparators,
	Title:   TitleOptionsCyanOnBlack,
}

var StyleBoxDefault = table.BoxStyle{
	BottomLeft:       "+",
	BottomRight:      "+",
	BottomSeparator:  "+",
	EmptySeparator:   text.RepeatAndTrim(" ", text.RuneWidthWithoutEscSequences("+")),
	Left:             "|",
	LeftSeparator:    "+",
	MiddleHorizontal: "-",
	MiddleSeparator:  "+",
	MiddleVertical:   "|",
	PaddingLeft:      " ",
	PaddingRight:     " ",
	PageSeparator:    "\n",
	Right:            "|",
	RightSeparator:   "+",
	TopLeft:          "+",
	TopRight:         "+",
	TopSeparator:     "+",
	UnfinishedRow:    " ~",
}

var ColorOptionsCyanWhiteOnBlack = table.ColorOptions{
	Footer:       text.Colors{text.FgCyan, text.BgHiBlack},
	Header:       text.Colors{text.FgHiCyan, text.BgHiBlack},
	IndexColumn:  text.Colors{text.FgHiCyan, text.BgHiBlack},
	Row:          text.Colors{text.FgHiWhite, text.BgBlack},
	RowAlternate: text.Colors{text.FgWhite, text.BgBlack},
}

var ColorOptionsBlackOnCyanWhite = table.ColorOptions{
	Footer:       text.Colors{text.BgCyan, text.FgBlack},
	Header:       text.Colors{text.BgHiCyan, text.FgBlack},
	IndexColumn:  text.Colors{text.BgHiCyan, text.FgBlack},
	Row:          text.Colors{text.BgHiWhite, text.FgBlack},
	RowAlternate: text.Colors{text.BgWhite, text.FgBlack},
}

var FormatOptionsDefault = table.FormatOptions{
	Footer: text.FormatUpper,
	Header: text.FormatUpper,
	Row:    text.FormatDefault,
}

var OptionsNoBordersAndSeparators = table.Options{
	DrawBorder:      false,
	SeparateColumns: false,
	SeparateFooter:  false,
	SeparateHeader:  false,
	SeparateRows:    false,
}

var TitleOptionsBlackOnCyan = table.TitleOptions{
	Colors: append(table.ColorOptionsBlackOnCyanWhite.Header, text.Bold),
}

var TitleOptionsCyanOnBlack = table.TitleOptions{
	Colors: append(ColorOptionsCyanWhiteOnBlack.Header, text.Bold),
}
