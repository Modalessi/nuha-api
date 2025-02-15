# Nuha api

competitive programming



### To do


- CORS configuration
- Auth (4 / 5)
  - register [DONE]
  - login [DONE]
  - logout [DONE]
  - restore password method
  - verify email method [DONE]


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
- public bool
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

****
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



#### Contest
- id uuid
- title string
- problems []Problems.id
- start_time data
- particpants []Users.id
- duration time.Duration


## Contest Flow

### creating
1. create new contest POST /contest
2. add problems to the contest PUT /contest
3. only can add not public problems


### particpating

#### before start
1. get /contest (returns all contests titles, number of participants ,ids)
2. post /participate?id=<contest_id>


1. contest start
2. GET /contest?id=<contest_id 
   1. check the user is participated in the contest 
   2. if participated returns all the contest deatils (title, problems, elapsed time)
3. user chooses a problem
4. GET /problem?id=<problem_id>
   1. check the user is particpated 
   2. if participated he get the problem back
   3. if not he get no access error
5. submit solution POST /submit 
   1. same check the user is particpated in the combeition
   2. if not return an error
   3. if yes register the submission with the contest id
6. 




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
