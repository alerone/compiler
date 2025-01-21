package emitter

import (
	"fmt"
	"os"
)

type Emitter struct {
    // FullPath to output file.
	FullPath string
    // Contains things that we will prepend to the code later on.
	Header   string
    // String containing the C code that is emitted.
	Code     string
}

func NewEmitter(fullPath string) Emitter{
    return Emitter{FullPath: fullPath, Header: "", Code: ""}    
}

// Add a fragment of C code.
func (self *Emitter) Emit(code string) {
	self.Code += code
}

// Add a fragment of C code that ends a line.
func (self *Emitter) EmitLine(code string) {
	self.Code += code + "\n"
}

// Adding a line of C code to the top of the C code file (including library, the main function, variable declarations,...).
func (self *Emitter) HeaderLine(code string) {
	self.Header += code + "\n"
}

// Writes the C code to a file.
func (self *Emitter) WriteFile() {
    fileOutput, err := os.Create(self.FullPath)
    if err != nil {
        panic(fmt.Sprintf("Error. can't open/create output file: %v", self.FullPath))
    }
    _, err = fileOutput.Write([]byte(self.Header + self.Code))  
    if err != nil {
        panic(fmt.Sprintf("Error. can't write to output file: %v", self.FullPath))
    }
}
