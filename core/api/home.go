package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ossman11/sip/core/index"
)

const (
	// HomePageContent The text to be displayed on the home page
	HomePageContent string = `
<html>

<head>
	<meta charset="utf-8">
	<script src="wasm.js"></script>
	<script>
		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
			go.run(result.instance);
		});
	</script>
</head>

<body style="margin:0px"></body>

</html>
`
)

// Home the Api interface implementation for the Home Api
type Home struct{}

func NewHome() API {
	return Home{}
}

// Action Implements the Home Api behavior
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, HomePageContent)
}

func wasmJS(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./wasm_exec.js")

	if err != nil {

		file, err = os.Open(filepath.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js"))
		if err != nil {
			http.Error(w, "Failed to find wasm.js file", 404)
			return
		}
	}

	defer file.Close()
	io.Copy(w, file)
}

func wasmGO(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("main.wasm")

	if err != nil {
		err := index.Build("js", "wasm")
		if err != nil {
			http.Error(w, "Failed to compile 'main.wasm' file.", 500)
			return
		}

		tmpFile, err := os.Open(".tmp/js-wasm")
		if err != nil {
			http.Error(w, "Failed to open 'main.wasm' file.", 500)
			return
		}
		defer tmpFile.Close()

		file, err = os.Create("main.wasm")
		if err != nil {
			http.Error(w, "Failed to create 'main.wasm' file.", 500)
			return
		}

		io.Copy(file, tmpFile)
		file.Close()

		file, err = os.Open("main.wasm")
		if err != nil {
			http.Error(w, "Failed to open 'main.wasm' file.", 500)
			return
		}
	}
	defer file.Close()
	io.Copy(w, file)
}

// Get Implements the Get API for the Home definition
func (h Home) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/":          home,
		"/wasm.js":   wasmJS,
		"/main.wasm": wasmGO,
	}
}

// Post Implements the Post API for the Home definition
func (h Home) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}

func (h Home) Running() func() {
	return func() {}
}
