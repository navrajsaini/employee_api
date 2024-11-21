Setup:

- set PORT env variable
- start the api (I used go run . in the api working directory)

All api expoints are under:/api/v1

Creating a new entry:
- make POST request to /api/v1/employee
- request body must be a in json format
- the employee_id must be unique otherwise a 400 response will be returned
ex: 
'''
{
    "employee_id": 8,
    "department": "Engineering",
    "job_title": "Senior Enginer",
    "gender": "male"
}
'''

Reading employee entries:
Read all employee entries:
- make GET reques to /api/v1/employees

Read unique employee entry:
- make GET request to /api/v1/employee/<id>
- 200 response with the employee data will be returned

Updating employee entry:
- make a PUT request to /api/v1/employee/<id>
- Request body must be in json format with the new data needed
ex: 
'''
{
    "employee_id": 8,
    "department": "Engineering",
    "job_title": "Senior Enginer",
    "gender": "male"
}
'''

Deleting a employee entry:
- make a DELETE request to /api/v1/employee/<id>
- with the specific id to be deleted
- 200 response with a short description will be returned on success

