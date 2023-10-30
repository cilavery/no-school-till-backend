# No School Till Backend
A service that interacts with Teachable's public API.

# External APIs that it interacts with
- `v1/courses`
- `v1/courses/{course_id}/enrollments`
- `v1/users`

# Getting started
1. Make a copy of the ```.sample.env``` file and rename it to ```.env```.
2. Obtain a Teachable API key and set the ```API_KEY``` environment variable with your key.
3. Ensure you are using go `1.21.3` and get all dependencies with `go get`
4. Start the application locally
  ```make serve```
5. Make a request to ```localhost:8080/``` in your API testing platform. A list of published courses with enrolled students will be returned.
6. To run tests
  ```make test```

