package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	var source string
	var objsStr string
	var output string
	var include string

	flag.StringVar(&source, "path", "", "")
	flag.StringVar(&objsStr, "objs", "", "")
	flag.StringVar(&output, "output", "", "")
	flag.StringVar(&include, "include", "", "")

	flag.Parse()

	targets := decodeList(objsStr)
	includeList := decodeList(include)

	if _, err := process(source, targets, output, includeList); err != nil {
		fmt.Printf("[ERR]: %v", err)
	}
}

func decodeList(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(input), ",")
}

func process(source string, targetObjs []string, output string, includePaths []string) (string, error) {
	// read ast file
	file, err := parser.ParseFile(token.NewFileSet(), source, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	// read package name
	packName := file.Name.Name

	// decode the struct fields
	structs := decodeStructs(file)

	// get the target structs
	targets := []*astStruct{}
	if len(targetObjs) == 0 {
		targets = append(targets, structs...)
	} else {
		for _, i := range structs {
			if contains(targetObjs, i.name) {
				targets = append(targets, i)
			}
		}
	}

	// process the targets
	outs := []string{}
	for _, s := range targets {
		v, err := createIR(s, structs)
		if err != nil {
			return "", err
		}
		outs = append(outs,
			generateMarshal(v),
			"",
			generateUnmarshal(v),
		)
	}

	tmpl := `package {{.Package}}
	
	import (
		"fmt"

		"github.com/umbracle/fastrlp"
	)

	{{.Outs}}
	`
	res := execTmpl(tmpl, map[string]interface{}{
		"Package": packName,
		"Outs":    strings.Join(outs, "\n"),
	})

	// escape the &
	res = strings.Replace(res, "xx", "&", -1)

	resB := []byte(res)

	/*
		// format
		resB, err = format.Source(resB)
		if err != nil {
			return "", err
		}
	*/

	// write output file
	if err := ioutil.WriteFile(output, resB, 0644); err != nil {
		return "", err
	}

	return string(resB), nil
}

func isExportedField(str string) bool {
	return str[0] <= 90
}

func isByte(obj ast.Expr) bool {
	if ident, ok := obj.(*ast.Ident); ok {
		if ident.Name == "byte" {
			return true
		}
	}
	return false
}

func getStruct(name string, contexts []*astStruct) *astStruct {
	for _, i := range contexts {
		if i.name == name {
			return i
		}
	}
	return nil
}

func getObjLen(obj *ast.ArrayType) uint64 {
	if obj.Len == nil {
		return 0
	}
	value := obj.Len.(*ast.BasicLit).Value
	num, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		panic(fmt.Sprintf("BUG: Failed to convert to uint64 %s: %v", value, err))
	}
	return num
}

func createIR(target *astStruct, contexts []*astStruct) (*value, error) {
	if target.obj == nil {
		return nil, fmt.Errorf("cannot do it for alias values")
	}

	fields := []*field{}
	for _, f := range target.obj.Fields.List {
		if len(f.Names) != 1 {
			continue
		}
		name := f.Names[0].Name
		if !isExportedField(name) {
			continue
		}

		f := parseASTField(f.Type, contexts)
		f.name = name
		f.indx = uint64(len(fields))

		fields = append(fields, f)
	}

	v := &value{
		fields: fields,
		name:   target.name,
	}
	return v, nil
}

func findStructByName(contexts []*astStruct, name string) *astStruct {
	for _, obj := range contexts {
		if obj.name == name {
			return obj
		}
	}
	return nil
}

func parseStructRef(name string, contexts []*astStruct) *field {
	// struct reference
	target := findStructByName(contexts, name)
	if target == nil {
		panic("NOT FOUND")
	}
	if target.typ != nil {
		// alias of basic type
		return parseASTField(target.typ, nil)
	}
	// another struct
	return &field{typ: typeObj, obj: name}
}

func parseASTField(expr ast.Expr, contexts []*astStruct) *field {
	switch obj := expr.(type) {
	case *ast.ArrayType:
		if isByte(obj.Elt) {
			if fixedlen := getObjLen(obj); fixedlen != 0 {
				// [fixedLen]byte
				return &field{typ: typeFixedBytes, len: fixedlen}
			}
			// []byte
			return &field{typ: typeDynamicBytes}
		}
		// []*Struct
		vv := parseASTField(obj.Elt, contexts)
		return &field{typ: typeArray, obj: vv.obj}

	case *ast.Ident:
		name := obj.Name
		if name == "uint64" {
			// uint64
			return &field{typ: typeUint}
		} else if name == "byte" {
			// byte
			return &field{typ: typeByte}
		}
		// struct reference
		return parseStructRef(name, contexts)

	case *ast.StarExpr:
		// *Struct
		var f *field
		switch elem := obj.X.(type) {
		case *ast.Ident:
			// reference to a local package
			f = parseStructRef(elem.Name, contexts)
		}

		if f.typ != typeObj {
			f.ptr = true
		}
		return f
	}
	panic(fmt.Sprintf("BUG: %s", reflect.TypeOf(expr)))
}

func contains(i []string, j string) bool {
	for _, o := range i {
		if o == j {
			return true
		}
	}
	return false
}

type astStruct struct {
	name string
	obj  *ast.StructType
	typ  ast.Expr
}

func decodeStructs(file *ast.File) []*astStruct {
	res := []*astStruct{}

	for _, dec := range file.Decls {
		if genDecl, ok := dec.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					obj := &astStruct{
						name: typeSpec.Name.Name,
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok {
						obj.obj = structType
					} else {
						obj.typ = typeSpec.Type
					}
					res = append(res, obj)
				}
			}
		}
	}
	return res
}

// --- IR ---

