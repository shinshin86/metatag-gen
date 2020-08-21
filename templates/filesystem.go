// Code generated by shogo82148/assets-life v1.0.0. DO NOT EDIT.

//go:generate go run assets-life.go "." . templates

package templates

import (
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// Root is the root of the file system.
var Root http.FileSystem = fileSystem{
	file{
		name:    "/",
		content: "",
		mode:    0755 | os.ModeDir,
		next:    0,
		child:   1,
	},
	file{
		name:    "/assets-life.go",
		content: "// Copyright (C) 2019 Ichinose Shogo All rights reserved.\n// Use of this source code is governed by a MIT-style\n// license that can be found in https://github.com/shogo82148/assets-life/blob/master/LICENSE\n\n// +build ignore\n\n// assets-life is a very simple embedding asset generator.\n// It generates an embed small in-memory file system that is served from an http.FileSystem.\n// Install the command line tool first.\n//\n//     go get github.com/shogo82148/assets-life\n//\n// The assets-life command generates a package that have embed small in-memory file system.\n//\n//     assets-life /path/to/your/project/public public\n//\n// You can access the file system by accessing a public variable Root of the generated package.\n//\n//     import (\n//         \"net/http\"\n//         \"./public\" // TODO: Replace with the absolute import path\n//     )\n//\n//     func main() {\n//         http.Handle(\"/\", http.FileServer(public.Root))\n//         http.ListenAndServe(\":8080\", nil)\n//     }\n//\n// Visit http://localhost:8080/path/to/file to see your file.\n//\n// The assets-life command also embed go:generate directive into generated code, and assets-life itself.\n// It allows you to re-generate the package using go generate.\n//\n//     go generate ./public\n//\n// The assets-life command is no longer needed because it is embedded into the generated package.\npackage main\n\nimport (\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"log\"\n\t\"os\"\n\t\"path\"\n\t\"path/filepath\"\n\t\"strings\"\n)\n\nconst version = \"1.0.0\"\n\nfunc main() {\n\tif len(os.Args) <= 2 {\n\t\tlog.Println(\"Usage:\")\n\t\tlog.Println(os.Args[0] + \" INPUT_DIR OUTPUT_DIR [PACKAGE_NAME]\")\n\t\tos.Exit(2)\n\t}\n\tin, err := filepath.Abs(os.Args[1])\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tout, err := filepath.Abs(os.Args[2])\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tvar name string\n\tif len(os.Args) > 3 {\n\t\tname = os.Args[3]\n\t}\n\tif name == \"\" {\n\t\tname = filepath.Base(out)\n\t}\n\tif err := build(in, out, name); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n\nfunc build(in, out, name string) error {\n\tfilename := \"assets-life.go\"\n\trel, err := filepath.Rel(out, in)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif err := os.MkdirAll(out, 0755); err != nil {\n\t\treturn err\n\t}\n\tf, err := os.OpenFile(filepath.Join(out, \"filesystem.go\"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\n\tif err != nil {\n\t\treturn err\n\t}\n\theader := `// Code generated by shogo82148/assets-life v%s. DO NOT EDIT.\n\n//%s\n\npackage %s\n\nimport (\n\t\"io\"\n\t\"net/http\"\n\t\"os\"\n\t\"path\"\n\t\"sort\"\n\t\"strings\"\n\t\"time\"\n)\n\n// Root is the root of the file system.\nvar Root http.FileSystem = fileSystem{\n`\n\trel = filepath.ToSlash(rel)\n\tfmt.Fprintf(f, header, version, \"go:generate go run \"+filename+\" \\\"\"+rel+\"\\\" . \"+name, name)\n\n\ttype file struct {\n\t\tpath     string\n\t\tmode     os.FileMode\n\t\tchildren []int\n\t\tnext     int\n\t}\n\tindex := map[string]int{}\n\tfiles := []file{}\n\n\tvar i int\n\terr = filepath.Walk(in, func(path string, info os.FileInfo, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\t// ignore hidden files\n\t\tif strings.HasPrefix(info.Name(), \".\") {\n\t\t\treturn nil\n\t\t}\n\n\t\tif (info.Mode()&os.ModeType)|os.ModeDir != os.ModeDir {\n\t\t\treturn fmt.Errorf(\"unsupported file type: %s, mode %s\", path, info.Mode())\n\t\t}\n\n\t\tindex[path] = i\n\t\tfiles = append(files, file{\n\t\t\tpath: path,\n\t\t\tmode: info.Mode(),\n\t\t})\n\t\tparent := filepath.Dir(path)\n\t\tif idx, ok := index[parent]; ok {\n\t\t\tfiles[idx].children = append(files[idx].children, i)\n\t\t}\n\t\ti++\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, ff := range files {\n\t\t// search neighborhood\n\t\tfor i := range ff.children {\n\t\t\tnext := -1\n\t\t\tif i+1 < len(ff.children) {\n\t\t\t\tnext = ff.children[i+1]\n\t\t\t}\n\t\t\tfiles[ff.children[i]].next = next\n\t\t}\n\n\t\tfmt.Fprintf(f, \"\\tfile{\\n\")\n\t\trel, err := filepath.Rel(in, ff.path)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfmt.Fprintf(f, \"\\t\\tname:    %q,\\n\", path.Clean(\"/\"+filepath.ToSlash(rel)))\n\t\tif ff.mode.IsDir() {\n\t\t\tfmt.Fprintln(f, \"\\t\\tcontent: \\\"\\\",\")\n\t\t} else {\n\t\t\tb, err := ioutil.ReadFile(ff.path)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\tfmt.Fprintf(f, \"\\t\\tcontent: %q,\\n\", string(b))\n\t\t}\n\t\tswitch {\n\t\tcase ff.mode.IsDir(): // directory\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0755 | os.ModeDir,\")\n\t\tcase ff.mode&0100 != 0: // executable file\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0755,\")\n\t\tdefault:\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0644,\")\n\t\t}\n\t\tfmt.Fprintf(f, \"\\t\\tnext:    %d,\\n\", ff.next)\n\t\tif len(ff.children) > 0 {\n\t\t\tfmt.Fprintf(f, \"\\t\\tchild:   %d,\\n\", ff.children[0])\n\t\t} else {\n\t\t\tfmt.Fprint(f, \"\\t\\tchild:   -1,\\n\")\n\t\t}\n\t\tfmt.Fprint(f, \"\\t},\\n\")\n\t}\n\tfooter := `}\n\ntype fileSystem []file\n\nfunc (fs fileSystem) Open(name string) (http.File, error) {\n\tname = path.Clean(\"/\" + name)\n\ti := sort.Search(len(fs), func(i int) bool { return fs[i].name >= name })\n\tif i >= len(fs) || fs[i].name != name {\n\t\treturn nil, &os.PathError{\n\t\t\tOp:   \"open\",\n\t\t\tPath: name,\n\t\t\tErr:  os.ErrNotExist,\n\t\t}\n\t}\n\tf := &fs[i]\n\treturn &httpFile{\n\t\tReader: strings.NewReader(f.content),\n\t\tfile:   f,\n\t\tfs:     fs,\n\t\tidx:    i,\n\t\tdirIdx: f.child,\n\t}, nil\n}\n\ntype file struct {\n\tname    string\n\tcontent string\n\tmode    os.FileMode\n\tchild   int\n\tnext    int\n}\n\nvar _ os.FileInfo = (*file)(nil)\n\nfunc (f *file) Name() string {\n\treturn path.Base(f.name)\n}\n\nfunc (f *file) Size() int64 {\n\treturn int64(len(f.content))\n}\n\nfunc (f *file) Mode() os.FileMode {\n\treturn f.mode\n}\n\nvar zeroTime time.Time\n\nfunc (f *file) ModTime() time.Time {\n\treturn zeroTime\n}\n\nfunc (f *file) IsDir() bool {\n\treturn f.Mode().IsDir()\n}\n\nfunc (f *file) Sys() interface{} {\n\treturn nil\n}\n\ntype httpFile struct {\n\t*strings.Reader\n\tfile   *file\n\tfs     fileSystem\n\tidx    int\n\tdirIdx int\n}\n\nvar _ http.File = (*httpFile)(nil)\n\nfunc (f *httpFile) Stat() (os.FileInfo, error) {\n\treturn f.file, nil\n}\n\nfunc (f *httpFile) Readdir(count int) ([]os.FileInfo, error) {\n\tret := []os.FileInfo{}\n\tif !f.file.IsDir() {\n\t\treturn ret, nil\n\t}\n\n\tif count <= 0 {\n\t\tfor f.dirIdx >= 0 {\n\t\t\tentry := &f.fs[f.dirIdx]\n\t\t\tret = append(ret, entry)\n\t\t\tf.dirIdx = entry.next\n\t\t}\n\t\treturn ret, nil\n\t}\n\n\tret = make([]os.FileInfo, 0, count)\n\tfor f.dirIdx >= 0 {\n\t\tentry := &f.fs[f.dirIdx]\n\t\tret = append(ret, entry)\n\t\tf.dirIdx = entry.next\n\t\tif len(ret) == count {\n\t\t\treturn ret, nil\n\t\t}\n\t}\n\treturn ret, io.EOF\n}\n\nfunc (f *httpFile) Close() error {\n\treturn nil\n}`\n\tfmt.Fprintln(f, footer)\n\tif err := f.Close(); err != nil {\n\t\treturn err\n\t}\n\n\tf, err = os.OpenFile(filepath.Join(out, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\n\tformat := `// Copyright (C) 2019 Ichinose Shogo All rights reserved.\n// Use of this source code is governed by a MIT-style\n// license that can be found in https://github.com/shogo82148/assets-life/blob/master/LICENSE\n\n// +build ignore\n\n// assets-life is a very simple embedding asset generator.\n// It generates an embed small in-memory file system that is served from an http.FileSystem.\n// Install the command line tool first.\n//\n//     go get github.com/shogo82148/assets-life\n//\n// The assets-life command generates a package that have embed small in-memory file system.\n//\n//     assets-life /path/to/your/project/public public\n//\n// You can access the file system by accessing a public variable Root of the generated package.\n//\n//     import (\n//         \"net/http\"\n//         \"./public\" // TODO: Replace with the absolute import path\n//     )\n//\n//     func main() {\n//         http.Handle(\"/\", http.FileServer(public.Root))\n//         http.ListenAndServe(\":8080\", nil)\n//     }\n//\n// Visit http://localhost:8080/path/to/file to see your file.\n//\n// The assets-life command also embed go:generate directive into generated code, and assets-life itself.\n// It allows you to re-generate the package using go generate.\n//\n//     go generate ./public\n//\n// The assets-life command is no longer needed because it is embedded into the generated package.\npackage main\n\nimport (\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"log\"\n\t\"os\"\n\t\"path\"\n\t\"path/filepath\"\n\t\"strings\"\n)\n\nconst version = \"%s\"\n\nfunc main() {\n\tif len(os.Args) <= 2 {\n\t\tlog.Println(\"Usage:\")\n\t\tlog.Println(os.Args[0] + \" INPUT_DIR OUTPUT_DIR [PACKAGE_NAME]\")\n\t\tos.Exit(2)\n\t}\n\tin, err := filepath.Abs(os.Args[1])\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tout, err := filepath.Abs(os.Args[2])\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tvar name string\n\tif len(os.Args) > 3 {\n\t\tname = os.Args[3]\n\t}\n\tif name == \"\" {\n\t\tname = filepath.Base(out)\n\t}\n\tif err := build(in, out, name); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n\nfunc build(in, out, name string) error {\n\tfilename := \"assets-life.go\"\n\trel, err := filepath.Rel(out, in)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif err := os.MkdirAll(out, 0755); err != nil {\n\t\treturn err\n\t}\n\tf, err := os.OpenFile(filepath.Join(out, \"filesystem.go\"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\n\tif err != nil {\n\t\treturn err\n\t}\n\theader := %c%s%c\n\trel = filepath.ToSlash(rel)\n\tfmt.Fprintf(f, header, version, \"go:generate go run \"+filename+\" \\\"\"+rel+\"\\\" . \"+name, name)\n\n\ttype file struct {\n\t\tpath     string\n\t\tmode     os.FileMode\n\t\tchildren []int\n\t\tnext     int\n\t}\n\tindex := map[string]int{}\n\tfiles := []file{}\n\n\tvar i int\n\terr = filepath.Walk(in, func(path string, info os.FileInfo, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\t// ignore hidden files\n\t\tif strings.HasPrefix(info.Name(), \".\") {\n\t\t\treturn nil\n\t\t}\n\n\t\tif (info.Mode()&os.ModeType)|os.ModeDir != os.ModeDir {\n\t\t\treturn fmt.Errorf(\"unsupported file type: %%s, mode %%s\", path, info.Mode())\n\t\t}\n\n\t\tindex[path] = i\n\t\tfiles = append(files, file{\n\t\t\tpath: path,\n\t\t\tmode: info.Mode(),\n\t\t})\n\t\tparent := filepath.Dir(path)\n\t\tif idx, ok := index[parent]; ok {\n\t\t\tfiles[idx].children = append(files[idx].children, i)\n\t\t}\n\t\ti++\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, ff := range files {\n\t\t// search neighborhood\n\t\tfor i := range ff.children {\n\t\t\tnext := -1\n\t\t\tif i+1 < len(ff.children) {\n\t\t\t\tnext = ff.children[i+1]\n\t\t\t}\n\t\t\tfiles[ff.children[i]].next = next\n\t\t}\n\n\t\tfmt.Fprintf(f, \"\\tfile{\\n\")\n\t\trel, err := filepath.Rel(in, ff.path)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfmt.Fprintf(f, \"\\t\\tname:    %%q,\\n\", path.Clean(\"/\"+filepath.ToSlash(rel)))\n\t\tif ff.mode.IsDir() {\n\t\t\tfmt.Fprintln(f, \"\\t\\tcontent: \\\"\\\",\")\n\t\t} else {\n\t\t\tb, err := ioutil.ReadFile(ff.path)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\tfmt.Fprintf(f, \"\\t\\tcontent: %%q,\\n\", string(b))\n\t\t}\n\t\tswitch {\n\t\tcase ff.mode.IsDir(): // directory\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0755 | os.ModeDir,\")\n\t\tcase ff.mode&0100 != 0: // executable file\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0755,\")\n\t\tdefault:\n\t\t\tfmt.Fprintln(f, \"\\t\\tmode:    0644,\")\n\t\t}\n\t\tfmt.Fprintf(f, \"\\t\\tnext:    %%d,\\n\", ff.next)\n\t\tif len(ff.children) > 0 {\n\t\t\tfmt.Fprintf(f, \"\\t\\tchild:   %%d,\\n\", ff.children[0])\n\t\t} else {\n\t\t\tfmt.Fprint(f, \"\\t\\tchild:   -1,\\n\")\n\t\t}\n\t\tfmt.Fprint(f, \"\\t},\\n\")\n\t}\n\tfooter := %c%s%c\n\tfmt.Fprintln(f, footer)\n\tif err := f.Close(); err != nil {\n\t\treturn err\n\t}\n\n\tf, err = os.OpenFile(filepath.Join(out, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\n\tformat := %c%s%c\n\tfmt.Fprintf(f, format, version, 96, header, 96, 96, footer, 96, 96, format, 96)\n\tif err := f.Close(); err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n`\n\tfmt.Fprintf(f, format, version, 96, header, 96, 96, footer, 96, 96, format, 96)\n\tif err := f.Close(); err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n",
		mode:    0644,
		next:    2,
		child:   -1,
	},
	file{
		name:    "/filesystem.go",
		content: "// Code generated by shogo82148/assets-life v1.0.0. DO NOT EDIT.\n\n//go:generate go run assets-life.go \".\" . templates\n\npackage templates\n\nimport (\n\t\"io\"\n\t\"net/http\"\n\t\"os\"\n\t\"path\"\n\t\"sort\"\n\t\"strings\"\n\t\"time\"\n)\n\n// Root is the root of the file system.\nvar Root http.FileSystem = fileSystem{\n\tfile{\n\t\tname:    \"/\",\n\t\tcontent: \"\",\n\t\tmode:    0755 | os.ModeDir,\n\t\tnext:    0,\n\t\tchild:   1,\n\t},\n\tfile{\n\t\tname:    \"/assets-life.go\",\n\t\tcontent: \"// Copyright (C) 2019 Ichinose Shogo All rights reserved.\\n// Use of this source code is governed by a MIT-style\\n// license that can be found in https://github.com/shogo82148/assets-life/blob/master/LICENSE\\n\\n// +build ignore\\n\\n// assets-life is a very simple embedding asset generator.\\n// It generates an embed small in-memory file system that is served from an http.FileSystem.\\n// Install the command line tool first.\\n//\\n//     go get github.com/shogo82148/assets-life\\n//\\n// The assets-life command generates a package that have embed small in-memory file system.\\n//\\n//     assets-life /path/to/your/project/public public\\n//\\n// You can access the file system by accessing a public variable Root of the generated package.\\n//\\n//     import (\\n//         \\\"net/http\\\"\\n//         \\\"./public\\\" // TODO: Replace with the absolute import path\\n//     )\\n//\\n//     func main() {\\n//         http.Handle(\\\"/\\\", http.FileServer(public.Root))\\n//         http.ListenAndServe(\\\":8080\\\", nil)\\n//     }\\n//\\n// Visit http://localhost:8080/path/to/file to see your file.\\n//\\n// The assets-life command also embed go:generate directive into generated code, and assets-life itself.\\n// It allows you to re-generate the package using go generate.\\n//\\n//     go generate ./public\\n//\\n// The assets-life command is no longer needed because it is embedded into the generated package.\\npackage main\\n\\nimport (\\n\\t\\\"fmt\\\"\\n\\t\\\"io/ioutil\\\"\\n\\t\\\"log\\\"\\n\\t\\\"os\\\"\\n\\t\\\"path\\\"\\n\\t\\\"path/filepath\\\"\\n\\t\\\"strings\\\"\\n)\\n\\nconst version = \\\"1.0.0\\\"\\n\\nfunc main() {\\n\\tif len(os.Args) <= 2 {\\n\\t\\tlog.Println(\\\"Usage:\\\")\\n\\t\\tlog.Println(os.Args[0] + \\\" INPUT_DIR OUTPUT_DIR [PACKAGE_NAME]\\\")\\n\\t\\tos.Exit(2)\\n\\t}\\n\\tin, err := filepath.Abs(os.Args[1])\\n\\tif err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n\\tout, err := filepath.Abs(os.Args[2])\\n\\tif err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n\\tvar name string\\n\\tif len(os.Args) > 3 {\\n\\t\\tname = os.Args[3]\\n\\t}\\n\\tif name == \\\"\\\" {\\n\\t\\tname = filepath.Base(out)\\n\\t}\\n\\tif err := build(in, out, name); err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n}\\n\\nfunc build(in, out, name string) error {\\n\\tfilename := \\\"assets-life.go\\\"\\n\\trel, err := filepath.Rel(out, in)\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.MkdirAll(out, 0755); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tf, err := os.OpenFile(filepath.Join(out, \\\"filesystem.go\\\"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\theader := `// Code generated by shogo82148/assets-life v%s. DO NOT EDIT.\\n\\n//%s\\n\\npackage %s\\n\\nimport (\\n\\t\\\"io\\\"\\n\\t\\\"net/http\\\"\\n\\t\\\"os\\\"\\n\\t\\\"path\\\"\\n\\t\\\"sort\\\"\\n\\t\\\"strings\\\"\\n\\t\\\"time\\\"\\n)\\n\\n// Root is the root of the file system.\\nvar Root http.FileSystem = fileSystem{\\n`\\n\\trel = filepath.ToSlash(rel)\\n\\tfmt.Fprintf(f, header, version, \\\"go:generate go run \\\"+filename+\\\" \\\\\\\"\\\"+rel+\\\"\\\\\\\" . \\\"+name, name)\\n\\n\\ttype file struct {\\n\\t\\tpath     string\\n\\t\\tmode     os.FileMode\\n\\t\\tchildren []int\\n\\t\\tnext     int\\n\\t}\\n\\tindex := map[string]int{}\\n\\tfiles := []file{}\\n\\n\\tvar i int\\n\\terr = filepath.Walk(in, func(path string, info os.FileInfo, err error) error {\\n\\t\\tif err != nil {\\n\\t\\t\\treturn err\\n\\t\\t}\\n\\n\\t\\t// ignore hidden files\\n\\t\\tif strings.HasPrefix(info.Name(), \\\".\\\") {\\n\\t\\t\\treturn nil\\n\\t\\t}\\n\\n\\t\\tif (info.Mode()&os.ModeType)|os.ModeDir != os.ModeDir {\\n\\t\\t\\treturn fmt.Errorf(\\\"unsupported file type: %s, mode %s\\\", path, info.Mode())\\n\\t\\t}\\n\\n\\t\\tindex[path] = i\\n\\t\\tfiles = append(files, file{\\n\\t\\t\\tpath: path,\\n\\t\\t\\tmode: info.Mode(),\\n\\t\\t})\\n\\t\\tparent := filepath.Dir(path)\\n\\t\\tif idx, ok := index[parent]; ok {\\n\\t\\t\\tfiles[idx].children = append(files[idx].children, i)\\n\\t\\t}\\n\\t\\ti++\\n\\t\\treturn nil\\n\\t})\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\n\\tfor _, ff := range files {\\n\\t\\t// search neighborhood\\n\\t\\tfor i := range ff.children {\\n\\t\\t\\tnext := -1\\n\\t\\t\\tif i+1 < len(ff.children) {\\n\\t\\t\\t\\tnext = ff.children[i+1]\\n\\t\\t\\t}\\n\\t\\t\\tfiles[ff.children[i]].next = next\\n\\t\\t}\\n\\n\\t\\tfmt.Fprintf(f, \\\"\\\\tfile{\\\\n\\\")\\n\\t\\trel, err := filepath.Rel(in, ff.path)\\n\\t\\tif err != nil {\\n\\t\\t\\treturn err\\n\\t\\t}\\n\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tname:    %q,\\\\n\\\", path.Clean(\\\"/\\\"+filepath.ToSlash(rel)))\\n\\t\\tif ff.mode.IsDir() {\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tcontent: \\\\\\\"\\\\\\\",\\\")\\n\\t\\t} else {\\n\\t\\t\\tb, err := ioutil.ReadFile(ff.path)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn err\\n\\t\\t\\t}\\n\\t\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tcontent: %q,\\\\n\\\", string(b))\\n\\t\\t}\\n\\t\\tswitch {\\n\\t\\tcase ff.mode.IsDir(): // directory\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0755 | os.ModeDir,\\\")\\n\\t\\tcase ff.mode&0100 != 0: // executable file\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0755,\\\")\\n\\t\\tdefault:\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0644,\\\")\\n\\t\\t}\\n\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tnext:    %d,\\\\n\\\", ff.next)\\n\\t\\tif len(ff.children) > 0 {\\n\\t\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tchild:   %d,\\\\n\\\", ff.children[0])\\n\\t\\t} else {\\n\\t\\t\\tfmt.Fprint(f, \\\"\\\\t\\\\tchild:   -1,\\\\n\\\")\\n\\t\\t}\\n\\t\\tfmt.Fprint(f, \\\"\\\\t},\\\\n\\\")\\n\\t}\\n\\tfooter := `}\\n\\ntype fileSystem []file\\n\\nfunc (fs fileSystem) Open(name string) (http.File, error) {\\n\\tname = path.Clean(\\\"/\\\" + name)\\n\\ti := sort.Search(len(fs), func(i int) bool { return fs[i].name >= name })\\n\\tif i >= len(fs) || fs[i].name != name {\\n\\t\\treturn nil, &os.PathError{\\n\\t\\t\\tOp:   \\\"open\\\",\\n\\t\\t\\tPath: name,\\n\\t\\t\\tErr:  os.ErrNotExist,\\n\\t\\t}\\n\\t}\\n\\tf := &fs[i]\\n\\treturn &httpFile{\\n\\t\\tReader: strings.NewReader(f.content),\\n\\t\\tfile:   f,\\n\\t\\tfs:     fs,\\n\\t\\tidx:    i,\\n\\t\\tdirIdx: f.child,\\n\\t}, nil\\n}\\n\\ntype file struct {\\n\\tname    string\\n\\tcontent string\\n\\tmode    os.FileMode\\n\\tchild   int\\n\\tnext    int\\n}\\n\\nvar _ os.FileInfo = (*file)(nil)\\n\\nfunc (f *file) Name() string {\\n\\treturn path.Base(f.name)\\n}\\n\\nfunc (f *file) Size() int64 {\\n\\treturn int64(len(f.content))\\n}\\n\\nfunc (f *file) Mode() os.FileMode {\\n\\treturn f.mode\\n}\\n\\nvar zeroTime time.Time\\n\\nfunc (f *file) ModTime() time.Time {\\n\\treturn zeroTime\\n}\\n\\nfunc (f *file) IsDir() bool {\\n\\treturn f.Mode().IsDir()\\n}\\n\\nfunc (f *file) Sys() interface{} {\\n\\treturn nil\\n}\\n\\ntype httpFile struct {\\n\\t*strings.Reader\\n\\tfile   *file\\n\\tfs     fileSystem\\n\\tidx    int\\n\\tdirIdx int\\n}\\n\\nvar _ http.File = (*httpFile)(nil)\\n\\nfunc (f *httpFile) Stat() (os.FileInfo, error) {\\n\\treturn f.file, nil\\n}\\n\\nfunc (f *httpFile) Readdir(count int) ([]os.FileInfo, error) {\\n\\tret := []os.FileInfo{}\\n\\tif !f.file.IsDir() {\\n\\t\\treturn ret, nil\\n\\t}\\n\\n\\tif count <= 0 {\\n\\t\\tfor f.dirIdx >= 0 {\\n\\t\\t\\tentry := &f.fs[f.dirIdx]\\n\\t\\t\\tret = append(ret, entry)\\n\\t\\t\\tf.dirIdx = entry.next\\n\\t\\t}\\n\\t\\treturn ret, nil\\n\\t}\\n\\n\\tret = make([]os.FileInfo, 0, count)\\n\\tfor f.dirIdx >= 0 {\\n\\t\\tentry := &f.fs[f.dirIdx]\\n\\t\\tret = append(ret, entry)\\n\\t\\tf.dirIdx = entry.next\\n\\t\\tif len(ret) == count {\\n\\t\\t\\treturn ret, nil\\n\\t\\t}\\n\\t}\\n\\treturn ret, io.EOF\\n}\\n\\nfunc (f *httpFile) Close() error {\\n\\treturn nil\\n}`\\n\\tfmt.Fprintln(f, footer)\\n\\tif err := f.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\n\\tf, err = os.OpenFile(filepath.Join(out, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\\n\\tformat := `// Copyright (C) 2019 Ichinose Shogo All rights reserved.\\n// Use of this source code is governed by a MIT-style\\n// license that can be found in https://github.com/shogo82148/assets-life/blob/master/LICENSE\\n\\n// +build ignore\\n\\n// assets-life is a very simple embedding asset generator.\\n// It generates an embed small in-memory file system that is served from an http.FileSystem.\\n// Install the command line tool first.\\n//\\n//     go get github.com/shogo82148/assets-life\\n//\\n// The assets-life command generates a package that have embed small in-memory file system.\\n//\\n//     assets-life /path/to/your/project/public public\\n//\\n// You can access the file system by accessing a public variable Root of the generated package.\\n//\\n//     import (\\n//         \\\"net/http\\\"\\n//         \\\"./public\\\" // TODO: Replace with the absolute import path\\n//     )\\n//\\n//     func main() {\\n//         http.Handle(\\\"/\\\", http.FileServer(public.Root))\\n//         http.ListenAndServe(\\\":8080\\\", nil)\\n//     }\\n//\\n// Visit http://localhost:8080/path/to/file to see your file.\\n//\\n// The assets-life command also embed go:generate directive into generated code, and assets-life itself.\\n// It allows you to re-generate the package using go generate.\\n//\\n//     go generate ./public\\n//\\n// The assets-life command is no longer needed because it is embedded into the generated package.\\npackage main\\n\\nimport (\\n\\t\\\"fmt\\\"\\n\\t\\\"io/ioutil\\\"\\n\\t\\\"log\\\"\\n\\t\\\"os\\\"\\n\\t\\\"path\\\"\\n\\t\\\"path/filepath\\\"\\n\\t\\\"strings\\\"\\n)\\n\\nconst version = \\\"%s\\\"\\n\\nfunc main() {\\n\\tif len(os.Args) <= 2 {\\n\\t\\tlog.Println(\\\"Usage:\\\")\\n\\t\\tlog.Println(os.Args[0] + \\\" INPUT_DIR OUTPUT_DIR [PACKAGE_NAME]\\\")\\n\\t\\tos.Exit(2)\\n\\t}\\n\\tin, err := filepath.Abs(os.Args[1])\\n\\tif err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n\\tout, err := filepath.Abs(os.Args[2])\\n\\tif err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n\\tvar name string\\n\\tif len(os.Args) > 3 {\\n\\t\\tname = os.Args[3]\\n\\t}\\n\\tif name == \\\"\\\" {\\n\\t\\tname = filepath.Base(out)\\n\\t}\\n\\tif err := build(in, out, name); err != nil {\\n\\t\\tlog.Fatal(err)\\n\\t}\\n}\\n\\nfunc build(in, out, name string) error {\\n\\tfilename := \\\"assets-life.go\\\"\\n\\trel, err := filepath.Rel(out, in)\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.MkdirAll(out, 0755); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tf, err := os.OpenFile(filepath.Join(out, \\\"filesystem.go\\\"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\theader := %c%s%c\\n\\trel = filepath.ToSlash(rel)\\n\\tfmt.Fprintf(f, header, version, \\\"go:generate go run \\\"+filename+\\\" \\\\\\\"\\\"+rel+\\\"\\\\\\\" . \\\"+name, name)\\n\\n\\ttype file struct {\\n\\t\\tpath     string\\n\\t\\tmode     os.FileMode\\n\\t\\tchildren []int\\n\\t\\tnext     int\\n\\t}\\n\\tindex := map[string]int{}\\n\\tfiles := []file{}\\n\\n\\tvar i int\\n\\terr = filepath.Walk(in, func(path string, info os.FileInfo, err error) error {\\n\\t\\tif err != nil {\\n\\t\\t\\treturn err\\n\\t\\t}\\n\\n\\t\\t// ignore hidden files\\n\\t\\tif strings.HasPrefix(info.Name(), \\\".\\\") {\\n\\t\\t\\treturn nil\\n\\t\\t}\\n\\n\\t\\tif (info.Mode()&os.ModeType)|os.ModeDir != os.ModeDir {\\n\\t\\t\\treturn fmt.Errorf(\\\"unsupported file type: %%s, mode %%s\\\", path, info.Mode())\\n\\t\\t}\\n\\n\\t\\tindex[path] = i\\n\\t\\tfiles = append(files, file{\\n\\t\\t\\tpath: path,\\n\\t\\t\\tmode: info.Mode(),\\n\\t\\t})\\n\\t\\tparent := filepath.Dir(path)\\n\\t\\tif idx, ok := index[parent]; ok {\\n\\t\\t\\tfiles[idx].children = append(files[idx].children, i)\\n\\t\\t}\\n\\t\\ti++\\n\\t\\treturn nil\\n\\t})\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\n\\tfor _, ff := range files {\\n\\t\\t// search neighborhood\\n\\t\\tfor i := range ff.children {\\n\\t\\t\\tnext := -1\\n\\t\\t\\tif i+1 < len(ff.children) {\\n\\t\\t\\t\\tnext = ff.children[i+1]\\n\\t\\t\\t}\\n\\t\\t\\tfiles[ff.children[i]].next = next\\n\\t\\t}\\n\\n\\t\\tfmt.Fprintf(f, \\\"\\\\tfile{\\\\n\\\")\\n\\t\\trel, err := filepath.Rel(in, ff.path)\\n\\t\\tif err != nil {\\n\\t\\t\\treturn err\\n\\t\\t}\\n\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tname:    %%q,\\\\n\\\", path.Clean(\\\"/\\\"+filepath.ToSlash(rel)))\\n\\t\\tif ff.mode.IsDir() {\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tcontent: \\\\\\\"\\\\\\\",\\\")\\n\\t\\t} else {\\n\\t\\t\\tb, err := ioutil.ReadFile(ff.path)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn err\\n\\t\\t\\t}\\n\\t\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tcontent: %%q,\\\\n\\\", string(b))\\n\\t\\t}\\n\\t\\tswitch {\\n\\t\\tcase ff.mode.IsDir(): // directory\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0755 | os.ModeDir,\\\")\\n\\t\\tcase ff.mode&0100 != 0: // executable file\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0755,\\\")\\n\\t\\tdefault:\\n\\t\\t\\tfmt.Fprintln(f, \\\"\\\\t\\\\tmode:    0644,\\\")\\n\\t\\t}\\n\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tnext:    %%d,\\\\n\\\", ff.next)\\n\\t\\tif len(ff.children) > 0 {\\n\\t\\t\\tfmt.Fprintf(f, \\\"\\\\t\\\\tchild:   %%d,\\\\n\\\", ff.children[0])\\n\\t\\t} else {\\n\\t\\t\\tfmt.Fprint(f, \\\"\\\\t\\\\tchild:   -1,\\\\n\\\")\\n\\t\\t}\\n\\t\\tfmt.Fprint(f, \\\"\\\\t},\\\\n\\\")\\n\\t}\\n\\tfooter := %c%s%c\\n\\tfmt.Fprintln(f, footer)\\n\\tif err := f.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\n\\tf, err = os.OpenFile(filepath.Join(out, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)\\n\\tformat := %c%s%c\\n\\tfmt.Fprintf(f, format, version, 96, header, 96, 96, footer, 96, 96, format, 96)\\n\\tif err := f.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn nil\\n}\\n`\\n\\tfmt.Fprintf(f, format, version, 96, header, 96, 96, footer, 96, 96, format, 96)\\n\\tif err := f.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn nil\\n}\\n\",\n\t\tmode:    0644,\n\t\tnext:    2,\n\t\tchild:   -1,\n\t},\n\tfile{\n\t\tname:    \"/filesystem.go\",\n",
		mode:    0644,
		next:    3,
		child:   -1,
	},
	file{
		name:    "/haml.tmpl",
		content: "%meta{:charset => \"utf-8\"}\n%meta{\"http-equiv\"=>\"X-UA-Compatible\", :content=>\"IE=edge\"}\n%title {{ .Title }}\n%meta{:name => \"viewport\", :content => \"width=device-width,initial-scale=1\"}\n%meta{:name => \"title\", :content => \"{{ .Title }}\"}\n%meta{:name => \"description\", :content => \"{{ .Description }}\"}\n%meta{:name => \"keywords\", :content => \"{{ .Keywords }}\"}\n\n/ Open Graph Meta Tags\n%meta{:property => \"og:type\", :content => \"website\"}\n%meta{:property => \"og:url\", :content => \"{{ .Url }}\"}\n%meta{:property => \"og:title\", :content => \"{{ .Title }}\"}\n%meta{:property => \"og:description\", :content => \"{{ .Description }}\"}\n%meta{:property => \"og:image\", :content => \"{{ .ImgPath }}\"}\n\n/ Twitter\n%meta{:property => \"twitter:card\", :content => \"summary_large_image\"}\n%meta{:property => \"twitter:url\", :content => \"{{ .Url }}\"}\n%meta{:property => \"twitter:title\", :content => \"{{ .Title }}\"}\n%meta{:property => \"twitter:description\", :content => \"{{ .Description }}\"}\n%meta{:property => \"twitter:image\", :content => \"{{ .ImgPath }}\"}",
		mode:    0644,
		next:    4,
		child:   -1,
	},
	file{
		name:    "/html.tmpl",
		content: "<meta charset=\"utf-8\">\n<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n<title>{{ .Title }}</title>\n<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\">\n<meta name=\"title\" content=\"{{ .Title }}\">\n<meta name=\"description\" content=\"{{ .Description }}\">\n<meta name=\"keywords\" content=\"{{ .Keywords }}\">\n\n<!-- Open Graph Meta Tags -->\n<meta property=\"og:type\" content=\"website\">\n<meta property=\"og:url\" content=\"{{ .Url }}\">\n<meta property=\"og:title\" content=\"{{ .Title }}\">\n<meta property=\"og:description\" content=\"{{ .Description }}\">\n<meta property=\"og:image\" content=\"{{ .ImgPath }}\">\n\n<!-- Twitter -->\n<meta property=\"twitter:card\" content=\"summary_large_image\">\n<meta property=\"twitter:url\" content=\"{{ .Url }}\">\n<meta property=\"twitter:title\" content=\"{{ .Title }}\">\n<meta property=\"twitter:description\" content=\"{{ .Description }}\">\n<meta property=\"twitter:image\" content=\"{{ .ImgPath }}\">",
		mode:    0644,
		next:    5,
		child:   -1,
	},
	file{
		name:    "/jsx.tmpl",
		content: "<meta charSet=\"utf-8\" />\n<meta httpEquiv=\"X-UA-Compatible\" content=\"IE=edge\" />\n<title>{{ .Title }}</title>\n<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" />\n<meta name=\"title\" content=\"{{ .Title }}\" />\n<meta name=\"description\" content=\"{{ .Description }}\" />\n<meta name=\"keywords\" content=\"{{ .Keywords }}\" />\n\n{/* Open Graph Meta Tags */}\n<meta property=\"og:type\" content=\"website\" />\n<meta property=\"og:url\" content=\"{{ .Url }}\" />\n<meta property=\"og:title\" content=\"{{ .Title }}\" />\n<meta property=\"og:description\" content=\"{{ .Description }}\" />\n<meta property=\"og:image\" content=\"{{ .ImgPath }}\" />\n\n{/* Twitter */}\n<meta property=\"twitter:card\" content=\"summary_large_image\" />\n<meta property=\"twitter:url\" content=\"{{ .Url }}\" />\n<meta property=\"twitter:title\" content=\"{{ .Title }}\" />\n<meta property=\"twitter:description\" content=\"{{ .Description }}\" />\n<meta property=\"twitter:image\" content=\"{{ .ImgPath }}\" />",
		mode:    0644,
		next:    6,
		child:   -1,
	},
	file{
		name:    "/pug.tmpl",
		content: "meta(charset=\"utf-8\")\nmeta(http-equiv=\"X-UA-Compatible\" content=\"IE=edge\")\ntitle {{ .Title }}\nmeta(name=\"viewport\" content=\"width=device-width,initial-scale=1\")\nmeta(name=\"title\" content=\"{{ .Title }}\")\nmeta(name=\"description\" content=\"{{ .Description }}\")\nmeta(name=\"keywords\" content=\"{{ .Keywords }}\")\n\n// Open Graph Meta Tags\nmeta(property=\"og:type\" content=\"website\")\nmeta(property=\"og:url\" content=\"{{ .Url }}\")\nmeta(property=\"og:title\" content=\"{{ .Title }}\")\nmeta(property=\"og:description\" content=\"{{ .Description }}\")\nmeta(property=\"og:image\" content=\"{{ .ImgPath }}\")\n\n// Twitter\nmeta(property=\"twitter:card\" content=\"summary_large_image\")\nmeta(property=\"twitter:url\" content=\"{{ .Url }}\")\nmeta(property=\"twitter:title\" content=\"{{ .Title }}\")\nmeta(property=\"twitter:description\" content=\"{{ .Description }}\")\nmeta(property=\"twitter:image\" content=\"{{ .ImgPath }}\")",
		mode:    0644,
		next:    7,
		child:   -1,
	},
	file{
		name:    "/slim.tmpl",
		content: "meta charset=\"utf-8\"\nmeta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"\ntitle {{ .Title }}\nmeta name=\"viewport\" content=\"width=device-width,initial-scale=1\"\nmeta name=\"title\" content=\"{{ .Title }}\"\nmeta name=\"description\" content=\"{{ .Description }}\"\nmeta name=\"keywords\" content=\"{{ .Keywords }}\"\n\n// Open Graph Meta Tags\nmeta property=\"og:type\" content=\"website\"\nmeta property=\"og:url\" content=\"{{ .Url }}\"\nmeta property=\"og:title\" content=\"{{ .Title }}\"\nmeta property=\"og:description\" content=\"{{ .Description }}\"\nmeta property=\"og:image\" content=\"{{ .ImgPath }}\"\n\n// Twitter\nmeta property=\"twitter:card\" content=\"summary_large_image\"\nmeta property=\"twitter:url\" content=\"{{ .Url }}\"\nmeta property=\"twitter:title\" content=\"{{ .Title }}\"\nmeta property=\"twitter:description\" content=\"{{ .Description }}\"\nmeta property=\"twitter:image\" content=\"{{ .ImgPath }}\"",
		mode:    0644,
		next:    -1,
		child:   -1,
	},
}

type fileSystem []file

func (fs fileSystem) Open(name string) (http.File, error) {
	name = path.Clean("/" + name)
	i := sort.Search(len(fs), func(i int) bool { return fs[i].name >= name })
	if i >= len(fs) || fs[i].name != name {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  os.ErrNotExist,
		}
	}
	f := &fs[i]
	return &httpFile{
		Reader: strings.NewReader(f.content),
		file:   f,
		fs:     fs,
		idx:    i,
		dirIdx: f.child,
	}, nil
}

