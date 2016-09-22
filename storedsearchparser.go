//# Copyright 2015 by Pasquale J. Rinaldi, Jr.
//# Public Domain
//# If you encounter any issues, email me at pjrinaldi@gmail.com
//# Directions: the script looks for the met file in the same directory as the script and dumps the output to the same directory.

package main
import (
	"encoding/binary"
	"bytes"
	"strconv"
	"fmt"
	"os"
	"io"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NumberToString(n int, sep rune) string {

    s := strconv.Itoa(n)

    startOffset := 0
    var buff bytes.Buffer

    if n < 0 {
        startOffset = 1
        buff.WriteByte('-')
    }


    l := len(s)

    commaIndex := 3 - ((l - startOffset) % 3) 

    if (commaIndex == 3) {
        commaIndex = 0
    }

    for i := startOffset; i < l; i++ {

        if (commaIndex == 3) {
            buff.WriteRune(sep)
            commaIndex = 0
        }
        commaIndex++

        buff.WriteByte(s[i])
    }

    return buff.String()
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Print("./storedsearches infile outfile")
		os.Exit(0)
	}
	metfile, err := os.Open(args[0])
	check(err)
	outfile, err := os.Create(args[1])
	check(err)
	io.WriteString(outfile, "<html>\n<head>\n<style>")
	io.WriteString(outfile, ".tablehead { text-transform: uppercase; border-top: 1px solid black; border-bottom: 1px solid black;}\n")
	io.WriteString(outfile, ".oddrow { background-color: ddd; }\n")
	io.WriteString(outfile, ".endrow { border-top: 1px solid black;}\n")
	io.WriteString(outfile, ".cell { padding: 3 10 10 3; }\n")
	io.WriteString(outfile, "</style>\n</head>\n<body>\n")
	io.WriteString(outfile, "<h2>Parsed StoredSearches.met File</h2>\n")
	metfile.Seek(0, 0)
	var filenamestring string
	b1 := make([]byte, 1)
	b2 := make([]byte, 2)
	b4 := make([]byte, 4)
	bname10 := make([]byte, 10)
	bname11 := make([]byte, 11)
	bname12 := make([]byte, 12)
	bname13 := make([]byte, 13)
	bname14 := make([]byte, 14)
	bname15 := make([]byte, 15)
	bname16 := make([]byte, 16)
	var filesize uint32
	b16 := make([]byte, 16)
	metfile.Read(b1)
	if b1[0] != 0x0F {
		fmt.Print("Not a storedsearches.met file\n")
		os.Exit(0)
	}
	metfile.Read(b1)
	if b1[0] != 0x01 {
		fmt.Print("Not the right version\n")
		os.Exit(0)
	}
	metfile.Read(b2)
	searchcount := binary.LittleEndian.Uint16(b2)
	if searchcount <= 0 {
		fmt.Print("No Searches Found.\n")
		os.Exit(0)
	} else {
		io.WriteString(outfile, "<h3>Number of Open User Search Tabs: " + fmt.Sprintf("%d", searchcount) + "</h3>\n")
	}
	for i := 1; i <= int(searchcount); i++ {
		metfile.Seek(6, 1)
		metfile.Read(b2)
		titlelength := binary.LittleEndian.Uint16(b2)
		metfile.Seek(int64(titlelength), 1)
		metfile.Read(b2)
		exprlength := binary.LittleEndian.Uint16(b2)
		bexpr := make([]byte, exprlength)
		metfile.Read(bexpr)
		io.WriteString(outfile, "<h4>Search Expression " + fmt.Sprintf("%d", i) + ": \"" + fmt.Sprintf("%s", string(bexpr)) + "\" had ")
		metfile.Read(b2)
		typelength := binary.LittleEndian.Uint16(b2)
		metfile.Seek(int64(typelength), 1)
		metfile.Read(b4)
		hitcount := binary.LittleEndian.Uint32(b4)
		io.WriteString(outfile, fmt.Sprintf("%s", NumberToString(int(hitcount), ',')) + " search results returned.</h4>\n")
		io.WriteString(outfile, "<table style=\"border-collapse: collapse;\">\n")
		io.WriteString(outfile, "<tr class=\"tablehead\"><th>Hit</th><th>File name</th><th>File Size (bytes)</th><th>File Hash</th></tr>\n")
		for j := 1; j <= int(hitcount); j++ {
			metfile.Read(b16)
			metfile.Seek(6, 1)
			metfile.Read(b4)
			tagcount := binary.LittleEndian.Uint32(b4)
			for k := 1; k <= int(tagcount); k++ {
				metfile.Read(b1)
				if b1[0] == 0x82 {
					metfile.Seek(1, 1)
					metfile.Read(b2)
					filenamelength := binary.LittleEndian.Uint16(b2)
					bname := make([]byte, filenamelength)
					metfile.Read(bname)
					filenamestring = string(bname)
				}
				if b1[0] == 0x83 {
					metfile.Seek(1, 1)
					metfile.Read(b4)
					filesize = binary.LittleEndian.Uint32(b4)
				}
				if b1[0] == 0x89 {
					metfile.Seek(2, 1)
				}
				if b1[0] == 0x88 {
					metfile.Seek(3, 1)
				}
				if b1[0] == 0x94 {
					metfile.Seek(5, 1)
				}
				if b1[0] == 0x93 {
					metfile.Seek(4, 1)
				}
				if b1[0] == 0x92 {
					metfile.Seek(3, 1)
				}
				if b1[0] == 0x9C {
					metfile.Seek(1, 1)
					metfile.Read(bname12)
					filenamestring = string(bname12)
				}
				if b1[0] == 0x9E {
					metfile.Seek(1, 1)
					metfile.Read(bname14)
					filenamestring = string(bname14)
				}
				if b1[0] == 0x9D {
					metfile.Seek(1, 1)
					metfile.Read(bname13)
					filenamestring = string(bname13)
				}
				if b1[0] == 0x9F {
					metfile.Seek(1, 1)
					metfile.Read(bname15)
					filenamestring = string(bname15)
				}
				if b1[0] == 0x9A {
					metfile.Seek(1, 1)
					metfile.Read(bname10)
					filenamestring = string(bname10)
				}
				if b1[0] == 0x9B {
					metfile.Seek(1, 1)
					metfile.Read(bname11)
					filenamestring = string(bname11)
				}
				if b1[0] == 0x8B {
					metfile.Seek(9, 1)
				}
				if b1[0] == 0xA0 {
					metfile.Seek(1, 1)
					metfile.Read(bname16)
					filenamestring = string(bname16)
				}
				if j % 2 != 0 {
					io.WriteString(outfile, "<tr class=\"oddrow\">")
				} else {
					io.WriteString(outfile, "<tr>")
				}
			}
			io.WriteString(outfile, "<td align=\"center\" class=\"cell\">" + fmt.Sprintf("%d", j) + "</td><td class=\"cell\">" + fmt.Sprintf("%s", filenamestring) + "</td><td class=\"cell\" align=\"center\">" + fmt.Sprintf("%s", NumberToString(int(filesize), ',')) + "</td><td style=\"font-family: monospace;\" class=\"cell\" align=\"center\">" + fmt.Sprintf("%X", string(b16)) + "</td></tr>\n")
		}
		io.WriteString(outfile, "<tr><td colspan=\"4\" class=\"endrow\">&nbsp;</td></tr></table>\n")
		io.WriteString(outfile, "<br/><br/>\n")
	}
	io.WriteString(outfile, "</body></html>\n")
	fi, err := metfile.Stat()
	check(err)
	pos, err := metfile.Seek(0, 1)
	check(err)
	if fi.Size() == pos {
		fmt.Print("Stored Search Parsing completed successfully.\n")
	} else {
		fmt.Print("An error was encountered, but you may have some results to review.")
	}
	metfile.Close()
	outfile.Close()
}
