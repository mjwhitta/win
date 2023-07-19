package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/pathname"
)

type cacheEntry struct {
	C    string // C-style var name
	Go   string // CamelCase var name
	Type string // Type of var
	Val  string // Value for var
}

// Cache, indexed by scope
var cache = map[string][]*cacheEntry{
	"": { // Global scope
		{"FALSE", "False", "uintptr", "0"},
		{"NULL", "Null", "uintptr", "0"},
		{"TRUE", "True", "uintptr", "1"},
	},
}
var lookup = map[string]*cacheEntry{}

// Regular expressions
var bitwisenot = regexp.MustCompile(`\~`)
var camel = regexp.MustCompile(`[A-Z][a-z]+[A-Z][a-z]+`)
var comment = regexp.MustCompile(`\s*\/\*.*\*\/\s*`)
var fixcast = regexp.MustCompile(
	`\([A-Za-z_]+\)\s*(\(?\-?[0-9A-Fa-fXx]+|NULL)`,
)
var fixhex = regexp.MustCompile(`(0[Xx][0-9A-Fa-f]+)[LlUu]+`)
var fixlen = regexp.MustCompile(`(len\(.+?\))`)
var fixmsabi = regexp.MustCompile(`__MSABI_LONG([^)]+)`)
var fixnum = regexp.MustCompile(`(\d+)[LlUu]+`)
var fixsizeof = regexp.MustCompile(`sizeof\s*\(`)
var spaces = regexp.MustCompile(`\s+`)
var uselessOpRepl = []*regexp.Regexp{
	regexp.MustCompile(`\s*\^\s*0([^x])`), // ^ 0
	regexp.MustCompile(`\s*\|\s*0([^x])`), // | 0
}
var uselessOpRm = []*regexp.Regexp{
	regexp.MustCompile(`\s*\^\s*\(0\)`),             // ^ (0)
	regexp.MustCompile(`\s*\|\s*\(0\)`),             // | (0)
	regexp.MustCompile(`(0x0{8}\s*\|\s*)+`),         // 0x00000000 |
	regexp.MustCompile(`(\(0x0{8}\)\s*\|\s*)+`),     // (0x00000000) |
	regexp.MustCompile(`(\(\(0x0{8}\)\)\s*\|\s*)+`), // ((0x0...0)) |
	regexp.MustCompile(`(\s*\|\s*0x0{8})+`),         // | 0x00000000
	regexp.MustCompile(`(\s*\|\s*\(0x0{8}\))+`),     // | (0x00000000)
	regexp.MustCompile(`(\s*\|\s*\(\(0x0{8}\)\))+`), // | ((0x0...0))
}

// Skips
var skipLContains = map[string][]string{
	"": {
		"(",
		")",
		"DECLSPEC",
		"EXTERN_C",
	},
	"shellapi.h": {
		"DUMMY",
	},
	"winnt.h": {
		"DUMMY",
		"XSTATE_MASK_ALLOWED",
	},
}
var skipRContains = map[string][]string{
	"": {
		"__declspec",
		"__MINGW_NAME",
		"DECLSPEC",
		"HRESULT",
		"len(DWORD)",
		"len(ULONGLONG)",
		"WINAPI",
	},
	"shellapi.h": {
		"FIELD_OFFSET",
	},
	"wininet.h": {
		"INTERNET_STATUS_CALLBACK",
	},
	"winnt.h": {
		"FIELD_OFFSET",
		"inline",
		"MAKELANGID(",
		"MAKELCID(",
	},
	"nb30.h": {
		"\\0",
	},
	"wincrypt.h": {
		"\\0",
	},
	"winuser.h": {
		"len(LRESULT)",
		"MAKEINTATOM(",
	},
}
var skipRStarts = map[string][]string{
	"": {
		":",
		"_",
		"extern",
		"void",
	},
	"ddeml.h": {
		"CALLBACK",
	},
	"mmsystem.h": {
		"mmioFOURCC",
		"OutputDebugString",
		"sndAlias(",
	},
	"rpc.h": {
		"MIDL_user",
		"struct",
		"}",
	},
	"shellapi.h": {
		"SHGetDiskFreeSpaceEx",
		"STDAPI",
	},
	"wincrypt.h": {
		"const",
	},
	"winscard.h": {
		"(&g_rgSCard",
		"SCardGetAttrib",
		"SCardSetAttrib",
	},
	"winuser.h": {
		"GET_DEVICE_LPARAM",
		"(UINT_MAX)",
	},
}

