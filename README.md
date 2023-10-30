# No School Till Backend
A service that interacts with Teachable's public API.

# External APIs that it interacts with
- `v1/courses`
- `v1/courses/{course_id}/enrollments`
- `v1/users`

# Getting started
1. Start the application locally
  ```make serve```

2. Make a request to ```localhost:8080/``` in your API testing platform. A list of published courses with enrolled students will be returned.

3. To run tests
  ```make test```

