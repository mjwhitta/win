//nolint:godox // I'll do the one todo later maybe
package main

import (
	"flag"
	"fmt"
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

var (
	// Cache, indexed by scope
	cache = map[string][]*cacheEntry{
		"": { // Global scope
			{"FALSE", "False", "uintptr", "0"},
			{"NULL", "Null", "uintptr", "0"},
			{"TRUE", "True", "uintptr", "1"},
		},
	}
	// WIN32_LEAN_AND_MEAN... and a bunch of other files
	headers []string = []string{
		"accctrl.h",
		"cderr.h",
		"dde.h",
		"ddeml.h",
		"dlgs.h",
		"joystickapi.h",
		"lzexpand.h",
		"mciapi.h",
		"mmeapi.h",
		"mmiscapi.h",
		"mmiscapi2.h",
		"mmsyscom.h",
		"mmsystem.h",
		"nb30.h",
		"playsoundapi.h",
		"poppack.h",
		"pshpack1.h",
		"rpc.h",
		"shellapi.h",
		"stdlib.h",
		"timeapi.h",
		"tlhelp32.h",
		"winapifamily.h",
		"wincrypt.h",
		"winefs.h",
		"winhttp.h",
		"wininet.h",
		"winnt.h",
		"winperf.h",
		"winscard.h",
		"winsock.h",
		"winuser.h",
	}
	lookup = map[string]*cacheEntry{}
	// Regular expressions
	reBitwisenot = regexp.MustCompile(`\~`)
	reCamel      = regexp.MustCompile(`[A-Z][a-z]+[A-Z][a-z]+`)
	reComment    = regexp.MustCompile(`\s*\/\*.*\*\/\s*`)
	reFixCast    = regexp.MustCompile(
		`\([A-Za-z_]+\)\s*(\(?\-?[0-9A-Fa-fXx]+|NULL)`,
	)
	reFixHex        = regexp.MustCompile(`(0[Xx][0-9A-Fa-f]+)[LlUu]+`)
	reFixLen        = regexp.MustCompile(`(len\(.+?\))`)
	reFixMSABI      = regexp.MustCompile(`__MSABI_LONG([^)]+)`)
	reFixNum        = regexp.MustCompile(`(\d+)[LlUu]+`)
	reFixSizeOf     = regexp.MustCompile(`sizeof\s*\(`)
	reSpaces        = regexp.MustCompile(`\s+`)
	reUselessOpRepl = []*regexp.Regexp{
		regexp.MustCompile(`\s*\^\s*0([^x])`), // ^ 0
		regexp.MustCompile(`\s*\|\s*0([^x])`), // | 0
	}
	reUselessOpRm = []*regexp.Regexp{
		regexp.MustCompile(`\s*\^\s*\(0\)`),         // ^ (0)
		regexp.MustCompile(`\s*\|\s*\(0\)`),         // | (0)
		regexp.MustCompile(`(0x0{8}\s*\|\s*)+`),     // 0x00000000 |
		regexp.MustCompile(`(\(0x0{8}\)\s*\|\s*)+`), // (0x00000000) |
		// ((0x0...0)) |
		regexp.MustCompile(`(\(\(0x0{8}\)\)\s*\|\s*)+`),
		regexp.MustCompile(`(\s*\|\s*0x0{8})+`),     // | 0x00000000
		regexp.MustCompile(`(\s*\|\s*\(0x0{8}\))+`), // | (0x00000000)
		// | ((0x0...0))
		regexp.MustCompile(`(\s*\|\s*\(\(0x0{8}\)\))+`),
	}
	// Skips
	skipLContains = map[string][]string{
		"": {
			"(",
			")",
			"DECLSPEC",
			"EXTERN_C",
		},
		"shellapi.h": {
			"DUMMY",
		},
		"stdlib.h": {
			"__",
			"errno",
		},
		"winnt.h": {
			"DUMMY",
			"XSTATE_MASK_ALLOWED",
		},
	}
	skipRContains = map[string][]string{
		"": {
			"__declspec",
			"__MINGW_NAME",
			"DECLSPEC",
			"HRESULT",
			"len(double)",
			"len(DWORD)",
			"len(ULONGLONG)",
			"WINAPI",
		},
		"mciapi.h": {
			"DRV_",
			"MM_",
		},
		"mmeapi.h": {
			"MM_",
		},
		"nb30.h": {
			"\\0",
		},
		"shellapi.h": {
			"FIELD_OFFSET",
		},
		"stdlib.h": {
			"__",
		},
		"wincrypt.h": {
			"\\0",
		},
		"wininet.h": {
			"INTERNET_STATUS_CALLBACK",
		},
		"winnt.h": {
			"DWORD64",
			"FIELD_OFFSET",
			"HANDLE",
			"inline",
			"MAKELANGID(",
			"MAKELCID(",
		},
		"winuser.h": {
			"len(LRESULT)",
			"MAKEINTATOM(",
		},
	}
	skipRStarts = map[string][]string{
		"": {
			":",
			"_",
			"extern",
			"void",
		},
		"accctrl.h": {
			"LocalFree",
		},
		"ddeml.h": {
			"CALLBACK",
		},
		"mmiscapi.h": {
			"mmioFOURCC",
		},
		"playsoundapi.h": {
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
)

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

	if reFixLen.MatchString(v) {
		v = reFixLen.ReplaceAllString(v, "uintptr($1)")
	}

	switch {
	case strings.HasPrefix(v, "\""):
		t = "string"
	case strings.HasPrefix(v, "L\""):
		v = v[1:]
		t = "string"
	case strings.HasPrefix(v, "TEXT("):
		v = strings.Replace(v[5:], ")", "", 1)
		t = "string"
	case strings.Contains(v, "L\""):
		v = "[]string{" + strings.ReplaceAll(v, " L\"", ", \"") + "}"
		t = "[]string"
	case strings.Contains(v, "-"):
		t = "int"
	case strings.Contains(v, ".") && strings.HasSuffix(v, "f"):
		v = strings.TrimSuffix(v, "f")
		t = "float64"
	case strings.HasPrefix(v, "{"):
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
			switch {
			case strings.HasPrefix(entry.Val, "\""):
				entry.Type = "string"
			case strings.HasPrefix(entry.Val, "[]string{"):
				entry.Type = "[]string"
			case strings.HasPrefix(entry.Val, "[]uintptr{"):
				entry.Type = "[]uintptr"
			case strings.Contains(entry.Val, "-"):
				entry.Type = "int"
			}
		}
	}
}

func format(str string) string {
	var out []string
	var tmp []string

	if str == "" {
		return str
	}

	for _, s := range strings.Split(str, " ") {
		// Replace _ with CamelCase
		switch {
		case strings.Contains(s, "_"):
			// Split on "_"
			tmp = strings.Split(strings.ToLower(s), "_")

			// Capitalize every part
			for i := range tmp {
				if tmp[i] == "" {
					continue
				}

				tmp[i] = strings.ToUpper(tmp[i][:1]) + tmp[i][1:]
			}

			// Join together for camelcase
			s = strings.Join(tmp, "")
		case reCamel.MatchString(s):
			// Do nothing
		default:
			s = strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		}

		// Fix some special cases
		s = strings.ReplaceAll(s, "Crlf", "CRLF")
		s = strings.ReplaceAll(s, "Ftp", "FTP")
		s = strings.ReplaceAll(s, "Hf32", "")
		s = strings.ReplaceAll(s, "Http", "HTTP")
		s = strings.ReplaceAll(s, "Lf32", "")
		s = strings.ReplaceAll(s, "LOGICAL", "Logical")
		s = strings.ReplaceAll(s, "LOGNAME", "Logname")
		s = strings.ReplaceAll(s, "Snapall", "SnapAll")
		s = strings.ReplaceAll(s, "Snapheaplist", "SnapHeapList")
		s = strings.ReplaceAll(s, "Snapmodule", "SnapModule")
		s = strings.ReplaceAll(s, "Snapprocess", "SnapProcess")
		s = strings.ReplaceAll(s, "Snapthread", "SnapThread")
		s = strings.ReplaceAll(s, "Th32cs", "")

		// Fix type-related
		s = strings.ReplaceAll(s, "BYTE", "Byte")
		s = strings.ReplaceAll(s, "CHAR", "Char")
		s = strings.ReplaceAll(s, "DWORD", "Dword")
		s = strings.ReplaceAll(s, "MAX", "Max")
		s = strings.ReplaceAll(s, "MIN", "Min")
		s = strings.ReplaceAll(s, "LONG", "Long")
		s = strings.ReplaceAll(s, "SHORT", "Short")
		s = strings.ReplaceAll(s, "SIZE", "Size")
		s = strings.ReplaceAll(s, "WORD", "Word")

		out = append(out, s)
	}

	return strings.Join(out, " ")
}

func genFile() (e error) {
	var entries []*cacheEntry
	var f *os.File
	var scopes []string

	// Build lookup table and replace C-style var names with Go-style
	// var names
	buildLookup()
	replaceVars()
	fixVarTypes()

	// Open file to write
	if f, e = genHeader(); e != nil {
		return e
	}
	defer func() {
		if e == nil {
			e = f.Close()
		}
	}()

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
		_, _ = fmt.Fprintf(
			f,
			"\n// %s contains constants from %s.h\n",
			format(scope),
			scope,
		)
		_, _ = fmt.Fprintf(f, "var %s = struct {\n", format(scope))

		// Define struct
		for _, entry := range entries {
			_, _ = fmt.Fprintf(f, "\t%s %s\n", entry.Go, entry.Type)
		}

		_, _ = fmt.Fprintln(f, "}{")

		// Create struct
		for _, entry := range entries {
			_, _ = fmt.Fprintf(f, "\t%s: %s,\n", entry.Go, entry.Val)
		}

		_, _ = fmt.Fprintln(f, "}")
	}

	return nil
}

