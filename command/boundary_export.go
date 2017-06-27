package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"

	"github.com/codegangsta/cli"
	"github.com/missinglink/gosmparse"
)

// BoundaryExporter cli command
func BoundaryExporter(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 2 {
		log.Println("invalid arguments, expected: {leveldb} {geojson_dir}")
		os.Exit(1)
	}

	// stat leveldb destination
	lib.EnsureDirectoryExists(argv[0], "leveldb")

	// stat geojson destination
	lib.EnsureDirectoryExists(argv[1], "geojson")

	// open database connection
	conn := &leveldb.Connection{}
	conn.Open(argv[0])
	defer conn.Close()

	// worker function
	var worker = func(rel *gosmparse.Relation) {

		// create a new assembler
		var assembler = &lib.RelationAssembler{
			Relation: rel,
			Conn:     conn,
		}

		// generate json
		var json = assembler.GenerateJSON()

		// child process
		var child *exec.Cmd

		// increase v8 max memory limit to 8GB for json over 100MB
		// see: https://github.com/tyrasd/osmtogeojson#usage
		if len(json.Bytes()) > 104857600 /*(100 * 1024 * 1024)*/ {
			child = exec.Command("/usr/local/bin/node", "--max_old_space_size=8192", "/home/peter/.go/src/github.com/missinglink/pbf/nodejs/osmtogeojson.js")
		} else {
			child = exec.Command("/usr/local/bin/node", "/home/peter/.go/src/github.com/missinglink/pbf/nodejs/osmtogeojson.js")
		}

		// stdio
		stdin, _ := child.StdinPipe()
		stdout, _ := child.StdoutPipe()
		stderr, _ := child.StderrPipe()

		// start process
		if err := child.Start(); err != nil {
			log.Println("An error occured: ", err)
		}

		// write to stdin
		stdin.Write(json.Bytes())
		stdin.Close()

		// read from stdio
		stdoutBytes, _ := ioutil.ReadAll(stdout)
		stderrBytes, _ := ioutil.ReadAll(stderr)

		// wait for child to exit
		child.Wait()
		stdout.Close()
		stderr.Close()

		// debug stderr messages
		if len(stderrBytes) > 0 {
			log.Println("relation", rel.ID)
			log.Println(string(stderrBytes))
		}

		// pad id with leading zeros
		var id = fmt.Sprintf("%09d", rel.ID)
		var dir = fmt.Sprintf("%s/%s/%s/%s/", argv[1], id[0:3], id[3:6], id[6:9])

		// create directory if it doesn't exist
		os.MkdirAll(dir, 0777)

		// write geojson to disk (on success only)
		if len(stdoutBytes) > 0 {
			ioutil.WriteFile(fmt.Sprintf("%s/%s.geojson", dir, id), stdoutBytes, 0644)
		}

		// write errors and inputs to disk (on error only)
		if len(stderrBytes) > 0 {
			ioutil.WriteFile(fmt.Sprintf("%s/%s.in", dir, id), json.Bytes(), 0644)
			ioutil.WriteFile(fmt.Sprintf("%s/%s.err", dir, id), stderrBytes, 0644)
		}
	}

	// create a channel for relations
	var queue = make(chan *gosmparse.Relation, 256)

	// create waitgroup for routines
	var wg = &sync.WaitGroup{}

	// total amount of go routines to use
	// note: each one will spawn a separate node process
	var maxRoutines = runtime.NumCPU() - 1

	// use multiple goroutes
	for i := 0; i < maxRoutines; i++ {
		wg.Add(1)
		go func() {
			for rel := range queue {
				worker(rel)
			}
			wg.Done()
		}()
	}

	// iterate over relations, add each to the queue
	go func() {

		// iterate over relations, add each to the queue
		conn.IterateRelations(func(rel *gosmparse.Relation, err error) {
			queue <- rel
		})

		// close queue
		close(queue)
	}()

	// wait for all routines to finish
	wg.Wait()

	return nil
}
