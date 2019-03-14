# DESCRIPTION

This command is used to revoke a granted permission.

# SYNOPSIS

```
soma right grant ${category}::${permission} from user|admin|team|tool ${name} [on repository|bucket|monitoring|team ${object}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
category | string | Name of the category | | no
permission | string | Name of the permission | | no
name | string | Name of the subject | | no
object | string | Name of the object | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | right | revoke | yes | no

# EXAMPLES

```
soma right revoke global::browse from user jd
soma right revoke monitoring::worker from user jd on monitoring ExampleMonitoring
```
