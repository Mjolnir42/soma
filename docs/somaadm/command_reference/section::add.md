# somaadm section add

This command adds a new permission section to the system. Sections
are groups of actions and belong to one permission category. They
are used by the middleware workers to query if a user is allowed
to perform an action.

In which section a worker queries is hardcoded. Sections added beyond
this set of well known sections will be unused. By not creating a
section that the server uses it becomes impossible to allow a user
to perform any actions from that section.

# SYNOPSIS

```
soma section add ${section} to ${category}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
section | string | Name of the section | | no
category | string | Name of the category of the section | | no 

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | section | add | yes | no

# EXAMPLES

```
soma section add environment to global
```
