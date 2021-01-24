# Neoway Challenge

This is an API made for Neoway's Systems Developer Analyst challenge.

It has 1 endoint that receives a `multipart/form-data` csv or txt data file, containing lines with data separated by commas or spaces.

The goal is to split the data into a mapeable structure to save it in a PostgreSQL database.

Made using Go Modules, PostgreSQL and Docker.

## Challenge Requirements

1. [x] **Create a service in Go that recieves a csv/txt input file**:\
        The data file for the challenge is a txt file, but the API handles both CSV and TXT extensions.
2. [x] **The service must persist all data contained in the file into a relational database (postgresql)**:\
        Data is persisted in a containerized postgresql database.
3. [x] **The data inside the file must be splitted into database columns (can be done in either Go or SQL)**:\
        All string manipulations were made directly in Go using the standard libraries.
4. [x] **Sanitize data on persistance**:\
        All data are treated and sanitized before persistance. Separators are removed from the CPF and CNPJ columns, currency values are converted to `int` types, nullable fields become `NULL` in the db, and dates are converted to go's standard `time.Time` type and persisted as `DATE` on the database.
5. [x] **Validate all CPFs inside the file**:\
        All CPFs, as well as all the CNPJs, are validated using go validator2 library and the popular [brdoc](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwi_-tGa5vvsAhWtJLkGHc1KCeoQFjAAegQIAxAC&url=https%3A%2F%2Fgithub.com%2FNhanderu%2Fbrdoc&usg=AOvVaw0a6XEUu8T5tJX7ZiJoz8cm) library.
6. [x] **All code must be available in a public git repository**:\
        Code hosted in github

## Run

To run the API, use docker-compose.\
First, build it once:

```bash
$ docker-compose build
```

Then run it with:

```bash
$ docker-compose up
```

This command will:

1. Raise a PostgreSQL instance on port 5432.
2. Raise a PgAdmin server on port 16543.
3. Build the API and run it on port 3000.

## Usage

The API has 1 endpoint, `POST /send-data`, that takes a multipart/form-data file, reads, parses and validates the data inside the file, and then persists the information in the database.

You can acess it with the URL `http://localhost:3000/send-data`

To make a request to send a file to the api using CURL, enter the api root folder and type:

```bash
$ curl -v -F 'file=@./data/base_teste.txt' http://localhost:3000/send-data
```

The data file for the challenge is provided inside the `/data` folder to make this process easier

If the request is succesfull, the API will respond with `201: CREATED` and the following response payload:

```json
{
  "message": "Data saved succesfully into the database."
}
```

## Database Setup

When docker-compose starts, a container will be created with an instance of a postgresql database.

**To access the postgresql database:**\
**URL**: `http://localhost:5432`\
**User**: `postgres`\
**Pass**: `admin123`

The container will create a database called `neoway_challenge`, which is defined in the `environments` field inside `docker-compose.yml`.\
When the API is run, the DAO auto creates the `shopping_data` table by inferring the table name from the `ShoppingData` struct, and typing its fields accordindly.

A `pgAdmin` server container is also available to access the database using the browser.

**To access the pgAdmin server:**\
**URL**: `http://localhost:16543`\
**User**: `admin@pg.com`\
**Pass**: `pdadmin`

#### Database Modelling and Field/Column Validations

Analising the data file for the challenge, one could find several thousands of lines containing the following information on each line:

```
* CPF
* Private
* Incompleto
* DataUltimaCompra
* TicketMedio
* TicketUltimaCompra
* LojaMaisFrequente
```

On some of the lines, the last 4 fields are NULL\
On others, 2 fields hold a specific type:

- `DataUltimaCompra` holds a date
- `TicketMedio` holds a currency value

Considering this schema, the following struct was created to hold each line of data contained in the challenge data file:

```go
ShoppingData struct {
  ID                 uint   `gorm:"primaryKey"`
  CPF                string `gorm:"index" validate:"regexp=^\d{3}\.\d{3}\.\d{3}\-\d{2}$"`
  Private            int
  Incompleto         int
  DataUltimaCompra   *time.Time      `gorm:"type:date"`
  TicketMedio        int
  TicketUltimaCompra int
  LojaMaisFrequente  *string         `validate:"regexp=^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$"`
  LojaUltimaCompra   *string         `validate:"regexp=^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$"`
}
```

