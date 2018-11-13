# DESCRIPTION

This command is used to fetch updates for asynchronous server side
processing jobs.

The `soma` client stores the JobID of every 204/Accepted asynchronous
job in a local cache database. For every Job not in status `processed`,
an update is fetched and displayed. For jobs that have been completed,
the local cache data is updated as well.

# SYNOPSIS

```
soma job update
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | self | | no | yes
self | job | search | yes | no

# EXAMPLES

```
soma job update
```
