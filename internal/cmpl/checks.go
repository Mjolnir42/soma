package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// I'm sorry as well.
func CheckConfigCreate(c *cli.Context) {
	topArgs := []string{`in`, `on`, `with`, `interval`, `inheritance`, `childrenonly`, `extern`, `threshold`, `constraint`}
	thrArgs := []string{`predicate`, `level`, `value`}
	ctrArgs := []string{`service`, `oncall`, `attribute`, `system`, `native`, `custom`}
	onArgs := []string{`repository`, `bucket`, `group`, `cluster`, `node`}

	if c.NArg() == 0 {
		return
	}

	if c.NArg() == 1 {
		for _, t := range topArgs {
			fmt.Println(t)
		}
	}

	skipNext := 0
	subON := false
	subTHRESHOLD := false
	subCONSTRAINT := false

	hasIN := false
	hasON := false
	hasWITH := false
	hasINTERVAL := false
	hasINHERITANCE := false
	hasCHILDRENONLY := false
	hasEXTERN := false

	hasTHRPredicate := false
	hasTHRLevel := false
	hasTHRValue := false

	hasCTRService := false
	hasCTROncall := false
	hasCTRAttribute := false
	hasCTRSystem := false
	hasCTRNative := false
	hasCTRCustom := false
	hasCTRSelectedService := false
	hasCTRSelectedOncall := false

	for _, t := range c.Args().Tail() {
		if skipNext > 0 {
			skipNext--
			continue
		}
		if subON {
			skipNext = 1
			subON = false
		}
		if subTHRESHOLD {
			if hasTHRPredicate && hasTHRLevel && hasTHRValue {
				subTHRESHOLD = false
				hasTHRPredicate = false
				hasTHRLevel = false
				hasTHRValue = false
			} else {
				switch t {
				case `predicate`:
					skipNext = 1
					hasTHRPredicate = true
					continue
				case `level`:
					skipNext = 1
					hasTHRLevel = true
					continue
				case `value`:
					skipNext = 1
					hasTHRValue = true
					continue
				}
			}
		}
		if subCONSTRAINT {
			if hasCTRSelectedService {
				skipNext = 1
				hasCTRSelectedService = false
				continue
			}
			if hasCTRSelectedOncall {
				skipNext = 1
				hasCTRSelectedOncall = false
				continue
			}
			if hasCTRService || hasCTROncall || hasCTRAttribute || hasCTRSystem || hasCTRNative || hasCTRCustom {
				subCONSTRAINT = false
				hasCTRService = false
				hasCTROncall = false
				hasCTRAttribute = false
				hasCTRSystem = false
				hasCTRNative = false
				hasCTRCustom = false
				hasCTRSelectedService = false
				hasCTRSelectedOncall = false
			} else {
				switch t {
				case `service`:
					hasCTRSelectedService = true
					hasCTRService = true
					continue
				case `oncall`:
					hasCTRSelectedOncall = true
					hasCTROncall = true
					continue
				case `attribute`:
					skipNext = 2
					hasCTRAttribute = true
					continue
				case `system`:
					skipNext = 2
					hasCTRSystem = true
					continue
				case `native`:
					skipNext = 2
					hasCTRNative = true
					continue
				case `custom`:
					skipNext = 2
					hasCTRCustom = true
					continue
				}
			}

		}
		switch t {
		case `in`:
			skipNext = 1
			hasIN = true
			continue
		case `on`:
			hasON = true
			subON = true
			continue
		case `with`:
			skipNext = 1
			hasWITH = true
			continue
		case `interval`:
			skipNext = 1
			hasINTERVAL = true
			continue
		case `inheritance`:
			skipNext = 1
			hasINHERITANCE = true
			continue
		case `childrenonly`:
			skipNext = 1
			hasCHILDRENONLY = true
			continue
		case `extern`:
			skipNext = 1
			hasEXTERN = true
			continue
		case `threshold`:
			subTHRESHOLD = true
			continue
		case `constraint`:
			subCONSTRAINT = true
			continue
		}
	}
	// skipNext not yet consumed
	if skipNext > 0 {
		return
	}
	// in subchain: ON
	if subON {
		for _, t := range onArgs {
			fmt.Println(t)
		}
		return
	}
	// in subchain: CONSTRAINT
	if subCONSTRAINT {
		if hasCTRSelectedService || hasCTRSelectedOncall {
			fmt.Println(`name`)
			return
		}
		if !(hasCTRService || hasCTROncall || hasCTRAttribute || hasCTRSystem || hasCTRNative || hasCTRCustom) {
			for _, t := range ctrArgs {
				fmt.Println(t)
			}
			return
		}
	}
	// in subchain: THRESHOLD
	if subTHRESHOLD {
		if !(hasTHRPredicate && hasTHRLevel && hasTHRValue) {
			for _, t := range thrArgs {
				switch t {
				case `predicate`:
					if !hasTHRPredicate {
						fmt.Println(t)
					}
				case `level`:
					if !hasTHRLevel {
						fmt.Println(t)
					}
				case `value`:
					if !hasTHRValue {
						fmt.Println(t)
					}
				}
			}
			return
		}
	}
	// not in any subchain
	for _, t := range topArgs {
		switch t {
		case `in`:
			if !hasIN {
				fmt.Println(t)
			}
		case `on`:
			if !hasON {
				fmt.Println(t)
			}
		case `with`:
			if !hasWITH {
				fmt.Println(t)
			}
		case `interval`:
			if !hasINTERVAL {
				fmt.Println(t)
			}
		case `inheritance`:
			if !hasINHERITANCE {
				fmt.Println(t)
			}
		case `childrenonly`:
			if !hasCHILDRENONLY {
				fmt.Println(t)
			}
		case `extern`:
			if !hasEXTERN {
				fmt.Println(t)
			}
		default:
			fmt.Println(t)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
