## Description
Counts the total characters, words and lines of a single file or multiple files. Also gives a total words, characters and lines count of all the
files if multiple files are passed as arguments. Use flags `[-l|-w|-c]` to only show line, word and character count of the files respectively.

### How to run?
`go run . [-l|-w|-c] file1 file2...`

## Example:
**Without flag:**
`go run . *.txt`

**With flag:**
`go run . -l *.txt`
