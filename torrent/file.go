package torrent

import (
	"crypto/sha1"
	"time"
)

// File is a representation of https://wiki.theory.org/BitTorrentSpecification#Metainfo_File_Structure
type File struct {
	// Announce is the URL of a main tracker
	Announce string `bencode:"announce"`

	// AnnounceList lists additional trackers
	AnnounceList []string `bencode:"announce-list"`

	// CreationDate is epoch of the creation of this torrent
	CreationDate int64 `bencode:"creation date,omitempty"`

	// Info is the dictionary about this torrent, including files to be tracked
	Info InfoSection `bencode:"info,omitempty"`

	// Comment is free-form textual comments of the author
	Comment string `bencode:"comment,omitempty"`

	// CreatedBy is the name and version of the program used to create the .torrent
	CreatedBy string `bencode:"created by,omitempty"`

	// Encoding string encoding used to generate the `pieces` and `info` fields
	Encoding string `bencode:"encoding"`
}

// CreationDateTime returns a time from the File's CreationDate
func (f File) CreationDateTime() time.Time {
	return time.Unix(f.CreationDate, 0)
}

type InfoSection struct {
	// suggested file/directory name where the file(s) are to be saved
	Name string `bencode:"name,omitempty"`

	// hash list of joined SHA1 sums (160-bit length)
	Pieces string `bencode:"pieces"`

	// number of bytes per piece
	PieceLength int64 `bencode:"piece length"`

	// size of the file in bytes (only if this torrent is for a single file)
	Length int64 `bencode:"length,omitempty"`

	// 32-char hexadecimal string corresponding to the MD5 sum of the file (only if this torrent is for a single file)
	MD5 string `bencode:"md5sum,omitempty"`

	// list of information about the files
	Files []FileInfo `bencode:"files"`
}

func (is InfoSection) PiecesList() []string {
	pieces := []string{}
	for i := 0; i < (len(is.Pieces) / sha1.Size); i++ {
		pieces = append(pieces, is.Pieces[i*sha1.Size:(i+1)*sha1.Size])
	}
	return pieces
}

type FileInfo struct {
	// size of file in bytes
	Length int64 `bencode:"length"`

	// list of strings corresponding to subdirectory names, the last of which is the actual file name
	Path []string `bencode:"path"`

	// 32-char hexadecimal string corresponding to the MD5 sum of the file (only if this torrent is for a single file)
	MD5 string `bencode:"md5sum,omitempty"`
}

type torrentError struct {
	Msg string
}

func (te torrentError) Error() string {
	return te.Msg
}

var (
	ErrNotProperDataInterface = torrentError{"data does not look like map[string]interface{}"}
)

// DecocdeTorrentData translates a bencode decoded message from a *.torrent file.
// Assuming a generic map[string]interface{} data.
func DecocdeTorrentData(data interface{}) (*File, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, ErrNotProperDataInterface
	}
	announce := m["announce"].(string)
	creationDate := m["creation date"].(int64)

	info := m["info"].(map[string]interface{})
	pieceLength := info["piece length"].(int64)
	pieces := info["pieces"].(string)
	infoName := info["name"].(string)

	isSingleFileTorrent := true
	infoFiles, ok := info["files"].([]interface{})
	if ok {
		isSingleFileTorrent = false
	}
	infoLength := int64(0)
	if isSingleFileTorrent {
		infoLength = info["length"].(int64)
	}
	files := []FileInfo{}
	if !isSingleFileTorrent {
		for _, fileInterface := range infoFiles {
			fileInfo := fileInterface.(map[string]interface{})
			paths := []string{}
			for _, path := range fileInfo["path"].([]interface{}) {
				paths = append(paths, path.(string))
			}
			files = append(files, FileInfo{
				Length: fileInfo["length"].(int64),
				Path:   paths,
			})
		}
	}

	return &File{
		Announce:     announce,
		CreationDate: creationDate,
		Info: InfoSection{
			Name:        infoName,
			Length:      infoLength,
			Pieces:      pieces,
			PieceLength: pieceLength,
			Files:       files,
		},
	}, nil
}
