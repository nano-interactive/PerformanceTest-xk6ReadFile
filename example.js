import readFile from 'k6/x/read-file';


export function setup() {
  readFile.openFile('data_to_read.txt')

}
export default function () {
  console.log(readFile.readLine());
}

export function teardown() {
  readFile.close()
}