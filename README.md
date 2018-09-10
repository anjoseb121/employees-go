# employees-go

## Deploy
1. Build the binary `GOOS=linux go build -o main main.go`

2. Zip the binary
`zip deployment.zip main`

3. Create db instance
```
aws rds create-db-instance \
    --db-instance-identifier MySQLForLambdaTest \
    --db-instance-class db.t2.micro \
    --engine MySQL \
    --allocated-storage 5 \
    --no-publicly-accessible \
    --db-name ExampleDB \
    --master-username username \
    --master-user-password password \
    --backup-retention-period 3 
```