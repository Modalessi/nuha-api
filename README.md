# Nuha api

competitive programming



### To do

- register, sign in. [DONE]
- restore password method
- send submissions [MOSTLY DONE]
- CRUD operatiosn for problems 
  - admin create [DONE]
  - get problems list [DONE]
  - get problem [DONE]
  - admin DELETE problem [DONE]
  - admin EDIT problem [DONE]




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



### Submission
- id
- problemId (problems.id foreign key)
- user_id (user.id foreign key)
- language_id int
- source_code string
- status text (PENDING, ACCEPETED, WRONG ANSWER...)
- max_time
- max_memory
- created_at


### Submission Results
- id 
- submission_id (submission.id foriegn key)
- test_case_input string (used to determin which test case failed later)
- judge_token string
- time_used float64
- memory_used float64
- judge_respones JSONB
- created_at timestamp


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