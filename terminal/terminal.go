package terminal

import (
    "os"
    "syscall"
    "unsafe"
)
    
type Terminal struct {
    fd uintptr
    original Termios
    modified *Termios
    NCols int
    NRows int
}

func New() (*Terminal, error) {
    t := &Terminal{}
    
    t.fd = os.Stdout.Fd()
    termios, err := getTermios(t.fd)
    if err != nil {
        return nil, err
    }
    
    t.original = *termios
    t.modified = termios
    
    t.enableRawMode()
    err = t.getWindowSize()
    if err != nil {
        return nil, err
    }
    
    err = setTermios(t.fd, t.modified)
    if err != nil {
        return nil, err
    }
    
    return t, nil
}

func (t *Terminal) Restore() error {
    return setTermios(t.fd, &t.original)
}

func (t *Terminal) getWindowSize() error {
    ws := struct { 
        row uint16
        col uint16
    }{}
    
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, t.fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
    if errno != 0 {
        return errno
    }
    
    t.NCols = int(ws.col)
    t.NRows = int(ws.row)
    
    return nil
}

func (t *Terminal) enableRawMode() {
    t.modified.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
    t.modified.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
    t.modified.Cflag |= syscall.CS8
    t.modified.Oflag &^= syscall.OPOST
    t.modified.Cc[syscall.VMIN+1] = 0
    t.modified.Cc[syscall.VTIME+1] = 1
}
