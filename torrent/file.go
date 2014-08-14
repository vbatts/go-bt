package torrent

import (
	//"github.com/vbatts/go-bt/bencode"
	"crypto/sha1"
)

/*
map[string]interface {}{"announce":"http://torrent.fedoraproject.org:6969/announce", "creation date":1387244350, "info":map[string]interface {}{"files":[]interface {}{map[string]interface {}{"length":1125, "path":[]interface {}{"Fedora-20-x86_64-CHECKSUM"}}, map[string]interface {}{"length":4603248640, "path":[]interface {}{"Fedora-20-x86_64-DVD.iso"}}}, "name":"Fedora-20-x86_64-DVD", "piece length":262144, "pieces":"m\x
*/
type File struct {
	// URL of a main tracker
	Announce string `bencode:"announce"`

	// Epoch of the creation of this torrent
	CreationDate int64 `bencode:"creation date"`

	// Dictionary about this torrent, including files to be tracked
	Info TorrentFileInfo `bencode:"info"`
}

type TorrentFileInfo struct {
	// suggested file/directory name where the file(s) are to be saved
	Name string `bencode:"name"`

	// hash list of joined SHA1 sums (160-bit length)
	Pieces string `bencode:"pieces"`

	// number of bytes per piece
	PieceLength int64 `bencode:"piece length"`

	// size of the file in bytes (only if this torrent is for a single file)
	Length int64 `bencode:"length"`

	// list of information about the files
	Files []FileInfo `bencode:"files"`
}

func (tfi TorrentFileInfo) PiecesList() []string {
	pieces := []string{}
	for i := 0; i < (len(tfi.Pieces) / sha1.Size); i++ {
		pieces = append(pieces, tfi.Pieces[i*sha1.Size:(i+1)*sha1.Size])
	}
	return pieces
}

type FileInfo struct {
	// size of file in bytes
	Length int64 `bencode:"length"`

	// list of strings corresponding to subdirectory names, the last of which is the actual file name
	Path []string `bencode:"path"`
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
		Info: TorrentFileInfo{
			Name:        infoName,
			Length:      infoLength,
			Pieces:      pieces,
			PieceLength: pieceLength,
			Files:       files,
		},
	}, nil
}
