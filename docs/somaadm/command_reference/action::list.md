# somaadm action list

This command shows the actions for a specific section.

# SYNOPSIS

```
somaadm action list in ${section} [of ${category}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
section | string | Name of the section | | no
category | string | Name of the category | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | action | list | yes | no

# EXAMPLES

```
./somaadm action list in environment of global
```
