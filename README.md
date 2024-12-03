# Nuha api

competitive programming



### To do

- register, sign in. [DONE]
- send submissions [MOSTLY DONE]
- CRUD operatiosn for problems (admin only) 




### Data Design 

#### Problem

- id
- title
- description (in markdown) (stored in s3)
- test cases (in .in & .out) (stored in s3)
  - input
  - output
- tags []string
- difficulty string
- time limit float
- memory limit float
- created_at
- updated_at



some users

-- normal user
{
    "email": "mohammed@email.com",
    "password": "Moh12345"

}


-- admin
{
    "name": "the boss",
    "email": "admin@email.com",
    "password": "Adm12345"
}