# olin

> Intelligent networks are able to operate from a single language with translation interfaces that enable global intercourse. This means language is no longer a barrier to communication. Intelligent networks will introduce a meta-language that translates both real-time written and spoken applications. It will revolutionize the genetic mind's global construct, and facilitate the digitalization of your global economy.

olin is an environment to run and operate functions as a service projects using
event sourcing and webassembly under the hood. Your handler code shouldn't need
to care that there is an event queue involved. Your handler should just do what
your handler needs to do.

## Background

Very frequently, I end up needing to write applications that basically end up
waiting forever to make sure things get put in the right place and then the 
right code runs as a response. I then have to make sure these things get put
in the right places and then that the right versions of things are running for
each of the relevant services. This doesn't scale very well, not to mention is 
hard to secure. This leads to a lot of duplicate infrastructure over time and 
as things grow. Not to mention adding in tracing, metrics and log aggreation.

I would like to change this.

I would like to make a perscriptive environment kinda like [Google Cloud Functions][gcf]
or [AWS Lambda][lambda] backed by a durable message queue and with handlers
compiled to webassembly to ensure forward compatibility. As such, the ABI 
involved will be versioned, documented and tested. Multiple ABI's will eventually
need to be maintained in parallel, so it might be good to get used to that early
on.

I expect this project to last decades. I want the binary modules I upload today
to be still working in 5 years, assuming its dependencies outside of the module 
still work. 

## Dagger

> The dagger of light that renders your self-importance a decisive death

Dagger is the first ABI that will be used for interfacing with the outside world.
This will be mostly for an initial spike out of the basic ideas to see what it's 
like while the rest of the plan is being stabilized and implemented.
The core idea is that everything is a file, to the point that the file descriptor
and file handle array are the only real bits of persistent state for the process.
HTTP sessions, logging writers, TCP sockets, operating system files, cryptographic 
random readers, everything is done via filesystem system calls.

Consider this the first draft of Dagger, everything here is subject to change.
This is going to be the experimental phase.

### Base Environment

When a dagger process is opened, the following files are open:

- 0: standard input: the semantic "input" of the program.
- 1: standard output: the standard output of the program.
- 2: standard error: error output for the program.

### File Handlers

In the open call (defined later), a file URL is specified instead of a file name.
This allows for Dagger to natively offer programs using it quick access to common
services like HTTP, logging or pretty much anything else.

I'm playing with the following handlers currently:

- http and https (Write request as http/1.1 request and sync(), Read response as http/1.1 response and close())
- rand - cryptographically secure random data
- time - unix timestamp in a little-endian encoded int64 on every read() - `time://utc`

I'd like to add the following handlers in the future:

- file - filesystem files on the host OS (dangerous!)
- tcp - TCP connections
- tcp+tls - TCP connections with TLS
- meta - metadata about the runtime or the event

### Handler Function

Each Dagger module can only handle one data type. This is intentional. This 
forces users to make a separate handler for each type of data they want to 
handle. The handler function reads its input from standard input and then 
returns `0` if whatever it needs to do "worked" (for some definition of success).

TODO(Xe): Find better way to say:

This function could be exposed in clang by doing:

```c
__attribute__ ((visibility ("default")))
int handle() {
  // read all of standard input to memory and handle it
  return 0; // success
}
```

### System Calls

A [system call][syscall] is how computer programs interface with the outside
world. When a dagger program makes a system call, the amount of time the program
spends waiting for that system call is collected and recorded based on what
underlying resource took care of the call. This means, in theory, users of olin
could alert on HTTP requests from one service to another taking longer amounts
of time very trivially.

Dagger uses the following system calls:

- open
- close
- read
- write
- sync

#### open

```c
extern int open(const char *furl, int flags);
```

This opens a file with the given handler and flags (provided the handler supports
them) and returns its file descriptor. This descriptor is used in all later calls.

#### close

```c
extern int close(int fd);
```

Close closes a file and returns if it failed or not. If this call returns nonzero,
you don't know what state the world is in. Panic.

#### read

```c
extern int read(int fd, void *buf, int nbyte);
```

Read attempts to read up to count bytes from file descriptor fd into the buffer 
starting at buf.

#### write

```c
extern int write(int fd, void *buf, int nbyte);
```

Write writes up to count bytes from the buffer starting at buf to the file 
referred to by the file descriptor fd.

#### sync

```c
extern int sync(int fd);
```

This is for some backends to forcibly make async operations into sync operations.

## Project Meta

To follow the project, check it on GitHub [here][olin]. To talk about it on Slack,
join the [Go community Slack][goslack] and join `#olin`. 

[gcf]: https://cloud.google.com/functions/
[lambda]: https://aws.amazon.com/lambda/
[syscall]: https://en.wikipedia.org/wiki/System_call
[olin]: https://github.com/Xe/olin
[goslack]: https://invite.slack.golangbridge.org