func buildLookup() {
	for scope, entries := range cache {
		if scope == "" {
			continue
		}

		for _, entry := range entries {
			lookup[format(scope)+"."+entry.Go] = entry
		}
	}
}

// This only caches the first result, which may not be the best idea.
func cacheVar(fn string, c string, v string) {
	var g string
	var scope string = strings.TrimSuffix(filepath.Base(fn), ".h")
	var t string = "uintptr"

	// Remove leading and trailing whitespace, just in case
	c = strings.TrimSpace(c)
	g = format(c)
	v = strings.TrimSpace(v)

	// Check Go-style not C-style b/c multiple C vars might resolve to
	// the same Go var
	for _, entry := range cache[scope] {
		if entry.Go == g {
			return
		}
	}

	if fixlen.MatchString(v) {
		v = fixlen.ReplaceAllString(v, "uintptr($1)")
	}

	if strings.HasPrefix(v, "\"") {
		t = "string"
	} else if strings.HasPrefix(v, "L\"") {
		v = v[1:]
		t = "string"
	} else if strings.HasPrefix(v, "TEXT(") {
		v = strings.Replace(v[5:], ")", "", 1)
		t = "string"
	} else if strings.Contains(v, "L\"") {
		v = "[]string{" + strings.ReplaceAll(v, " L\"", ", \"") + "}"
		t = "[]string"
	} else if strings.Contains(v, "-") {
		t = "int"
	} else if strings.Contains(v, ".") && strings.HasSuffix(v, "f") {
		v = strings.TrimSuffix(v, "f")
		t = "float64"
	} else if strings.HasPrefix(v, "{") {
		v = "[]uintptr" + v
		t = "[]uintptr"
	}

	// No previous entry found
	cache[scope] = append(cache[scope], &cacheEntry{c, g, t, v})
}

func fixVarTypes() {
	for _, entries := range cache {
		for _, entry := range entries {
			// Fix types
			if strings.HasPrefix(entry.Val, "\"") {
				entry.Type = "string"
			} else if strings.HasPrefix(entry.Val, "[]string{") {
				entry.Type = "[]string"
			} else if strings.HasPrefix(entry.Val, "[]uintptr{") {
				entry.Type = "[]uintptr"
			} else if strings.Contains(entry.Val, "-") {
				entry.Type = "int"
			}
		}
	}
}

func format(str string) string {
	var tmp []string

	if str == "" {
		return str
	}

	// Replace _ with CamelCase
	if strings.Contains(str, "_") {
		// Split on "_"
		tmp = strings.Split(strings.ToLower(str), "_")

		// Capitalize every part
		for i := range tmp {
			if tmp[i] == "" {
				continue
			}

			tmp[i] = strings.ToUpper(tmp[i][:1]) + tmp[i][1:]
		}

		// Join together for camelcase
		str = strings.Join(tmp, "")
	} else if camel.MatchString(str) {
		// Do nothing
	} else {
		str = strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
	}

	// Fix some special cases
	str = strings.ReplaceAll(str, "Crlf", "CRLF")
	str = strings.ReplaceAll(str, "Ftp", "FTP")
	str = strings.ReplaceAll(str, "Http", "HTTP")
	str = strings.ReplaceAll(str, "LOGICAL", "Logical")
	str = strings.ReplaceAll(str, "LOGNAME", "Logname")

	// Fix type-related
	str = strings.ReplaceAll(str, "BYTE", "Byte")
	str = strings.ReplaceAll(str, "CHAR", "Char")
	str = strings.ReplaceAll(str, "DWORD", "Dword")
	str = strings.ReplaceAll(str, "MAX", "Max")
	str = strings.ReplaceAll(str, "MIN", "Min")
	str = strings.ReplaceAll(str, "LONG", "Long")
	str = strings.ReplaceAll(str, "SHORT", "Short")
	str = strings.ReplaceAll(str, "SIZE", "Size")
	str = strings.ReplaceAll(str, "WORD", "Word")

	return str
}

