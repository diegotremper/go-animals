# Go Animals

Sample implementation `gin`, `postges` using `sqlx`.

This project will create a database named `sampledb` along with a table named `Animal` :

```
sampledb=# SELECT * FROM animals;
 id | name | age | description
----+------+-----+-------------
(0 rows)
```

## Prerequisites :

- install `docker`

## Running :

- `make environment`
- create file `.env` based on `.env.sample`
- `make server`
- app will running in port 8080!
  Note : you can use command `make help` for showing list available commands.

## Testing Locally using curl :

1. Create animal

```bash
curl --location 'http://localhost:8080/animals' \
--header 'Content-Type: application/json' \
--data '{
    "name": "cow",
    "age": 20,
    "description": "beautiful cow"
}'
```

2. List all animal

```bash
curl --location 'http://localhost:8080/animals'
```

3. Detail animal

```bash
curl --location 'http://localhost:8080/animals/11'
```

4. Update animal

```bash
curl --location --request PATCH 'http://localhost:8080/animals/12' \
--header 'Content-Type: application/json' \
--data '{
    "name": "cat",
    "age": 15,
    "description": "beautiful cat update"
}'
```

5. Delete animal

```bash
curl --location --request DELETE 'http://localhost:8080/animals/13'
```
