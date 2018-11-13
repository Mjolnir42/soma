# DESCRIPTION

The list commands present different views on the list of known jobs.

Invoked as `list outstanding`, the client dumps all jobs from the local
cache database not in status `processed`. These are the jobs that will
be updated by issuing a `soma job update` command.

Invoked as `list local`, the client dumps all jobs from the local cache
database regardless of processing status.

Invoked as `list remote`, the client fetches all jobs from the server,
regardless of processing status.

# SYNOPSIS

```
soma job list outstanding
soma job list local
soma job list remote
```

# ARGUMENT TYPES

These commands take no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | self | | no | yes
self | job | list | yes | no

The commands `list outstanding` and `list local` require no permissions
since they only work with the local client cache.

# EXAMPLES

```
soma job list outstanding
soma job list local
soma job list remote
```