func genFile(pkg string) error {
	var e error
	var entries []*cacheEntry
	var f *os.File
	var scopes []string

	// Build lookup table and replace C-style var names with Go-style
	// var names
	buildLookup()
	replaceVars()
	fixVarTypes()

	// Open file to write
	if f, e = genHeader(pkg); e != nil {
		return e
	}
	defer f.Close()

	// Get all scopes
	for scope := range cache {
		if scope == "" {
			continue
		}

		scopes = append(scopes, scope)
	}

	// Sort alphabetically, case-insensitive
	sort.Slice(
		scopes,
		func(i int, j int) bool {
			var l string = strings.ToLower(scopes[i])
			var r string = strings.ToLower(scopes[j])

			if l == r {
				return scopes[i] < scopes[j]
			}

			return l < r
		},
	)

	// Loop thru scopes
	for _, scope := range scopes {
		entries = cache[scope]

		// Sort alphabetically, case-insensitive
		sort.Slice(
			entries,
			func(i int, j int) bool {
				var l string = strings.ToLower(entries[i].Go)
				var r string = strings.ToLower(entries[j].Go)

				if l == r {
					return entries[i].Go < entries[j].Go
				}

				return l < r
			},
		)

		// Start scoped section
		f.WriteString(
			"\n// " + format(scope) + " contains constants from " +
				scope + ".h\n",
		)
		f.WriteString("var " + format(scope) + " = struct {\n")

		// Define struct
		for _, entry := range entries {
			f.WriteString(
				fmt.Sprintf("\t%s %s\n", entry.Go, entry.Type),
			)
		}

		f.WriteString("}{\n")

		// Create struct
		for _, entry := range entries {
			f.WriteString(
				fmt.Sprintf("\t%s: %s,\n", entry.Go, entry.Val),
			)
		}

		f.WriteString("}\n")
	}

	return nil
}

func genHeader(pkg string) (*os.File, error) {
	var e error
	var f *os.File
	var fn string = "generated.go"

	// Open file to write
	if f, e = os.Create(fn); e != nil {
		return nil, errors.Newf("failed to create %s: %w", fn, e)
	}

	// Create header
	f.WriteString(
		"// Code generated by tools/defines.go; DO NOT EDIT.\n",
	)
	f.WriteString("package " + pkg + "\n\n")

	sort.Slice(
		cache[""], // Global scope
		func(i int, j int) bool {
			var l string = strings.ToLower(cache[""][i].Go)
			var r string = strings.ToLower(cache[""][j].Go)

			return l < r
		},
	)

	f.WriteString("const (\n")

	for _, entry := range cache[""] {
		f.WriteString(
			fmt.Sprintf(
				"\t%s %s = %s\n",
				entry.Go,
				entry.Type,
				entry.Val,
			),
		)
	}

	f.WriteString(")\n")

	return f, nil
}

func ignoreType(fn string, l string, sep string) {
	for _, t := range strings.Split(l, sep) {
		t = strings.TrimSpace(t)
		t = strings.TrimPrefix(t, "*")
		skipRContains[fn] = append(skipRContains[fn], "len("+t+")")
		skipRStarts[fn] = append(skipRStarts[fn], t)
	}
}

func init() {
	flag.Parse()
}

// This is by no means perfect, but I try to grab as many constants as
// possible.
func main() {
	if flag.NArg() == 0 {
		return
	}

	// Find all the things to ignore/skip first (probably can remove
	// later after implementing struct parsing)
	for i, arg := range flag.Args() {
		if i == 0 {
			continue
		}

		arg = "/usr/x86_64-w64-mingw32/include/" + arg

		if ok, e := pathname.DoesExist(arg); e != nil {
			fmt.Println(e.Error())
			os.Exit(1)
		} else if !ok {
			fmt.Printf("file %s not found\n", arg)
			os.Exit(1)
		}

		if e := processFileSkips(arg); e != nil {
			panic(e)
		}
	}

	// Then process each file for #defines and enums
	for i, arg := range flag.Args() {
		if i == 0 {
			continue
		}

		arg = "/usr/x86_64-w64-mingw32/include/" + arg

		if e := processFileDefines(arg); e != nil {
			panic(e)
		}

		// TODO process structs too

		if e := processFileTypedefs(arg); e != nil {
			panic(e)
		}
	}

	if e := genFile(flag.Arg(0)); e != nil {
		panic(e)
	}
}

