package proc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)
/*
 the pid info are maintained and updated in three files in /proc/pid :
 stat ,status and 
 */





/*
 * Fields from https://www.kernel.org/doc/Documentation/filesystems/proc.txt
 */

type Stat_t struct {
	// fields from /proc/<pid>/stat -- order here is critical
	// Remember: lower case fields are not exposed (ie. ignored)

	Pid           uint64 // process id
	Tcomm         string // filename of the executable '(executable)'
	State         string // state (R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced/stopped)
	Ppid          int64  // process id of the parent process
	Pgrp          int64  // pgrp of the process
	Sid           int64  // session id
	Tty_nr        int64  // tty the process uses
	Tty_pgrp      int64  // pgrp of the tty
	Flags         uint64 // task flags
	Min_flt       uint64 // number of minor faults
	Cmin_flt      uint64 // number of minor faults with child's
	Maj_flt       uint64 // number of major faults
	Cmaj_flt      uint64 // number of major faults with child's
	Utime         uint64 // user mode jiffies
	Stime         uint64 // kernel mode jiffies
	Cutime        uint64 // user mode jiffies with child's
	Cstime        uint64 // kernel mode jiffies with child's
	Priority      int32  // priority level
	Nice          int32  // nice level
	Num_threads   uint32 // number of threads
	it_real_value uint32 // (obsolete, always 0) -- IGNORED
	Start_time    uint64 // time the process started after system boot (* 1000000 for no decimals)
	Vsize         uint64 // virtual memory size
	Rss           uint64 // resident set memory size
	Rsslim        uint64 // current limit in bytes on the rss
	Start_code    uint64 // address above which program text can run
	End_code      uint64 // address below which program text can run
	Start_stack   uint64 // address of the start of the main process stack
	Esp           uint64 // current value of ESP
	Eip           uint64 // current value of EIP
	Pending       string // bitmap of pending signals
	Blocked       string // bitmap of blocked signals
	Sigign        string // bitmap of ignored signals
	Sigcatch      string // bitmap of caught signals
	Wchan         uint64 // address where process went to sleep
	placeholder1  uint64 // (place holder) -- IGNORED
	placeholder2  uint64 // (place holder) -- IGNORED
	Exit_signal   uint64 // signal to send to parent thread on exit
	Task_cpu      uint64 // which CPU the task is scheduled on
	Rt_priority   uint64 // realtime priority
	Policy        uint64 // scheduling policy (man sched_setscheduler)
	Blkio_ticks   uint64 // time spent waiting for block IO
	Gtime         uint64 // guest time of the task in jiffies
	Cgtime        uint64 // guest time of the task children in jiffies
	Start_data    uint64 // address above which program data+bss is placed
	End_data      uint64 // address below which program data+bss is placed
	Start_brk     uint64 // address above which program heap can be expanded with brk()
	Arg_start     uint64 // address above which program command line is placed
	Arg_end       uint64 // address below which program command line is placed
	Env_start     uint64 // address above which program environment is placed
	Env_end       uint64 // address below which program environment is placed
	Exit_code     uint64 // the thread's exit_code in the form reported by the waitpid system call
}

type Statm_t struct {
	// fields from /proc/<pid>/statm -- order here is critical

	Size     uint64 // total program size (pages) (same as VmSize in status)
	Resident uint64 // size of memory portions (pages) (same as VmRSS in status)
	Shared   uint64 // number of pages that are shared (i.e. backed by a file)
	Trs      uint64 // number of pages that are 'code' (not including libs; broken, includes data segment)
	Lrs      uint64 // number of pages of library (always 0 on 2.6)
	Drs      uint64 // number of pages of data/stack (including libs; broken, includes library text)
	Dt       uint64 // number of dirty pages (always 0 on 2.6)
}

type Ids struct {
	Real      uint64
	Effective uint64
	Saved     uint64
	FS        uint64
}

type SigQVal struct {
	Num uint64
	Max uint64
}

