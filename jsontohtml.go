// jsontohtml
package main

import (
	"fmt"
	"encoding/json"
	"os"
	"io"
	"io/ioutil"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseJSON(node map[string]interface{}, ofile* os.File, counter *int) {
	if *counter % 2 != 0 {
		io.WriteString(ofile, "<tr class=\"oddrow\">")
	} else {
		io.WriteString(ofile, "<tr>")
	}
	
	var curannos map[string]interface{}
	if node["type"].(string) == "text/x-moz-place" {
		if node["title"] != nil {
			io.WriteString(ofile, "<td class=\"cell\">" + node["title"].(string) + "</td>")
		}
		if node["uri"] != nil {
			io.WriteString(ofile, "<td width=\"200px\" class=\"cell\">" + node["uri"].(string) + "</td>")
		}
		if node["dateAdded"] != nil {
			dtime := time.Unix(int64(node["dateAdded"].(float64)/1000000), 0)
			io.WriteString(ofile, "<td width=\"350px\" class=\"cell\" align=\"center\">" + dtime.Format("01/02/2006 03:04:05 PM -0700 MST") + "</td>")
		}
		if node["lastModified"] != nil {
			ltime := time.Unix(int64(node["lastModified"].(float64)/1000000), 0)
			io.WriteString(ofile, "<td class=\"cell\" align=\"center\">" + ltime.Format("01/02/2006 03:04:05 PM -0700 MST") + "</td>")
		}
		if node["annos"] != nil {
			annos := node["annos"].([]interface{})
			if len(annos) > 0 {
				for j := 0; j < len(annos); j++ {
					curannos = annos[j].(map[string]interface{})
					if curannos["name"].(string) == "bookmarkProperties/description" {
						if curannos["value"] != nil {
							io.WriteString(ofile, "<td class=\"cell\">" + curannos["value"].(string) + "</td></tr>\n")
						} else {
							io.WriteString(ofile, "<td class=\"cell\">&nbsp;</td></tr>\n")
						}
					} else {
						io.WriteString(ofile, "<td class=\"cell\">&nbsp;</td></tr>\n")
					}
				}
			} else {
				io.WriteString(ofile, "<td class=\"cell\">&nbsp;</td></tr>\n")
			}
		} else {
			io.WriteString(ofile, "<td class=\"cell\">&nbsp;</td></tr>\n")
		}
		*counter++
	}
	if node["children"] != nil {
		tmpchild := node["children"].([]interface{})
		if len(tmpchild) > 0 {
			for i := 0; i < len(tmpchild); i++ {
				ParseJSON(tmpchild[i].(map[string]interface{}), ofile, counter)
			}
		}
	}
}

func main() {
	var cnt int
	cnt = 1
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Print("./jsontohtml infile outfile\n")
		os.Exit(0)
	}
	infile, err := ioutil.ReadFile(args[0])
	check(err)
	outfile, err := os.Create(args[1])
	check(err)
	io.WriteString(outfile, "<html>\n<head>\n<style>\n")
	io.WriteString(outfile, ".tablehead { text-transform: uppercase; border-top: 1px solid black; border-bottom: 1px solid black;}\n")
	io.WriteString(outfile, ".oddrow { background-color: ddd; }\n")
	io.WriteString(outfile, ".endrow { border-top: 1px solid black;}\n")
	io.WriteString(outfile, ".cell { padding: 3 10 10 3; word-wrap: break-word; }\n")
	io.WriteString(outfile, "</style>\n</head>\n<body>\n")
	io.WriteString(outfile, "<h2>Parsed Bookmarks Backup JSON File</h2>\n")
	io.WriteString(outfile, "<table style=\"border-collapse: collapse; table-layout: fixed; empty-cells: show;\" width=\"100%\">\n")
	io.WriteString(outfile, "<tr class=\"tablehead\"><th>Title</th><th>URI</th><th>Date Added</th><th>Last Modified</th><th>Description</th></tr>\n")
	
	var dat map[string]interface{}
	
	err2 := json.Unmarshal(infile, &dat)
	check(err2)
	ParseJSON(dat, outfile, &cnt)
	io.WriteString(outfile, "<tr><td colspan=\"5\" class=\"endrow\">&nbsp;</td></tr></table>\n")
	io.WriteString(outfile, "<br/><br/>\n")
	io.WriteString(outfile, "</body></html>\n")
	outfile.Close()
	fmt.Print("JSON File was successfully Parsed.\n")
}
