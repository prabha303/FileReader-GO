

```

Create Employee
http://localhost:9000/v1/employees

{
    "id" : 1,
    "name": "Prabha",
    "address" : "new street",
    "employeeNumber": "EMP001"

}

```


```

To List All Employee
http://localhost:9000/v1/employees

```


```

Find by Employee ID
http://localhost:9000/v1/employees/search/id/:id

Find by name
http://localhost:9000/v1/employees/search/name/:name



```



```

Update Employee by Employee ID
http://localhost:9000/v1/employee/update/id/0

{
    "address" : "New Address"
}
```