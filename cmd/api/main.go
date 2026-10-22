// package main tells the Go compiler: "this file compiles into an executable binary,
// not an importable library." In TypeScript/Node, this is like your top-level index.ts.
package main

// import brings in packages from the Go standard library.
// "fmt" stands for "format" (I/O formatting).
// In Node/TypeScript, this is like having built-in console / process.stdout utilities.
import "fmt"

// main() is the universal entry point where the operating system starts execution.
// Unlike Node/Express where top-level code runs immediately on file load,
// Go requires an explicit func main() inside package main.
func main() {
	// fmt.Println writes a string followed by a newline to standard output (stdout).
	// This is our simplest log output, exactly like console.log() in JavaScript.
	// In v0.14, we will upgrade this to structured JSON logging (like Pino/Winston).
	fmt.Println("Hookline API v0.0 starting...")
}