type Status_t struct {
	// fields from /proc/<pid>/status

	Name                       string   // filename of the executable
	State                      string   // state (R=running, S=sleeping, D=sleeping in an uninterruptible wait, Z=zombie, T=(traced or stopped))
	Tgid                       uint64   // thread group ID
	Ngid                       uint64   // numa group ID
	Pid                        uint64   // process id
	PPid                       uint64   // process id of the parent process
	TracerPid                  uint64   // PID of process tracing this process (0 if not)
	Uid                        Ids      // Real, effective, saved set, and  file system UIDs
	Gid                        Ids      // Real, effective, saved set, and  file system GIDs
	FDSize                     uint64   // number of file descriptor slots currently allocated
	Groups                     []uint64 // supplementary group list
	VmPeak                     uint64   // peak virtual memory size
	VmSize                     uint64   // total program size
	VmLck                      uint64   // locked memory size
	VmPin                      uint64   // locked memory size
	VmHWM                      uint64   // peak resident set size ("high water mark")
	VmRSS                      uint64   // size of memory portions
	VmData                     uint64   // size of data, stack, and text segments
	VmStk                      uint64   // size of data, stack, and text segments
	VmExe                      uint64   // size of text segment
	VmLib                      uint64   // size of shared library code
	VmPTE                      uint64   // size of page table entries
	VmSwap                     uint64   // size of swap usage (the number of referred swapents)
	Threads                    uint64   // number of threads
	SigQ                       SigQVal  // number of signals queued (Num) / limit (Max)
	SigPnd                     string   // bitmap of pending signals for the thread
	ShdPnd                     string   // bitmap of shared pending signals for the process
	SigBlk                     string   // bitmap of blocked signals
	SigIgn                     string   // bitmap of ignored signals
	SigCgt                     string   // bitmap of caught signals
	CapInh                     string   // bitmap of inheritable capabilities
	CapPrm                     string   // bitmap of permitted capabilities
	CapEff                     string   // bitmap of effective capabilities
	CapBnd                     string   // bitmap of capabilities bounding set
	Seccomp                    uint64   // seccomp mode, like prctl(PR_GET_SECCOMP, ...)
	Cpus_allowed               string   // mask of CPUs on which this process may run "mask format"
	Cpus_allowed_list          string   // Same as previous, but in "list format"
	Mems_allowed               string   // mask of memory nodes allowed to this process "mask format"
	Mems_allowed_list          string   // Same as previous, but in "list format"
	Voluntary_ctxt_switches    uint64   // number of voluntary context switches
	Nonvoluntary_ctxt_switches uint64   // number of non voluntary context switches
}

type Proc struct {
	Stat   Stat_t
	Statm  Statm_t
	Status Status_t

	// Environ and Cmdline are from /proc/<pid>/{environ,cmdline}
	Cmdline []string
	Environ []string
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(cfg *procConfig, pid uint64, filename string) ([]string, error) {
	var lines []string
	var scanner *bufio.Scanner

	if contents, ok := cfg.contents[filename]; ok {
		scanner = bufio.NewScanner(strings.NewReader(contents))
	} else {
		fn := fmt.Sprintf("%s/%d/%s", cfg.basepath, pid, filename)
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, wrapError(err)
		}
		// for generating test cases, having the input is required
		cfg.contents[filename] = string(data)
		scanner = bufio.NewScanner(strings.NewReader(cfg.contents[filename]))
	}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, wrapError(scanner.Err())
}