// fieldType if the type of the field
type fieldType int

const (
	typeUint fieldType = iota + 1
	typeFixedBytes
	typeDynamicBytes
	typeByte
	typeArray // array of non-fixed length of objects
	typeObj   // reference to another object
)

// field represents a field in a value
type field struct {
	name string    // name of the field
	obj  string    // name of the struct
	typ  fieldType // type of the field
	len  uint64    // len != 0 if the field has fixed size
	ptr  bool      // the field is a pointer (only possible it not typeArray)
	indx uint64    // indes of the field in the struct
}

// value represents a Go struct
type value struct {
	name   string
	fields []*field
}

// --- Generation ---

func generateUnmarshal(v *value) string {
	tmpl := `func (:: *{{.Name}}) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := ::.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (:: *{{.Name}}) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != {{.Len}} {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		{{range $val := .Fields}}
		{{$val}}
		{{end}}
		return nil
	}`

	fieldsOut := []string{}
	for _, f := range v.fields {
		tmpl := generateUnmarshalField(f)
		fieldOut := execTmpl(tmpl, map[string]interface{}{
			"Name": f.name,
			"Indx": f.indx,
			"Len":  f.len,
			"Obj":  f.obj,
		})

		out := fmt.Sprintf("// Field '%s'\n%s", f.name, fieldOut)
		fieldsOut = append(fieldsOut, out)
	}

	return execTmpl(tmpl, map[string]interface{}{
		"Name":   v.name,
		"Fields": fieldsOut,
		"Len":    len(v.fields),
		"Sig":    getSignature(v),
	})
}

func generateUnmarshalField(f *field) string {
	switch f.typ {
	case typeUint:
		return `if ::.{{.Name}}, err = elems[{{.Indx}}].GetUint64(); err != nil {
			return err
		}`

	case typeFixedBytes:
		var tmpl string
		if f.len == 32 {
			// hash
			tmpl = `if err = elems[{{.Indx}}].GetHash(::.{{.Name}}[:]); err != nil {
				return err
			}`
		} else if f.len == 20 {
			// address
			tmpl = `if err = elems[{{.Indx}}].GetAddr(::.{{.Name}}[:]); err != nil {
				return err
			}`
		} else {
			tmpl = `if _, err = elems[{{.Indx}}].GetBytes(::.{{.Name}}[:0], {{.Len}}); err != nil {
				return err
			}`
		}
		return tmpl

	case typeDynamicBytes:
		return `if ::.{{.Name}}, err = elems[{{.Indx}}].GetBytes(::.{{.Name}}[:0]); err != nil {
			return err
		}`

	case typeByte:
		return `if ::.{{.Name}}, err = elems[{{.Indx}}].GetByte(); err != nil {
			return err
		}`

	case typeObj:
		return `{
			::.{{.Name}} = xx{{.Obj}}{}
			if err := ::.{{.Name}}.UnmarshalRLPFrom(elems[{{.Indx}}]); err != nil {
				return err
			}
		}`

	case typeArray:
		return `{
			subElems, err := elems[{{.Indx}}].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := xx{{.Obj}}{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				::.{{.Name}} = append(::.{{.Name}}, bb)
			}
		}`

	default:
		panic("TODO")
	}
}

func getSignature(v *value) string {
	return strings.ToLower(string(v.name[0]))
}

func generateMarshal(v *value) string {
	tmpl := `func (:: *{{.Name}}) MarshalRLP() []byte {
		return ::.MarshalRLPTo(nil)
	}
	
	func (:: *{{.Name}}) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = ::.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (:: *{{.Name}}) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		{{range $val := .Fields}}
		{{$val}}
		{{end}}
		return vv
	}`

	fieldsOut := []string{}
	for _, f := range v.fields {
		fieldOut := execTmpl(generateMarshalField(f), map[string]interface{}{
			"Name": f.name,
			"Indx": f.indx,
			"Len":  f.len,
			"Obj":  f.obj,
		})

		if f.ptr {
			// if its a pointer, the encoding sets null if the
			// value does not exists
			fieldOut = `if ::.` + f.name + ` == nil {
                vv.Set(ar.NewNull())
            } else {
                ` + fieldOut + `
            }`
		}

		out := fmt.Sprintf("// Field '%s'\n%s", f.name, fieldOut)
		fieldsOut = append(fieldsOut, out)
	}

	return execTmpl(tmpl, map[string]interface{}{
		"Name":   v.name,
		"Fields": fieldsOut,
		"Sig":    getSignature(v),
	})
}

func generateMarshalField(f *field) string {
	switch f.typ {
	case typeUint:
		return "vv.Set(ar.NewUint(::.{{.Name}}))"

	case typeFixedBytes:
		return "vv.Set(ar.NewBytes(::.{{.Name}}[:]))"

	case typeDynamicBytes:
		return "vv.Set(ar.NewCopyBytes(::.{{.Name}}))"

	case typeByte:
		return "vv.Set(ar.NewUint(uint64(::.{{.Name}})))"

	case typeObj:
		return "vv.Set(::.{{.Name}}.MarshalRLPWith(ar))"

	case typeArray:
		return `{
			if len(::.{{.Name}}) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range ::.{{.Name}} {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}`

	default:
		panic("TODO")
	}
}

func execTmpl(tpl string, input map[string]interface{}) string {
	tmpl, err := template.New("tmpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, input); err != nil {
		panic(err)
	}
	out := buf.String()

	// escape quotes
	out = strings.Replace(out, "&#34;", "\"", -1)
	out = strings.Replace(out, "&amp;#39;", "'", -1)

	sig, ok := input["Sig"]
	if ok {
		out = strings.Replace(out, "::", sig.(string), -1)
	}
	return out
}