func genHeader() (*os.File, error) {
	var e error
	var f *os.File
	var fn string = "generated.go"

	// Open file to write
	if f, e = os.Create(fn); e != nil {
		return nil, errors.Newf("failed to create %s: %w", fn, e)
	}

	// Create header
	_, _ = fmt.Fprintln(
		f,
		"// Code generated by tools/defines.go; DO NOT EDIT.",
	)
	_, _ = fmt.Fprintf(f, "package api\n\n")

	sort.Slice(
		cache[""], // Global scope
		func(i int, j int) bool {
			var l string = strings.ToLower(cache[""][i].Go)
			var r string = strings.ToLower(cache[""][j].Go)

			return l < r
		},
	)

	_, _ = fmt.Fprintln(f, "const (")

	for _, entry := range cache[""] {
		_, _ = fmt.Fprintf(
			f,
			"\t%s %s = %s\n",
			entry.Go,
			entry.Type,
			entry.Val,
		)
	}

	_, _ = fmt.Fprintln(f, ")")

	return f, nil
}

func ignoreType(fn string, line string, sep string) {
	for _, t := range strings.Split(line, sep) {
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
	// Find all the things to ignore/skip first (probably can remove
	// later after implementing struct parsing)
	for _, header := range headers {
		header = "/usr/x86_64-w64-mingw32/include/" + header

		if ok, e := pathname.DoesExist(header); e != nil {
			fmt.Println(e.Error())
			os.Exit(1)
		} else if !ok {
			fmt.Printf("file %s not found\n", header)
			os.Exit(1)
		}

		if e := processFileSkips(header); e != nil {
			panic(e)
		}
	}

	// Then process each file for #defines and enums
	for _, header := range headers {
		header = "/usr/x86_64-w64-mingw32/include/" + header

		if e := processFileDefines(header); e != nil {
			panic(e)
		}

		// TODO process structs too

		if e := processFileTypedefs(header); e != nil {
			panic(e)
		}
	}

	if e := genFile(); e != nil {
		panic(e)
	}
}

func processDefine(fn string, line string) {
	if !strings.HasPrefix(line, "#define") {
		return
	}

	// Remove define
	line = strings.Replace(line, "#define ", "", 1)

	// Remove comments
	if reComment.MatchString(line) {
		line = reComment.ReplaceAllString(line, "")
	}

	// Fix some Windows-isms
	if reBitwisenot.MatchString(line) {
		line = reBitwisenot.ReplaceAllString(line, "0xffffffff^")
	}

	if reFixSizeOf.MatchString(line) {
		line = reFixSizeOf.ReplaceAllString(line, "len(")
	}

	for _, r := range []*regexp.Regexp{
		reFixHex,
		reFixCast,
		reFixMSABI,
		reFixNum,
	} {
		if r.MatchString(line) {
			line = r.ReplaceAllString(line, "$1")
		}
	}

	// Cut and attempt to skip things we don't want
	if k, v, _ := strings.Cut(line, " "); !skip(fn, k, v) {
		cacheVar(fn, k, v)
	}
}

func processFileDefines(path string) error {
	var b []byte
	var e error
	var fn string = filepath.Base(path)
	var tmp string

	if b, e = os.ReadFile(filepath.Clean(path)); e != nil {
		return errors.Newf("failed to open %s: %w", path, e)
	}

	for _, line := range strings.Split(string(b), "\n") {
		line = strings.ReplaceAll(
			line,
			"| INTERNET_FLAG_BGUPDATE",
			"",
		)
		line = strings.TrimSpace(reSpaces.ReplaceAllString(line, " "))

		if strings.HasSuffix(line, "\\") {
			tmp += strings.TrimSuffix(line, "\\")
			continue
		}

		processDefine(fn, strings.TrimSpace(tmp+line))

		// Reset
		tmp = ""
	}

	return nil
}

func processFileSkips(path string) error {
	var b []byte
	var e error
	var fn string = filepath.Base(path)
	var inStructOrTypedef bool

	if b, e = os.ReadFile(filepath.Clean(path)); e != nil {
		return errors.Newf("failed to read %s: %w", path, e)
	}

	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(reSpaces.ReplaceAllString(line, " "))

		switch {
		case strings.HasPrefix(line, "#define"):
			line = strings.TrimPrefix(line, "#define")
			line = strings.TrimSpace(line)

			//nolint:mnd // Skip anything that's too long
			if !strings.Contains(line, " ") && (len(line) > 8) {
				skipRContains[fn] = append(skipRContains[fn], line)
			}
		case strings.HasPrefix(line, "}"):
			if inStructOrTypedef {
				// Extract type
				line = strings.TrimPrefix(line, "} ")
				line = strings.TrimSuffix(line, ";")

				ignoreType(fn, line, ",")
			}

			inStructOrTypedef = false
		case strings.HasPrefix(line, "typedef const"):
			// Extract type
			line = strings.TrimPrefix(line, "typedef const")
			line = strings.TrimPrefix(line, "struct")
			line = strings.TrimSuffix(line, ";")
			line = strings.TrimSpace(line)

			if line != "" {
				ignoreType(fn, line, " ")
			}
		case strings.HasPrefix(line, "typedef enum"),
			strings.HasPrefix(line, "typedef struct"):
			inStructOrTypedef = true

			// Extract type
			line = strings.TrimPrefix(line, "typedef enum")
			line = strings.TrimPrefix(line, "typedef struct")
			line = strings.TrimSuffix(line, "{")
			line = strings.TrimSpace(line)

			if line != "" {
				ignoreType(fn, line, ",")
			}
		}
	}

	return nil
}

