# psigo

A simple Go CLI tool for decoding Linux process signal masks.

## Description

psigo parses `/proc/[pid]/status` to extract and decode signal information for Linux processes. It can decode signal bitmasks and display signal names for better understanding of process signal states.

## Installation

```bash
go install github.com/realfatcat/psigo@latest
```

## Usage

### Decode signals for a specific process

```bash
psigo -p <PID>
```

Example:
```bash
psigo -p 1234
```

Output:
```
SigPnd: 1) SIGHUP 
SigBlk: 
SigIgn: 2) SIGINT 
SigCgt: 1) SIGHUP 2) SIGINT 
```

### Decode a hexadecimal signal mask

```bash
psigo -d <hex_value>
```

Example:
```bash
psigo -d 0x3
```

Output:
```
1) SIGHUP 2) SIGINT 
```

## Command Line Options

- `-p <PID>`: Decode signals for the given process ID
- `-d <hex>`: Decode the provided hexadecimal signal mask value

## Requirements

- Linux operating system
- Go 1.25.2 or later
- Appropriate permissions to read `/proc/[pid]/status`

## Signal Types

The tool displays four types of signal masks:

- **SigPnd**: Pending signals
- **SigBlk**: Blocked signals  
- **SigIgn**: Ignored signals
- **SigCgt**: Caught signals

## License

This project is open source. Please check the license file for details.
