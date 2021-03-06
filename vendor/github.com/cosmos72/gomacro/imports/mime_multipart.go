// this file was generated by gomacro command: import _b "mime/multipart"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"mime/multipart"
)

// reflection: allow interpreted code to import "mime/multipart"
func init() {
	Packages["mime/multipart"] = Package{
	Binds: map[string]Value{
		"NewReader":	ValueOf(multipart.NewReader),
		"NewWriter":	ValueOf(multipart.NewWriter),
	},Types: map[string]Type{
		"File":	TypeOf((*multipart.File)(nil)).Elem(),
		"FileHeader":	TypeOf((*multipart.FileHeader)(nil)).Elem(),
		"Form":	TypeOf((*multipart.Form)(nil)).Elem(),
		"Part":	TypeOf((*multipart.Part)(nil)).Elem(),
		"Reader":	TypeOf((*multipart.Reader)(nil)).Elem(),
		"Writer":	TypeOf((*multipart.Writer)(nil)).Elem(),
	},Proxies: map[string]Type{
		"File":	TypeOf((*File_mime_multipart)(nil)).Elem(),
	},
	}
}

// --------------- proxy for mime/multipart.File ---------------
type File_mime_multipart struct {
	Object	interface{}
	Close_	func(interface{}) error
	Read_	func(_proxy_obj_ interface{}, p []byte) (n int, err error)
	ReadAt_	func(_proxy_obj_ interface{}, p []byte, off int64) (n int, err error)
	Seek_	func(_proxy_obj_ interface{}, offset int64, whence int) (int64, error)
}
func (Proxy *File_mime_multipart) Close() error {
	return Proxy.Close_(Proxy.Object)
}
func (Proxy *File_mime_multipart) Read(p []byte) (n int, err error) {
	return Proxy.Read_(Proxy.Object, p)
}
func (Proxy *File_mime_multipart) ReadAt(p []byte, off int64) (n int, err error) {
	return Proxy.ReadAt_(Proxy.Object, p, off)
}
func (Proxy *File_mime_multipart) Seek(offset int64, whence int) (int64, error) {
	return Proxy.Seek_(Proxy.Object, offset, whence)
}
