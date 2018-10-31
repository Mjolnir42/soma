# user management

User management are the functions for maintaining the user database
within SOMA.

# SYNOPSIS OVERVIEW

```
soma user-mgmt add ${uname} firstname ${fname} lastname ${lname} employeenr ${num} mailaddr ${addr} team ${team} [system ${bool}]
soma user-mgmt update ${userID} username ${uname} firstname ${fname} lastname ${lname} employeenr ${num} mailaddr ${addr} team ${team} [deleted ${bool}]
soma user-mgmt remove ${uname}
soma user-mgmt purge ${uname}
soma user-mgmt show ${uname}
soma user-mgmt list
soma user-mgmt sync
```

See `soma user-mgmt help ${command}` for detailed help.
