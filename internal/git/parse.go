package git

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func ParseTreeEntries(data []byte) ([]*TreeEntry, error) {
	return parseTreeEntries(data, nil)
}

func parseTreeEntries(data []byte, parentTree *Tree) ([]*TreeEntry, error) {
	entries := make([]*TreeEntry, 0, 10)
	for position := 0; position < len(data); {
		entry := new(TreeEntry)
		entry.parentTree = parentTree
		if position+6 > len(data) {
			return nil, fmt.Errorf("invalid ls-tree output: %s", string(data))
		}
		switch string(data[position : position+6]) {
		case "100644":
			entry.entryMode = EntryModeBlob
			position += 12 // skip over "100644 blob "
		case "100755":
			entry.entryMode = EntryModeExec
			position += 12 // skip over "100755 blob "
		case "120000":
			entry.entryMode = EntryModeSymlink
			position += 12 // skip over "120000 blob "
		case "160000":
			entry.entryMode = EntryModeCommit
			position += 14 // skip over "160000 object "
		case "040000":
			entry.entryMode = EntryModeTree
			position += 12 // skip over "040000 tree "
		default:
			return nil, fmt.Errorf("unknown type: %v", string(data[position:position+6]))
		}
		if position+40 > len(data) {
			return nil, fmt.Errorf("Invalid ls-tree output: %s", string(data))
		}
		id := string(data[position : position+40])
		entry.ID = id
		position += 41
		end := position + bytes.IndexByte(data[position:], '\t')
		if end < position {
			return nil, fmt.Errorf("Invalid ls-tree -l output: %s", string(data))
		}
		entry.size, _ = strconv.ParseInt(strings.TrimSpace(string(data[position:end])), 10, 64)
		entry.sized = true
		position = end + 1
		end = position + bytes.IndexByte(data[position:], '\n')
		if end < position {
			return nil, fmt.Errorf("Invalid ls-tree output: %s", string(data))
		}
		var err error
		if data[position] == '"' {
			entry.name, err = strconv.Unquote(string(data[position:end]))
			if err != nil {
				return nil, fmt.Errorf("Invalid ls-tree output: %v", err)
			}
		} else {
			entry.name = string(data[position:end])
		}
		position = end + 1
		entries = append(entries, entry)
	}
	return entries, nil
}
