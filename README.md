# PostgreSQL Setup with Docker

### 1. **Install Docker**
- Install Docker Desktop by visiting [Docker Download Page](https://www.docker.com/).
- If you face any installation issues, follow this [Docker Troubleshooting Guide](https://docs.docker.com/desktop/cert-revoke-solution/).
- Useful Tutorial https://www.youtube.com/watch?v=Hs9Fh1fr5s8&t=155s
### 2. **Run PostgreSQL in Docker**
Run the following command to start a PostgreSQL container with Docker:

```bash
docker run --name some_postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql/data \
  -d postgres
```

**Explanation:**
- `--name some_postgres`: Assigns a name to the container.
- `-e POSTGRES_USER=myuser`: Sets the PostgreSQL user.
- `-e POSTGRES_PASSWORD=mypassword`: Sets the PostgreSQL password.
- `-e POSTGRES_DB=mydb`: Specifies the database to be created.
- `-p 5432:5432`: Maps the PostgreSQL container port (5432) to your host machine's port (5432).
- `-v pgdata:/var/lib/postgresql/data`: Mounts a Docker volume `pgdata` to store PostgreSQL data persistently.
- `-d`: Runs the container in detached mode (in the background).

This command will **pull the latest PostgreSQL image** if it is not available locally. If it's already pulled, it will use the local image.

### 3. **Managing Data with Volumes**
- The volume `pgdata` is mounted to `/var/lib/postgresql/data` inside the container. This ensures that any data PostgreSQL writes will persist even if the container is removed or restarted.
- If the `pgdata` volume doesn't exist, Docker will **automatically create it** to store data.

### 4. **Create a Custom Volume (Optional)**
You can explicitly create a custom volume and mount it to PostgreSQL data directory:

```bash
docker volume create new-volume-name
docker run --name some_postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v new-volume-name:/var/lib/postgresql/data \
  -d postgres
```

### 5. **Port Conflict**
If there is any process already running on port 5432, you can check it with the following command:

```bash
sudo lsof -i :5432
```

Once you identify the process, kill it using its Process ID (PID):

```bash
sudo kill <PID>
```

For example:

```bash
sudo kill 315
```

After killing the conflicting process, restart your PostgreSQL container:

```bash
docker start some_postgres
```

---



**GORM Relationships in Golang with PostgreSQL**

GORM is a powerful ORM for Golang, allowing seamless integration with databases. To use GORM with PostgreSQL, install the necessary dependencies:

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

Then, initialize the database connection:

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
    dsn := "host=localhost user=myuser password=mypassword dbname=gofiberdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    return db, err
}
```

### Belongs To Relationship
A **belongs to** relationship establishes a many-to-one (or one-to-one) connection where each instance of a model belongs to another model. This means the foreign key is stored in the table of the dependent model.

```go
type User struct {
    gorm.Model
    Name      string
    CompanyID uint
    Company   Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:SET NULL;OnUpdate:CASCADE;"`
}

type Company struct {
    ID   uint `gorm:"primaryKey"`
    Name string
}
```

In this setup, `User` belongs to `Company`, and `CompanyID` is the foreign key. The `Company` struct does not need a reference to `User`. If a company is deleted, the `CompanyID` in `User` will be set to NULL, preventing orphaned references.

To migrate and create records:

```go
err = db.AutoMigrate(&Company{}, &User{})

company := Company{Name: "TechCorp"}
db.Create(&company)

user := User{Name: "John Dallas", CompanyID: company.ID}
db.Create(&user)
```

Alternatively, GORM automatically assigns `CompanyID` if we pass the `Company` object:

```go
db.Create(&User{Name: "Huraira", Company: company})
```

Fetching users and their associated company:

```go
var user User
db.First(&user, 1)
fmt.Println(user.Company) // Outputs an empty struct {}
```

To preload related data:

```go
var user1 User
db.Preload("Company").First(&user1, 1)
fmt.Println(user1.Company) // Now prints the related company details
```

Filtering users by company:

```go
var users []User
db.Where("company_id = ?", 1).Find(&users)
fmt.Println(users)
```

Dropping tables:

```go
err = db.Migrator().DropTable(Company{}, User{})
```

### Has One Relationship
A **has one** relationship means that one entity owns another. The owner model contains a reference to the owned model, but the foreign key is stored in the owned model.

```go
type Company struct {
    gorm.Model
    Name string
    CEO  CEO `gorm:"foreignKey:CompanyID"` // A company has one CEO
}

type CEO struct {
    gorm.Model
    Name      string
    CompanyID uint
}
```

Here, `Company` has a `CEO`, but the foreign key (`CompanyID`) is stored in the `CEO` table. When querying a company, its CEO is automatically fetched if preloaded:

```go
var company Company
db.Preload("CEO").First(&company, 1)
fmt.Println(company.CEO)
```

### Has Many Relationship
A **has many** relationship means that one model owns multiple instances of another model. The foreign key is stored in the owned model.

```go
type Company struct {
    gorm.Model
    Name  string
    Users []User `gorm:"foreignKey:CompanyID"` // A company has many users
}

type User struct {
    gorm.Model
    Name      string
    CompanyID uint
}
```

Fetching a company with its users:

```go
var company Company
db.Preload("Users").First(&company, 1)
fmt.Println(company.Users)
```

### Many-to-Many Relationship
A **many-to-many** relationship means that multiple instances of one model are associated with multiple instances of another model. This requires a join table.

```go
type User struct {
    gorm.Model
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name string
}
```

Here, `user_languages` is the join table that maps users to languages. This avoids placing a `UserID` in `Language` because a language can be associated with multiple users. To make both directions explicit:

```go
type User struct {
    gorm.Model
    Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name  string
    Users []*User `gorm:"many2many:user_languages;"`
}
```

Fetching users and their languages:

```go
var users []User
db.Model(&User{}).Preload("Languages").Find(&users)
fmt.Println(users)
```

Similarly, fetching a specific user’s languages:

```go
var user User
db.Preload("Languages").First(&user, 1)
fmt.Println(user.Languages)
```

### Handling Preloading and Eager Loading
Using `Preload("RelationName")` ensures that related entities are fetched along with the primary entity. If not used, related data will be an empty struct unless explicitly queried later.

```go
var user User
db.First(&user, 1)
fmt.Println(user.Company) // Empty struct {}

var user1 User
db.Preload("Company").First(&user1, 1)
fmt.Println(user1.Company) // Populated with company details
```

To delete a related entity, constraints like `OnDelete:CASCADE` or `OnDelete:SET NULL` dictate the behavior. For example:

```go
type Company struct {
    gorm.Model
    Name string
    CEO  CEO `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;"`
}
```

If a `Company` is deleted, its `CEO` is also deleted. In contrast:

```go
type User struct {
    gorm.Model
    Name      string
    CompanyID uint
    Company   Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:SET NULL;OnUpdate:CASCADE;"`
}
```

Here, deleting a `Company` sets `CompanyID` in `User` to NULL rather than deleting the user.

These relationships allow flexible data modeling in GORM, ensuring efficient and meaningful database interactions.