func readStat(cfg *procConfig, pid uint64, proc *Proc) error {
	var stat Stat_t

	lines, err := readLines(cfg, pid, "stat")
	if err != nil {
		return wrapError(err)
	}
	if len(lines) != 1 {
		return newError("readStat(): expected 1 line, got %d", len(lines))
	}

	s := reflect.ValueOf(&stat).Elem()
	typeOfS := s.Type()

	// field[1] is special in that it's '(<command>)' and we need to
	// watch out for things like ')' in the <command>. We do what
	// procps does and take everything between the first '(' and
	// last ')'
	cmd_end := strings.LastIndex(lines[0], ")")
	cmd_start := strings.Index(lines[0], "(") + 1

	// This mess:
	//
	//   * splits all the fields *after* the command
	//   * prepends the first field (Pid) and command (w/ '()' removed)
	//
	fields := strings.Split(lines[0][cmd_end+2:], " ")
	fields = append([]string{lines[0][0 : cmd_start-1],
		lines[0][cmd_start:cmd_end]}, fields...)

	// loop through the struct's fields and add them to 'stat'
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !f.CanSet() {
			continue
		}

		// Older kernels are missing newer fields
		if i >= len(fields) {
			continue
		}

		//fmt.Printf("%d: %s %s = '%s'\n", i,
		//typeOfS.Field(i).Name, f.Type(), fields[i])

		switch f.Type().String() {
		case "int32":
			u, err := strconv.ParseInt(strings.TrimSpace(fields[i]), 10, 32)
			if err != nil {
				return wrapError(err)
			}
			s.Field(i).SetInt(u)
		case "int64":
			u, err := strconv.ParseInt(strings.TrimSpace(fields[i]), 10, 64)
			if err != nil {
				return wrapError(err)
			}
			s.Field(i).SetInt(u)
		case "uint32":
			u, err := strconv.ParseUint(strings.TrimSpace(fields[i]), 10, 32)
			if err != nil {
				return wrapError(err)
			}
			s.Field(i).SetUint(u)
		case "uint64":
			u, err := strconv.ParseUint(strings.TrimSpace(fields[i]), 10, 64)
			if err != nil {
				return wrapError(err)
			}
			s.Field(i).SetUint(u)
		case "string":
			s.Field(i).SetString(fields[i])

		default:
			return newError("readStat(): unhandled type '%s' for '%s'",
				f.Type().String(), typeOfS.Field(i).Name)
		}
	}

	proc.Stat = stat
	return nil
}

func readStatm(cfg *procConfig, pid uint64, proc *Proc) error {
	var statm Statm_t

	lines, err := readLines(cfg, pid, "statm")
	if err != nil {
		return wrapError(err)
	}
	if len(lines) != 1 {
		return newError("readStatm(): expected 1 line, got %d", len(lines))
	}

	s := reflect.ValueOf(&statm).Elem()
	typeOfS := s.Type()

	fields := strings.Split(lines[0], " ")

	// loop through the struct's fields and add them to 'statm'
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !f.CanSet() {
			continue
		}

		// fmt.Printf("%d: %s %s = '%s'\n", i,
		//		typeOfS.Field(i).Name, f.Type(), fields[i])

		switch f.Type().String() {
		case "uint64":
			u, err := strconv.ParseUint(strings.TrimSpace(fields[i]), 10, 64)
			if err != nil {
				return wrapError(err)
			}
			s.Field(i).SetUint(u)
		default:
			return newError("readStatm(): unhandled type '%s' for '%s'",
				f.Type().String(), typeOfS.Field(i).Name)
		}
	}

	proc.Statm = statm
	return nil
}

func readStatus(cfg *procConfig, pid uint64, proc *Proc) error {
	var err error
	var statusMap = make(map[string]reflect.Value)
	var status Status_t

	lines, err := readLines(cfg, pid, "status")
	if err != nil {
		return wrapError(err)
	}

	s := reflect.ValueOf(&status).Elem()
	typeOfS := s.Type()
	for i := 0; i < s.NumField(); i++ {
		statusMap[typeOfS.Field(i).Name] = s.Field(i)
	}

	for line := range lines {
		fields := strings.SplitN(lines[line], ":\t", 2)
		name := fields[0]
		value := strings.TrimSpace(fields[1])

		f := statusMap[name]
		if !f.IsValid() {
			// not valid, see if capitalizing first character fixes
			a := []rune(name)
			a[0] = unicode.ToUpper(a[0])
			name = string(a)

			f = statusMap[name]
			if !f.IsValid() {
				// still not valid, can't move forward, though we'll skip fields we know
				// exist on older kernels that we don't care about any more.
				switch fields[0] {
				case "SleepAVG":
					continue
				}
				return newError("readStatus(): '%s' is unhandled", fields[0])
			}
		}

		switch f.Type().String() {
		case "procreader.Ids":
			if name == "Uid" {
				cnt, err := fmt.Sscanf(value, "%d\t%d\t%d\t%d",
					&status.Uid.Real, &status.Uid.Effective,
					&status.Uid.Saved, &status.Uid.FS,
				)
				if err != nil {
					return wrapError(err)
				}
				if cnt != 4 {
					return newError("readStatus(Uid): expected 4 fields, got %d: '%s'", cnt, value)
				}
			} else if name == "Gid" {
				cnt, err := fmt.Sscanf(value, "%d\t%d\t%d\t%d",
					&status.Gid.Real, &status.Gid.Effective,
					&status.Gid.Saved, &status.Gid.FS,
				)
				if err != nil {
					return wrapError(err)
				}
				if cnt != 4 {
					return newError("readStatus(Gid): expected 4 fields, got %d: '%s'", cnt, value)
				}
			} else {
				return newError("readStatus: Internal Error: %s not supported for type Ids", name)
			}
		case "procreader.SigQVal":
			cnt, err := fmt.Sscanf(value, "%d/%d", &status.SigQ.Num, &status.SigQ.Max)
			if err != nil {
				return wrapError(err)
			}
			if cnt != 2 {
				return newError("readStatus[%s]: expected 2 fields, got %d: '%s'", name, cnt, value)
			}
		case "[]uint64":
			if name != "Groups" {
				return newError("readStatus: Internal Error: %s not supported for type []uint64", name)
			}
			groups := strings.Split(value, " ")
			for g := range groups {
				if len(groups[g]) == 0 {
					continue
				}
				val, err := strconv.ParseUint(groups[g], 10, 64)
				if err != nil {
					return wrapError(err)
				}
				status.Groups = append(status.Groups, val)
			}
		case "string":
			f.SetString(value)
		case "uint64":
			if strings.HasSuffix(value, " kB") {
				value = value[:len(value)-3]
			}
			u, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return wrapError(err)
			}
			f.SetUint(u)
		default:
			return newError("readStatus(): unhandled type '%s' for '%s'", f.Type().String(), name)
		}
	}

	proc.Status = status

	return nil
}