func processDefine(fn string, l string) {
	var tmp []string

	if !strings.HasPrefix(l, "#define") {
		return
	}

	// Remove define
	l = strings.Replace(l, "#define ", "", 1)

	// Remove comments
	if comment.MatchString(l) {
		l = comment.ReplaceAllString(l, "")
	}

	// Fix some Windows-isms
	if bitwisenot.MatchString(l) {
		l = bitwisenot.ReplaceAllString(l, "0xffffffff^")
	}

	if fixsizeof.MatchString(l) {
		l = fixsizeof.ReplaceAllString(l, "len(")
	}

	for _, r := range []*regexp.Regexp{
		fixhex,
		fixcast,
		fixmsabi,
		fixnum,
	} {
		if r.MatchString(l) {
			l = r.ReplaceAllString(l, "$1")
		}
	}

	// Split and attempt to skip things we don't want
	if tmp = strings.SplitN(l, " ", 2); !skip(fn, tmp) {
		cacheVar(fn, tmp[0], tmp[1])
	}
}

func processFileDefines(fp string) error {
	var b []byte
	var e error
	var f *os.File
	var fn string = filepath.Base(fp)
	var tmp string

	if f, e = os.Open(pathname.ExpandPath(fp)); e != nil {
		return errors.Newf("failed to open %s: %w", fp, e)
	}
	defer f.Close()

	if b, e = io.ReadAll(f); e != nil {
		return errors.Newf("failed to read %s: %w", fp, e)
	}

	for _, l := range strings.Split(string(b), "\n") {
		l = strings.ReplaceAll(l, "| INTERNET_FLAG_BGUPDATE", "")
		l = strings.TrimSpace(spaces.ReplaceAllString(l, " "))

		if strings.HasSuffix(l, "\\") {
			tmp += l[:len(l)-1]
			continue
		}

		processDefine(fn, strings.TrimSpace(tmp+l))

		// Reset
		tmp = ""
	}

	return nil
}

func processFileSkips(fp string) error {
	var b []byte
	var e error
	var f *os.File
	var fn string = filepath.Base(fp)
	var inStructOrTypedef bool

	if f, e = os.Open(pathname.ExpandPath(fp)); e != nil {
		return errors.Newf("failed to open %s: %w", fp, e)
	}
	defer f.Close()

	if b, e = io.ReadAll(f); e != nil {
		return errors.Newf("failed to read %s: %w", fp, e)
	}

	for _, l := range strings.Split(string(b), "\n") {
		l = strings.TrimSpace(spaces.ReplaceAllString(l, " "))

		if strings.HasPrefix(l, "#define") {
			l = strings.TrimPrefix(l, "#define")
			l = strings.TrimSpace(l)

			if !strings.Contains(l, " ") && (len(l) > 8) {
				skipRContains[fn] = append(skipRContains[fn], l)
			}
		} else if strings.HasPrefix(l, "}") {
			if inStructOrTypedef {
				// Extract type
				l = strings.TrimPrefix(l, "} ")
				l = strings.TrimSuffix(l, ";")

				ignoreType(fn, l, ",")
			}

			inStructOrTypedef = false
		} else if strings.HasPrefix(l, "typedef const") {
			// Extract type
			l = strings.TrimPrefix(l, "typedef const")
			l = strings.TrimPrefix(l, "struct")
			l = strings.TrimSuffix(l, ";")
			l = strings.TrimSpace(l)

			if l != "" {
				ignoreType(fn, l, " ")
			}
		} else if strings.HasPrefix(l, "typedef enum") ||
			strings.HasPrefix(l, "typedef struct") {
			inStructOrTypedef = true

			// Extract type
			l = strings.TrimPrefix(l, "typedef enum")
			l = strings.TrimPrefix(l, "typedef struct")
			l = strings.TrimSuffix(l, "{")
			l = strings.TrimSpace(l)

			if l != "" {
				ignoreType(fn, l, ",")
			}
		}
	}

	return nil
}

func processFileTypedefs(fp string) error {
	var b []byte
	var e error
	var f *os.File
	var fn string = filepath.Base(fp)
	var tmp string
	var inTypedef bool

	if f, e = os.Open(pathname.ExpandPath(fp)); e != nil {
		return errors.Newf("failed to open %s: %w", fp, e)
	}
	defer f.Close()

	if b, e = io.ReadAll(f); e != nil {
		return errors.Newf("failed to read %s: %w", fp, e)
	}

	for _, l := range strings.Split(string(b), "\n") {
		l = strings.TrimSpace(spaces.ReplaceAllString(l, " "))

		// Ignore non-defines unless typedef enum
		if l == "{" {
			continue
		} else if strings.HasPrefix(l, "#") {
			// #ifdef and #endif
			continue
		} else if strings.HasPrefix(l, "{") {
			if inTypedef {
				tmp += l[1:]
			}

			continue
		} else if strings.HasPrefix(l, "}") {
			if !inTypedef {
				continue
			}

			inTypedef = false
		} else if strings.HasPrefix(l, "typedef enum") {
			inTypedef = true
			continue
		} else {
			if inTypedef {
				tmp += l
			}

			continue
		}

		// Remove comments
		if comment.MatchString(tmp) {
			tmp = comment.ReplaceAllString(tmp, "")
		}

		processTypedef(fn, tmp)

		// Reset
		tmp = ""
	}

	return nil
}

