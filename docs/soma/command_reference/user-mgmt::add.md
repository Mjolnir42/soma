# DESCRIPTION

This command is used to add users to SOMA. In addition to various other
attributes, users have a system flag. System users are internal and not
synchronized from external sources.
The system flag is a creation-time attribute of the user account that
can not be changed later on.

Usernames must not be formatted as UUIDs.

# SYNOPSIS

```
soma user-mgmt add ${uname} firstname ${fname} lastname ${lname} employeenr ${num} mailaddr ${addr} team ${team} [system ${bool}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
uname | string | Username of the user | | no
fname | string | Given name of the user | | no
lname | string | Family name of the user | | no
num | integer | Numeric employee ID number of the user | | no
addr | string | Email address of the user | | no
team | string | Name of the team the user belongs to | | no
bool | boolean | Boolean flag if this is a system user | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | user-mgmt | add | yes | no

# EXAMPLES

```
soma user-mgmt add root \
    firstname Charlie \
    lastname Root \
    employeenr 0 \
    mailaddr devnull@example.com \
    team wheel \
    system true
soma user-mgmt add jd \
    firstname Jon \
    lastname Doe \
    employeenr 1234 \
    mailaddr jon@example.com \
    team ExampleTeam
```