func readNullSeparated(cfg *procConfig, pid uint64, filename string) ([]string, error) {
	var r *bufio.Reader
	var strs []string

	if contents, ok := cfg.contents[filename]; ok {
		r = bufio.NewReader(strings.NewReader(contents))
	} else {
		fn := fmt.Sprintf("%s/%d/%s", cfg.basepath, pid, filename)
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, wrapError(err)
		}
		// for generating test cases, having the input is required
		cfg.contents[filename] = string(data)
		r = bufio.NewReader(strings.NewReader(cfg.contents[filename]))
	}

	for {
		str, err := r.ReadString(0)
		if err == io.EOF {
			// do something here
			break
		} else if err != nil {
			return nil, wrapError(err)
		}
		strs = append(strs, strings.TrimRight(str, "\000"))
	}

	return strs, nil
}

func readCmdline(cfg *procConfig, pid uint64, proc *Proc) error {
	var err error
	proc.Cmdline, err = readNullSeparated(cfg, pid, "cmdline")
	return wrapError(err)
}

func readEnviron(cfg *procConfig, pid uint64, proc *Proc) error {
	var err error
	proc.Environ, err = readNullSeparated(cfg, pid, "environ")
	return wrapError(err)
}

//
// This function dispatches the reading of the various /proc files but allows
// (via cfg) the replacement of the "readers" that actually read the files. This
// is mostly done to make this module more testable.
//
func readProc(cfg *procConfig, pid uint64) (Proc, error) {
	var err error
	var proc Proc

	err = readStat(cfg, pid, &proc)
	if err != nil {
		return proc, wrapError(err)
	}
	err = readStatm(cfg, pid, &proc)
	if err != nil {
		return proc, wrapError(err)
	}
	err = readStatus(cfg, pid, &proc)
	if err != nil {
		return proc, wrapError(err)
	}
	err = readCmdline(cfg, pid, &proc)
	if err != nil {
		return proc, wrapError(err)
	}
	err = readEnviron(cfg, pid, &proc)
	if err != nil {
		return proc, wrapError(err)
	}

	return proc, nil
}

// ReadProc
// This function reads /proc/<pid>/* files and returns a Proc object
// which contains the information for the specified process.
//
func ReadProc(pid uint64) (Proc, error) {
	proc, _, err := ReadProcData(pid)
	return proc, err
}

func ReadProcData(pid uint64) (Proc, map[string]string, error) {
	var cfg procConfig

	cfg.basepath = "/proc"
	cfg.contents = make(map[string]string)

	proc, err := readProc(&cfg, pid)

	return proc, cfg.contents, err
}
