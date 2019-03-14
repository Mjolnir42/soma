# DESCRIPTION

This command is used to grant a permission.

# SYNOPSIS

```
soma right grant ${category}::${permission} to user|admin|team|tool ${name} [on repository|bucket|monitoring|team ${object}]
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
permission | right | grant | yes | no

# EXAMPLES

```
soma right grant global::browse to user jd
soma right grant monitoring::worker to user jd on monitoring ExampleMonitoring
```
