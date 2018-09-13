# DESCRIPTION

This command show details about a specific action.

# SYNOPSIS

```
soma action show ${action} from ${section} [in ${category}]
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
permission | action | show | yes | no

# EXAMPLES

```
soma action show add from environment in global
soma action show remove from environment
soma action show list from environment
soma action show show from environment in global
```
