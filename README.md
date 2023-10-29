# No School Till Backend
A service that interacts with Teachables public API.

# APIs that it interacts with
- `v1/courses`
- `v1/courses/{course_id}/enrollments`
- `v1/users`

# Getting started
1. Run the application locally
  ```make serve```

2. Run tests
  ```make test```

3. To get student and course data:
To get all users of a school:
```localhost:8080/```

To get all courses:
```localhost:8080/courses```

To get all enrollments by course id:
```localhost:8080/courses/2002430/enrollments```

# Known issues
1. What if there are many students? How to optimize performance?
2. Caching in Redis and when to call all users?
