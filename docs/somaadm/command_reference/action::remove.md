# somaadm action remove

This command is used to delete an action from a section.

# SYNOPSIS

```
somaadm action remove ${action} from ${section} [in ${category}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
action | string | Name of the action | | no
section | string | Name of the section | | no
category | string | Name of the category | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | action | remove | yes | no

# EXAMPLES

```
./somaadm actions remove add from environment
./somaadm actions remove remove from environment
./somaadm actions remove list from environment
./somaadm actions remove show from environment in global
```
