package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":5555", "server address:port")
	flag.Parse()
	http.HandleFunc("/", rootHandle)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, page)
}

const page = `
<html>
	<head>
		<title>Testing WebAssembly</title>
		<script type="text/javascript">

		var mod = Wasm.instantiateModule(
			new Uint8Array([
				0x00, 0x61, 0x73, 0x6d, 0x0b, 0x00, 0x00, 0x00,
				0x04, 0x74, 0x79, 0x70, 0x65, 0x07, 0x01, 0x40,
				0x02, 0x01, 0x01, 0x01, 0x01, 0x08, 0x66, 0x75,
				0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x02, 0x01,
				0x00, 0x06, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
				0x06, 0x01, 0x00, 0x03, 0x73, 0x75, 0x6d, 0x04,
				0x63, 0x6f, 0x64, 0x65, 0x0a, 0x01, 0x08, 0x00,
				0x14, 0x00, 0x14, 0x01, 0x40, 0x09, 0x01
			])
		);

		window.onload = function() {
			var div = document.getElementById("wasm-result");
			div.innerHTML = "<code>sum(1, 2)= " + mod.exports.sum(1, 2) + "</code>";
			console.log("mod-sum: "+ mod.exports.sum);
		};

		</script>

		<style>
		</style>
	</head>

	<body>
		<div id="header">
			<h2>WebAssembly</h2>
		</div>

		<div id="wasm-result"></div>

		<br/>
		<div>
			<h3>References</h3>
			<ul>
				<li><a href="https://medium.com/@MadsSejersen/webassembly-the-missing-tutorial-95f8580b08ba">WebAssembly — The missing tutorial</a></li>
			</ul>
		</div>
	</body>
</html>
`
