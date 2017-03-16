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

		var mod = WebAssembly.instantiate(
			new Uint8Array([
				0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
				0x01, 0x07, 0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01,
				0x7f, 0x03, 0x02, 0x01, 0x00, 0x07, 0x07, 0x01,
				0x03, 0x61, 0x64, 0x64, 0x00, 0x00, 0x0a, 0x09,
				0x01, 0x07, 0x00, 0x20, 0x00, 0x20, 0x01, 0x6a,
				0x0b
			])
		).then(results =>
			results.instance
		);

		window.onload = function() {
			mod.then(function(instance) {
				var div = document.getElementById("wasm-result");
				div.innerHTML = "<code>sum(1, 2)= " + instance.exports.add(1, 2) + "</code>";
				console.log("mod-sum: "+ instance.exports.add);
			});
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
				<li><a href="https://developer.mozilla.org/en-US/docs/WebAssembly">MDN WebAssembly</a></li>
			</ul>
		</div>
	</body>
</html>
`
