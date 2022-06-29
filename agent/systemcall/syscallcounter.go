package systemcall

// type SyscallCounter []int

// const maxSyscalls = 303

// func (s SyscallCounter) Init() SyscallCounter {
// 	s = make(SyscallCounter, maxSyscalls)
// 	return s
// }

// func (s SyscallCounter) Inc(syscallID uint64) error {
// 	if syscallID > maxSyscalls {
// 		return fmt.Errorf("invalid syscall ID (%x)", syscallID)
// 	}

// 	s[syscallID]++
// 	return nil
// }

// func (s SyscallCounter) Print() {
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', tabwriter.AlignRight|tabwriter.Debug)
// 	for k, v := range s {
// 		if v > 0 {
// 			name, _ := sec.ScmpSyscall(k).GetName()
// 			fmt.Fprintf(w, "%d\t%s\n", v, name)
// 		}
// 	}
// 	w.Flush()
// }

// func (s SyscallCounter) getName(syscallID uint64) string {
// 	name, _ := sec.ScmpSyscall(syscallID).GetName()
// 	return name
// }
