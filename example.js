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