type file struct {
	name    string
	content string
	mode    os.FileMode
	child   int
	next    int
}

var _ os.FileInfo = (*file)(nil)

func (f *file) Name() string {
	return path.Base(f.name)
}

func (f *file) Size() int64 {
	return int64(len(f.content))
}

func (f *file) Mode() os.FileMode {
	return f.mode
}

var zeroTime time.Time

func (f *file) ModTime() time.Time {
	return zeroTime
}

func (f *file) IsDir() bool {
	return f.Mode().IsDir()
}

func (f *file) Sys() interface{} {
	return nil
}

type httpFile struct {
	*strings.Reader
	file   *file
	fs     fileSystem
	idx    int
	dirIdx int
}

var _ http.File = (*httpFile)(nil)

func (f *httpFile) Stat() (os.FileInfo, error) {
	return f.file, nil
}

func (f *httpFile) Readdir(count int) ([]os.FileInfo, error) {
	ret := []os.FileInfo{}
	if !f.file.IsDir() {
		return ret, nil
	}

	if count <= 0 {
		for f.dirIdx >= 0 {
			entry := &f.fs[f.dirIdx]
			ret = append(ret, entry)
			f.dirIdx = entry.next
		}
		return ret, nil
	}

	ret = make([]os.FileInfo, 0, count)
	for f.dirIdx >= 0 {
		entry := &f.fs[f.dirIdx]
		ret = append(ret, entry)
		f.dirIdx = entry.next
		if len(ret) == count {
			return ret, nil
		}
	}
	return ret, io.EOF
}

func (f *httpFile) Close() error {
	return nil
}
