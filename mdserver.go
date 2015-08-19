// mdserver.go (c) 2013-2015 David Rook - all rights reserved
//
// Serve markdown files over http
//
// TODO(mdr) need a way to watch fs for changes before re-starting loadFiles()
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	//
	"github.com/russross/blackfriday"
)

const (
	serverRoot = "/home/mdr/Desktop/GO/"
	imageRoot  = "/home/mdr/Desktop/webbies/"
	mdURL      = "/md/"
	imageURL   = "/images/"
	version    = "mdserver version 0.0.2 (c) 2015 David Rook"
)

var (
	flagServerIPStr  string
	flagServerPort   int
	flagRefreshDelay int // in seconds
	flagLocalHost    bool
	flagCommonMkdn   bool = true // if true then use Common version vs Basic
	flagBasicMkdn    bool        // if true then use Common version vs Basic
	flagVerbose      bool
	flagVersion      bool

	portNum       = 8281
	portNumString string
	hostIPstr     string
	wantLocal     = true
	listenOnStr   string

	nFiles           int
	g_fileNames      []string // files with md content
	loadingFilenames sync.Mutex
	delayReloadSecs  time.Duration //  = 300 to reload every 5 minutes, set from flag

	myMdDir  = []byte{}
	pathName string
)

// suppress generic label README.md in these dirs
// github.com wants a README.md but on my system these are just links
// to a more specific README-programName.md file
var myGOs []string = []string{"GoGit", "GoHub", "GoWork", "GoDoc"}

// skip these entirely
var skipDirs []string = []string{
	"/home/mdr/Desktop/GO/GoWork/src/hubmd/tests/",
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	flag.StringVar(&flagServerIPStr, "hostIP", "10.1.2.213", "host IP address")
	flag.IntVar(&flagServerPort, "port", portNum, "host port number")
	flag.IntVar(&flagRefreshDelay, "refresh", 300, "delay between refresh of index")
	flag.BoolVar(&flagLocalHost, "localhost", false, "serve on local host 127.0.0.1")
	flag.BoolVar(&flagCommonMkdn, "common", true, "Use common markddown styles")
	flag.BoolVar(&flagBasicMkdn, "basic", false, "Use basic markddown styles")
	flag.BoolVar(&flagVerbose, "verbose", false, "Use more messages to user")
	flag.BoolVar(&flagVersion, "version", false, "print version and exit")
}

func flagSetup() {
	if !flag.Parsed() {
		flag.Parse()
	}
	Verbose = VerboseType(flagVerbose)
	if flagCommonMkdn && flagBasicMkdn {
		fmt.Printf("Can't use both common and basic markdown flags at same time\n")
	}
	if flagLocalHost {
		hostIPstr = "127.0.0.1"
	} else {
		if len(flagServerIPStr) > 0 {
			hostIPstr = flagServerIPStr
		}
	}
	if flagServerPort > 1024 {
		portNum = flagServerPort
	}
	portNumString = fmt.Sprintf(":%d", portNum)
	listenOnStr = hostIPstr + portNumString
	Verbose.Printf("Will listen on IP:port %s\n", listenOnStr)
	var mkdnType string = "unknown"
	if flagCommonMkdn {
		mkdnType = "common"
	}
	if flagBasicMkdn {
		mkdnType = "basic"
	}
	Verbose.Printf("Will use %s markdown\n", mkdnType)
	if flagRefreshDelay <= 30 {
		flagRefreshDelay = 30
	}
	delayReloadSecs = time.Duration(flagRefreshDelay)
	Verbose.Printf("delay between refreshes will be %d seconds\n", flagRefreshDelay)
	if flagVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}
	checkInterfaces()
}

// loadFiles() will run through the serverRoot directory and build an index
// of all markdown type files.  This runs as a goroutine that never exits.
// it will sleep a fixed time after building an index (normally about 5 min)
// this should be replaced by a select on chan of int to trigger rebuild
// somehow need to know when dir tree has been modified.  fsnotify/fswatch or ?
func loadFiles() {
	pathName := serverRoot
	for {
		nFiles = 0
		// dont serve while building the index
		loadingFilenames.Lock()
		g_fileNames = make([]string, 0, 1000)
		// 8888 TODO ??? next update at ---
		myMdDir = []byte(`<html><!-- comment --><head><title>Test MD package</title>
			</head><body>click link to read<br>refresh page if files added to server but
			note list update may lag by up to 5 minutes<br><p>`)
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
		Verbose.Printf("g_fileNames = %v\n", g_fileNames)
		for ndx, val := range g_fileNames {
			//fmt.Printf("%v\n", val)
			nFiles++
			line := makeMdLine(ndx, val)
			myMdDir = append(myMdDir, line...)
		}
		t := []byte(`</body></html>`)
		myMdDir = append(myMdDir, t...)
		fmt.Printf("Loaded files from directory %s and found %d files to serve\n",
			serverRoot, nFiles)
		loadingFilenames.Unlock()
		fmt.Printf("Sleeping for %d seconds\n", int(delayReloadSecs))
		time.Sleep(delayReloadSecs * time.Second)
	}
}

