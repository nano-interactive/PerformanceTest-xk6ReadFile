import readFile from 'k6/x/read-file';

let counter = 0;

export function setup() {
  readFile.setRewindFileUrl("http://localhost:8080/clear-cache");
  readFile.openFile('data_to_read.txt')
}
export default function () {
  let line = readFile.readLine();
}

export function teardown() {
  readFile.close()
}
