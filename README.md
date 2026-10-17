Spreadsheet Calculator
========
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/pedroalbanese/gosc/blob/master/LICENSE.md) 
[![GoDoc](https://godoc.org/github.com/pedroalbanese/gosc?status.png)](http://godoc.org/github.com/pedroalbanese/gosc)
[![GitHub downloads](https://img.shields.io/github/downloads/pedroalbanese/gosc/total.svg?logo=github&logoColor=white)](https://github.com/pedroalbanese/gosc/releases)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pedroalbanese/gosc)](https://golang.org)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/pedroalbanese/gosc)](https://github.com/pedroalbanese/gosc/releases)

This is a text-based spreadsheet program compatible with [sc](http://www.ibiblio.org/pub/Linux/apps/financial/spreadsheet/sc-7.16.lsm) files. `sc` is a public domain curses-based speadsheet calculator program that has similar key bindings to vi and works on Unix-like OSs.  

This version of `sc` is cross-platform, allowing it to view/edit sc files on a Windows machine.

## Supported Features

* Vi-like key movements
* Expressions (B23+E23*5)
* Export to file, as in sc.

## Normal Mode Commands

* Movement: h, j, k, l
* Enter insert mode for numbers and expressions: =, i
* Enter insert mode for strings with alignment: < (Left-Aligned), > (Right-Aligned), \ (Centered)
* Delete (clear) cell: x
* Commands: (y)ank, (p)ut
* Set output format for current column: f
* Quit: q

## TODO
* Undo `u`
* Goto `g

## Contribute
**Use issues for everything**
- You can help and get help by:
  - Reporting doubts and questions
- You can contribute by:
  - Reporting issues
  - Suggesting new features or enhancements
  - Improve/fix documentation

## License

This project is licensed under the ISC License.

#### Copyright (c) 2020-2026 Pedro F. Albanese - ALBANESE Research Lab.  
Todos os direitos de propriedade intelectual sobre este software pertencem ao autor, Pedro F. Albanese. Vide Lei 9.610/98, Art. 7º, inciso XII.