// checkMdName treewalk the path and find files that match markdown extensions
// these file names get appended to g_fileNames as a side effect.
func checkMdName(pathname string, info os.FileInfo, err error) error {
	Verbose.Printf("checking %s\n", pathname)
	if info == nil {
		fmt.Printf("WARNING --->  no stat info: %s\n", pathname)
		os.Exit(1)
	}
	if info.IsDir() {
		// if dir is in our skiplist then skip it
		for i := 0; i < len(skipDirs); i++ {
			if len(pathname) >= len(skipDirs[i]) {
				if pathname[:len(skipDirs[i])] == skipDirs[i] {
					return filepath.SkipDir
				}
			}
		}
		return nil
	} else {
		// TODO(mdr) should test if it's really a regular file
		//log.Printf("found %s %s\n", pathname, filepath.Ext(pathname))
		ext := filepath.Ext(pathname)
		if ext == ".md" || ext == ".markdown" || ext == ".mdown" {
			//log.Printf("basename = %s\n", filepath.Base(pathname))

			/* originally used to ignore README.md since it was a copy of
			something else.  No longer true.  However, if you were going to
			ignore something, this would be the place to test for a regex match
			or whatever...

			if filepath.Base(pathname) == "README.md" {
				for _, v := range myGOs {
					if strings.Contains(pathname, v) {
						return nil
					}
				}
			}
			*/
			g_fileNames = append(g_fileNames, pathname)
		}
	}
	return nil
}

// makeMdLine create one line of the index file.
func makeMdLine(i int, s string) []byte {
	//workDir := serverRoot + mdURL[1:]
	s = s[len(serverRoot):]
	x := fmt.Sprintf("%d <a href=\"%s\">%s</a><br>", i, mdURL+s, s)
	Verbose.Printf("line: %s\n", x)
	return []byte(x)
}

// checkInterfaces - see if listener is bound to correct interface.
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
	fmt.Printf("ifa[] = %v\n", ifa)
	var myIfs []string
	if wantLocal {
		myIfs = strings.Split(ifa[0].String(), "/")
	} else {
		myIfs = strings.Split(ifa[1].String(), "/")
	}
	fmt.Printf("myIfs = %v\n", myIfs)
	myIf := myIfs[0]
	fmt.Printf("myIf = %v\n", myIf)
	if myIf != hostIPstr {
		log.Printf("Wanted %s got %s\n", hostIPstr, myIf)
		log.Fatalf("handler bound to wrong interface")

	}
}

// mdHandler recognizes markdown extensions and expands to html.
func mdHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == mdURL {
		loadingFilenames.Lock()
		w.Write(myMdDir)
		loadingFilenames.Unlock()
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

// htmlFromMd creates html from a markdown style document.
func htmlFromMd(fname string) []byte {
	var output []byte
	input, err := ioutil.ReadFile(fname)
	if err != nil {
		tmp := fmt.Sprintf("Problem reading input, can't open %s", fname)
		output = []byte(tmp)
	} else {
		if flagCommonMkdn { // what's different between these?
			output = blackfriday.MarkdownCommon(input)
		}
		if flagBasicMkdn {
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
	flagSetup()
	go loadFiles()
	http.HandleFunc(mdURL, mdHandler)
	http.HandleFunc(imageURL, imageHandler)
	log.Printf("Compiled on %s\n", CompileDateTime)
	log.Printf("Version = %s\n", version)
	log.Printf("Server root = %s\n", serverRoot)
	log.Printf("Image root = %s\n", imageRoot)
	log.Printf("image urls syntax is: http://:%s/images/x.png for example\n", listenOnStr)
	log.Printf("md server is ready at %s\n", listenOnStr)
	log.Printf("start browser with this url: %s%s\n", listenOnStr, mdURL)
	err := http.ListenAndServe(listenOnStr, nil)
	if err != nil {
		log.Printf("mdserver: error running webserver %v", err)
	}
}
