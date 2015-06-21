<center>
# mdserver
</center>

## OVERVIEW

mdserver reads a directory tree and picks out all files that are markdown documents.
These are then presented as filenames to user via browser.
When clicked, the file is run through the markdown process and shown to user.


### Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/mdserver```

Will copy github.com/hotei/mdserver to the first entry of your $GOPATH

or if go is not installed yet :

> ```cd DestinationDirectory```

> ```git clone https://github.com/hotei/mdserver.git```

### Configuration

* You can/should change the directory where the search for markdown docs begins.

### Features

* Serves files with the following extensions
 * .md
 * .mdown
 * .markdown
* Standalone webserver executable is about 6.5 MB with go 1.4.2
* supply a list of dirs if you want to exclude them from processing
* port number served is easily changed - default is 8281 (a random pick)
* can be run at "localhost" if appropriate (more secure, but can't be used from other computers)
* refreshes list every 5 minutes (setable) 

### Limitations

* <font color="red">If used in forensic environment you may need to adjust file permissions of the markdown files to allow the server to read all files.</font>

### Usage

Typical usage is :

```mdserver &```

### Bugs

None Known

### To-Do

* Essential:
 * TBD
* Nice:
 * TBD


### Change Log
* 2015-06-20 minor update pushed
	* skipDirs implemented
	* suppress generic duplicated README.md 
* 2015-04-30 revisited for go 1.4.2 update
* 2014-03-25 Working
* 2014-03-25 Started

### Notes



### Resources

* [go language reference docs] [1] 
* [go standard library package docs] [2]
* [Source code for the program] [3]

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go package docs"
[3]: http://github.com/hotei/mdserver "github.com/hotei/mdserver"

Comments can be sent to <hotei1352@gmail.com> or to user "hotei" at github.com.
License is BSD-two-clause, in file "LICENSE"

License
-------
The 'mdserver' go package/program is distributed under the Simplified BSD License:

> Copyright (c) 2014-2015 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Documentation (c) 2015 David Rook 

// EOF README-mdserver.md  (this markdown document tested OK with blackfriday)