func processTypedef(fn string, l string) {
	var lhs string
	var next string = "0"
	var rhs string
	var tmp []string
	var vals = map[string]string{}

	for _, d := range strings.Split(l, ",") {
		d = strings.TrimSpace(d)

		if d == "" {
			continue
		}

		if strings.Contains(d, "=") {
			tmp = strings.SplitN(d, "=", 2)

			lhs = format(strings.TrimSpace(tmp[0]))
			rhs = format(strings.TrimSpace(tmp[1]))

			// Determine if next needs to be looked up from prev vals
			next = rhs
			if _, ok := vals[next]; ok {
				next = vals[next]
			}
		} else {
			lhs = format(d)
			rhs = next
		}

		// Store in case we need to lookup
		vals[lhs] = rhs

		// Prepare next enum value
		if strings.HasPrefix(next, "0x") {
			i, _ := strconv.ParseInt(next[2:], 16, 64)
			next = fmt.Sprintf("0x%x", i+1)
		} else {
			i, _ := strconv.Atoi(next)
			next = fmt.Sprintf("%d", i+1)
		}

		// Kkip things we don't want
		if !skip(fn, []string{lhs, rhs}) {
			cacheVar(fn, lhs, rhs)
		}
	}
}

func replaceVars() {
	for scope, entries := range cache {
		if scope == "" {
			continue
		}

		// Sort C-style var names by longest to shortest
		sort.Slice(
			entries,
			func(i int, j int) bool {
				return len(entries[i].C) > len(entries[j].C)
			},
		)

		// Replace C-style var names with actual values
		for _, entry := range entries {
			// Dumb hack for recursive lookups
			for i := 0; i < 8; i++ {
				for _, repl := range entries {
					entry.Val = strings.ReplaceAll(
						entry.Val,
						repl.C,
						repl.Val,
					)
					entry.Val = strings.ReplaceAll(
						entry.Val,
						repl.Go,
						repl.Val,
					)
				}
			}

			// Global-scoped vars
			for _, repl := range cache[""] {
				entry.Val = strings.ReplaceAll(
					entry.Val,
					repl.C,
					repl.Go,
				)
			}

			// Other scopes via lookup table
			for ref, repl := range lookup {
				if !strings.HasPrefix(ref, scope) {
					entry.Val = strings.ReplaceAll(
						entry.Val,
						repl.C,
						ref,
					)
				}
			}

			for _, uselessOp := range uselessOpRepl {
				if uselessOp.MatchString(entry.Val) {
					entry.Val = uselessOp.ReplaceAllString(
						entry.Val,
						"$1",
					)
				}
			}

			for _, uselessOp := range uselessOpRm {
				if uselessOp.MatchString(entry.Val) {
					entry.Val = uselessOp.ReplaceAllString(
						entry.Val,
						"",
					)
				}
			}
		}
	}
}

func skip(fn string, tmp []string) bool {
	if len(tmp) != 2 {
		if len(tmp[0]) > 8 {
			skipRContains[fn] = append(skipRContains[fn], tmp[0])
		}

		return true
	}

	if strings.ToLower(tmp[1]) == strings.ToLower(tmp[0])+"w" {
		skipRContains[fn] = append(skipRContains[fn], tmp[0])
		return true
	}

	for _, scope := range []string{"", fn} {
		for _, s := range skipLContains[scope] {
			if strings.Contains(tmp[0], s) {
				skipRContains[fn] = append(skipRContains[fn], tmp[0])
				return true
			}
		}

		for _, s := range skipRContains[scope] {
			if strings.Contains(tmp[1], s) {
				skipRContains[fn] = append(skipRContains[fn], tmp[0])
				return true
			}
		}

		for _, s := range skipRStarts[scope] {
			if strings.HasPrefix(tmp[1], s) {
				skipRContains[fn] = append(skipRContains[fn], tmp[0])
				return true
			}
		}
	}

	return false
}