func processFileTypedefs(path string) error {
	var b []byte
	var e error
	var fn string = filepath.Base(path)
	var tmp string
	var inTypedef bool

	if b, e = os.ReadFile(filepath.Clean(path)); e != nil {
		return errors.Newf("failed to open %s: %w", path, e)
	}

	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(reSpaces.ReplaceAllString(line, " "))

		// Ignore non-defines unless typedef enum
		switch {
		case line == "{":
			continue
		case strings.HasPrefix(line, "#"):
			// #ifdef and #endif
			continue
		case strings.HasPrefix(line, "{"):
			if inTypedef {
				tmp += line[1:]
			}

			continue
		case strings.HasPrefix(line, "}"):
			if !inTypedef {
				continue
			}

			inTypedef = false
		case strings.HasPrefix(line, "typedef enum"):
			inTypedef = true
			continue
		default:
			if inTypedef {
				tmp += line
			}

			continue
		}

		// Remove comments
		if reComment.MatchString(tmp) {
			tmp = reComment.ReplaceAllString(tmp, "")
		}

		processTypedef(fn, tmp)

		// Reset
		tmp = ""
	}

	return nil
}

func processTypedef(fn string, line string) {
	var lhs string
	var next string = "0"
	var rhs string
	var vals map[string]string = map[string]string{}

	for _, d := range strings.Split(line, ",") {
		d = strings.TrimSpace(d)

		if d == "" {
			continue
		}

		if strings.Contains(d, "=") {
			lhs, rhs, _ = strings.Cut(d, "=")

			lhs = format(strings.TrimSpace(lhs))
			rhs = format(strings.TrimSpace(rhs))

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
			next = strconv.Itoa(i + 1)
		}

		// Kkip things we don't want
		if !skip(fn, lhs, rhs) {
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
			for range 8 {
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

			for _, reUselessOp := range reUselessOpRepl {
				if reUselessOp.MatchString(entry.Val) {
					entry.Val = reUselessOp.ReplaceAllString(
						entry.Val,
						"$1",
					)
				}
			}

			for _, reUselessOp := range reUselessOpRm {
				if reUselessOp.MatchString(entry.Val) {
					entry.Val = reUselessOp.ReplaceAllString(
						entry.Val,
						"",
					)
				}
			}
		}
	}
}

func skip(fn string, k string, v string) bool {
	if v == "" {
		//nolint:mnd // Skip anything that's too long
		if len(k) > 8 {
			skipRContains[fn] = append(skipRContains[fn], k)
		}

		return true
	}

	if strings.ToLower(v) == strings.ToLower(k)+"w" {
		skipRContains[fn] = append(skipRContains[fn], k)
		return true
	}

	for _, scope := range []string{"", fn} {
		for _, s := range skipLContains[scope] {
			if strings.Contains(k, s) {
				skipRContains[fn] = append(skipRContains[fn], k)
				return true
			}
		}

		for _, s := range skipRContains[scope] {
			if strings.Contains(v, s) {
				skipRContains[fn] = append(skipRContains[fn], k)
				return true
			}
		}

		for _, s := range skipRStarts[scope] {
			if strings.HasPrefix(v, s) {
				skipRContains[fn] = append(skipRContains[fn], k)
				return true
			}
		}
	}

	return false
}
