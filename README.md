nurikabe
========
This project is a qml graphical UI for playing the logic puzzle nurikabe.

Requirements
------------
* golang >= 1.3

    To install golang visit: https://golang.org/doc/install

* Qt >= 5

    To install Qt visit: http://qt-project.org/downloads

* go-qml

    run 'go get github.com/gopkg.in/qml.v1'
    for documentation visit: http://github.com/go-qml/qml

Building
--------
Once all requirements have been met, you should be able to run 'go build' from the command line.
This will build a binary which you can then execute. Note that you must run the binary in the
same directory as the qml folder.

Levels
----
Nurikabe uses json format for all its levels. You may also generate levels using the nurikabe/gen helper binary.
The gen utility also allows for solving levels by piping the json level via stdin and issuing the 'solve' flag.

    ie. cat my_level.json | gen -solve

    Usage of ./gen:
      -base=2: minimum garden size
      -debug=false: enable debug output
      -growth=4: garden growth. base + growth is max garden size
      -height=5: grid height
      -min=3: minimum gardens count
      -smart=true: solve using smart algorithm
      -solve=false: solve generated grid
      -v=false: Verbose
      -width=5: grid width
