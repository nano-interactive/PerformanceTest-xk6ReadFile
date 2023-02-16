# xk6-read-file

This is a [k6](https://go.k6.io/k6) extension using the
[xk6](https://github.com/grafana/xk6) system.

It is lightweight, fast and concurrent file reader. Each line will be read only once until 
the end of file, then it will start from the beginning.
The only way to preserve the read line order is to use one VU `-u 1`. 

It can be very helpful for reading very large files without storing it in memory 
with [SharedArray](https://k6.io/docs/javascript-api/k6-data/sharedarray/)

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  xk6 build --with github.com/nano-interactive/PerformanceTest-xk6ReadFile
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. 
It will help you create a `k6` binary with your local code rather than from GitHub.

```bash
make
```
Once built, you can run your newly extended `k6` using:
```shell
 ./k6 run -u 1 -i 200 example.js
 ```

## Available Functions
### readFile.openFile(`FILE_NAME`)
It should be initialised in `setup` javascript function. It points to a physical file from local folder.

### readFile.closeFile()
It must be called in `tearDown` javascript function to release pointer on opened file.

### readFile.readLine()
It reads next line from file. If the line is the end of file, it will rewind it to the beginning.

### readFile.setRewindFileUrl(`URL`)
It sends a GET request to requested URL everytime when file rewind happen. 
## Example

Make sure to open and close file in `setup()` and `teardown()`

```javascript

import readFile from 'k6/x/read-file';

let counter = 0;

export function setup() {
    readFile.setFileStartJsFunc(fileStarted)
    readFile.openFile('data_to_read.txt')
}
export default function () {
    let line = readFile.readLine();
    console.log(line);
}

export function teardown() {
    readFile.close()
}

// every time first row is read, it fires an event
// it needs to be initialised with readFile.setFileStartJsFunc(fileStarted)
export function fileStarted() {
    counter++;
    console.log("File started",counter,"time");
}
```

## Thanks
Credits: [SharedArray](https://k6.io/docs/javascript-api/k6-data/sharedarray/) and https://github.com/grafana/xk6-exec

## TODO
Make tests

## Licence
Apache License Version 2.0
