// mdserver.go (c) 2013-2015 David Rook - all rights reserved

package main

import (
	// go 1.4.2 std lib
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// non-local pkgs
	"github.com/russross/blackfriday"
)

const (
	portNum    = 8281
	wantLocal  = true
	hostIPstr = "127.0.0.1"
	// hostIPstr = "10.1.2.113" // loki is 112, mars is 113
	serverRoot = "/home/mdr/Desktop/GO/GoDoc/"
	mdURL      = "/md/"
)

var (
	portNumString = fmt.Sprintf(":%d", portNum)
	g_fileNames   []string // files with md content
	listenOnPort  = hostIPstr + portNumString
	nFiles int
)

var myMdDir = []byte{}

var pathName string

func checkMdName(pathname string, info os.FileInfo, err error) error {
	fmt.Printf("checking %s\n", pathname)
	if info == nil {
		fmt.Printf("WARNING --->  no stat info: %s\n", pathname)
		os.Exit(1)
	}
	if info.IsDir() {
		// return filepath.SkipDir
		// g_fileNames = append(g_fileNames, pathname)
		return nil
	} else { // regular file
		//fmt.Printf("found %s %s\n", pathname, filepath.Ext(pathname))
		ext := filepath.Ext(pathname)
		if ext == ".md" || ext == ".markdown" || ext == ".mdown" {
			//fmt.Printf("appending\n")
			g_fileNames = append(g_fileNames, pathname)
		}
	}
	return nil
}

func makeMdLine(i int, s string) []byte {
	//workDir := serverRoot + mdURL[1:]
	s = s[len(serverRoot):]
	x := fmt.Sprintf("%d <a href=\"%s\">%s</a><br>", i, mdURL+s, s)
	fmt.Printf("line: %s\n",x)
	nFiles++
	return []byte(x)
}

func init() {
	checkInterfaces()
	pathName := serverRoot
	g_fileNames = make([]string, 0, 20)
	myMdDir = []byte(`<html><!-- comment --><head><title>Test MD package</title></head><body>click to read<br>`) // {}
	stats, err := os.Stat(pathName)
	if err != nil {
		fmt.Printf("Can't get fileinfo for %s\n", pathName)
		os.Exit(1)
	}
	if stats.IsDir() {
		filepath.Walk(pathName, checkMdName)
	} else {
		fmt.Printf("this argument must be a directory (but %s isn't)\n", pathName)
		os.Exit(-1)
	}
	fmt.Printf("g_fileNames = %v\n", g_fileNames)
	for ndx, val := range g_fileNames {
		//fmt.Printf("%v\n", val)
		line := makeMdLine(ndx, val)
		myMdDir = append(myMdDir, line...)
	}
	t := []byte(`</body></html>`)
	myMdDir = append(myMdDir, t...)
	fmt.Printf("Init ran oks, found %d files to serve\n",nFiles)
}

// checkInterfaces - see if listener is bound to correct interface
// first is localhost, second should be IP4 of active card,
// third is IP6 localhost, fourth is IP6 for active card (on this system)
func checkInterfaces() {
	ifa, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("Can't list interfaces\n")
		os.Exit(1)
	}
	fmt.Printf("Interfaces (ifa[]) = %v\n", ifa)
	if len(ifa) < 2 {
		fmt.Printf("Can't list interfaces\n")
		os.Exit(1)
	}
	// check IP4 of active card
	fmt.Printf("ifa[] = %v\n",ifa)
	var myIfs []string
	if wantLocal {
	myIfs = strings.Split(ifa[0].String(), "/")
	} else {
	myIfs = strings.Split(ifa[1].String(), "/")
	}
	fmt.Printf("myIfs = %v\n",myIfs)
	myIf := myIfs[0]
	fmt.Printf("myIf = %v\n",myIf)
	if myIf != hostIPstr {
		log.Fatalf("handler bound to wrong interface")
	}
}

// mdHandler recognizes markdown extensions and expands to html
func mdHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == mdURL {
		w.Write(myMdDir)
		return
	}
	var output []byte
	var err error
	fileName := serverRoot + r.URL.Path[len(mdURL):]
	fmt.Printf("mdHandler: reading fname = %s\n", fileName)
	ext := filepath.Ext(fileName)
	if ext == ".md" || ext == ".markdown" || ext == ".mdown" {
		output = htmlFromMd(fileName)
		w.Write(output)
		return
	}
	fmt.Printf("%s isn't markdown file type\n", fileName)
	// if in this path but not md - then send to browser is naieve
	// not a markdown ext - what is it? shouldn't happen
	output, err = ioutil.ReadFile(fileName)
	if err != nil {
		errStr := fmt.Sprintf("mdHandler: %v is not readable\n", err)
		fmt.Printf("%s\n", errStr)
		w.Write([]byte(fmt.Sprintf("404 - Not Found\n")))
		return
	}
	w.Write(output)
}

// htmlFromMd creates html from a markdown style document
func htmlFromMd(fname string) []byte {
	var output []byte
	input, err := ioutil.ReadFile(fname)
	if err != nil {
		tmp := fmt.Sprintf("Problem reading input, can't open %s", fname)
		output = []byte(tmp)
	} else {
		if true { // what's different between these?
			output = blackfriday.MarkdownCommon(input)
		} else {
			output = blackfriday.MarkdownBasic(input)
		}
	}
	if false { // debug use only
		os.Stdout.Write(input)
		os.Stdout.Write(output)
	}
	return output
}

func main() {
	//	http.HandleFunc(virtualURL, html)
	// Handle(serverRoot, is like a dir missing an index "ftp-style"
	//http.Handle(serverRoot, http.StripPrefix(serverRoot, http.FileServer(http.Dir(serverRoot))))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	http.HandleFunc(mdURL, mdHandler)
	log.Printf("md server is ready at %s\n", listenOnPort)
	log.Printf("start browser with this url: %s%s\n",listenOnPort,mdURL)
	err := http.ListenAndServe(listenOnPort, nil)
	if err != nil {
		log.Printf("mdserver: error running webserver %v", err)
	}
}
