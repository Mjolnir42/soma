# DESCRIPTION

This command is used to block the client on the completion
of a specific asynchronous job.

Information about completed jobs is held by the server for
up to 2 hours after job completion, in which case the client
will unblock immediately.

Max blocking time is 5 minutes after which clients always get
unblocked.

# SYNOPSIS

```
soma job wait ${jobID}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
jobID | string | UUID of the job | | no


# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | self | | no | yes
self | job | wait | yes | no

# EXAMPLES

```
soma job wait 34e9ca9c-6a6b-400f-a400-000000000000
```
