# Nuha api

competitive programming



### To do


- CORS configuration
- Auth (1/ 3)
  - register, sign in. [DONE]
  - restore password method
  - verify email method


- Better Storage Interface


- Submissions (6/ 6)
  - send submissions [DONE]
  - get single submission [DONE]
  - get user submssions [DONE]
  - get all submsssions [DONE]
  - Pagination [DONE]

- Problems (5 / 5)
  - admin create [DONE]
  - get problems list [DONE]
  - get problem [DONE]
  - admin DELETE problem [DONE]
  - admin EDIT problem [DONE]
  - Add Testcases using files [DONE]



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

- Testcases
  - id
  - problem id (problem.id foregin key)
  - number int
  - stdin (blob)
  - expected_output (blob)
  - unique (problem id, number)

- problems_descriptions
  - problem_id primary (problem.id foregin key)
  - description blob


#### Submission
- id
- problemId (problems.id foreign key)
- user_id (user.id foreign key)
- language_id int
- source_code string
- status text (PENDING, ACCEPETED, WRONG ANSWER...)
- created_at


#### Submission Results
- id
- submission_id (submission.id foriegn key)
- judge_token string
- test_case_input string
- stdout string
- expected_output string
- status string
- time_used float64
- memory_used float64
- judge_respones JSONB
- created_at timestamp





### Submission flow

- endpoint recives submissiosn (langauge, code, problemid)
- create the submission
- store it
- pass the submissionJobs to (submissionsPipeline)
- return the submission id

**SubmissionPipeline**
a piple line that runs in the background
takes a submission job and proccess in three steps


submissionsChan                                                               resultsChan
----------------------->|send submssions to judge api, and take the tokens| ---------------->

resultsChan                                                                  dbUpdateChan
----------------------->|polling judge api, to get results for the tokens|------------------->

dbUpdateChan
------------------------>|store result from judge api to db|


so we have three chanels
- SubmissionChan
- resultsChan
- dbUpdateChan

and three background porccess
- submissionProcessor
- resultProcessor
- dbWriter







- the frontend then can get the result using get_submission route

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
