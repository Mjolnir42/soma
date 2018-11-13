# DESCRIPTION

This command is used to request a Job from the SOMA server.

# SYNOPSIS

```
soma job show ${jobID}
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
self | job | show | yes | no

# EXAMPLES

```
soma job show 34e9ca9c-6a6b-400f-a400-000000000000
```