#### Database Columns, Types and Validation

- An unique and self-incremental `ID` was added to each saved row as best practices mandates.
- The `TicketMedio` and `TicketUltimaCompra` fields, which hold currency values, are stored as `int` values.\
  The rationale behind this is that floating points aren't recommended for storing currency, and go lacks a native representation of the `DECIMAL` SQL type.\
  Storing as int is also a common practice in several languages, given that the value can be parsed easily back to currency format by making the last 2 digits the cents.
- The field with pointers **\*** in their types are the `nullable` fields that could come inside the data file. In order to make a struct field null, its type has to be a pointer.\
  The other fields which are garanteed to come are mapped as literal types.
- The `gorm` annotations in some of the struct fields are used to:
  - Define the ID as the primary key
  - Make `CPF` an index.
  - Specify the dates as SQL `DATE` type
- The struct fields are all validated and sanitized before going to the database.
- `CPF`, `LojaMaisFrequente`and`LojaUltimaCompra` fields are checked against a regex using **`go validator2`** library.\
  A last validation on these fields is also made in each of these fields before adding them to the map that is sent to be saved in the struct using the [brdoc](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwi_-tGa5vvsAhWtJLkGHc1KCeoQFjAAegQIAxAC&url=https%3A%2F%2Fgithub.com%2FNhanderu%2Fbrdoc&usg=AOvVaw0a6XEUu8T5tJX7ZiJoz8cm) library.

### Design and Architectural explanation

This is a simple API, which could have been done with a lot less files and folders.\
The reason i made it this way is to demostrate an API strucutured by following the guidelines of the [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/), in which the main 3 layers of the application, the **`Business`**, **`Infrastructure`** and **`Interface`** layers, are separated from each other in order to allow for better scalability and maintainability by isolating these layers as much as its possible.

This separation allows any modifications in said layers to be made independently, allowing, for instance, a complete change in the database driver without having to mess in the business logic or the user interface.

Altought this api is far from being 100% conforming to the principles and guidelines of this architecture, its based on it and, by doing so, does benefit from some of its advantages:

- The **`/api`** folder represents the user **`Interface`** layer, which define de routes that are used to interact with the application
- The **`/api/handlers`** folder represents the **`Business`** layer, and contains the business logic
- The **`/dao`** folder represents the **`Infraestructure`** layer, which holds the Database operations, monitoring and related concerns.

This means there is only one place to look for when adding new interfaces (routes) to the api, adding or changing business logic (handlers), or make infrastructural modifications (DAO).

### Database Access Object

The database is manipulated using a DAO _(Database Access Object)_, a standard used for acessing data in databases witch allows further separation between the business logic and the rules to acess the DBMS.

The DAO object is usually used in _Object Oriented_ languages, but it suits Go perfectly, since its design can be used to map the database entities into respective Go Structs fields. These structs are then translated into database entities, and their fields into columns, when database operations are performed.

For these and more reasons, its fairly common to see the DAO model being used among Go projects.

This API uses [GORM](https://gorm.io/index.html) in its DAO.\
GORM provides an API which makes database operations easier and safer, offering field validations, batch insert using transactions, among with several other benefits.

### Configuration

Saving configuration in environment variables is recommended as Rule III of the [Twelve-Factor App](https://12factor.net/pt_br/config).

All relevant configuration, such as external the API PORT and the Database information, are defined inside go files inside the `/config` folder.
This is Go's native way of dealing with `environment variables`, by assigning values to Constant types that can be read from anywhere in the application as a module.

These config files are annotated with go build tags, meaning a flag can be specified in order to determine which of the config files go will use to parse values: `config_dev.go`, which contains literal values for development, and `config_prod.go`, which reads values from the environment and can be used for deployment in production.

To run the API in `development` mode:

```bash
$ go run -tags dev cmd/main.go
```

To run using the production config file:

```bash
$ go run -tags prod cmd/main.go
```

The same rule goes for building:

```bash
$ go build -tags dev -o main cmd/main.go
$ go build -tags prod -o main cmd/main.go
```

To run the compiled build:

```bash
$ ./main
```

**_PS_\***: _In order to run the api this way, you need to start only the database first with docker-compose_

```bash
$ docker-compose up -d postgres
```
