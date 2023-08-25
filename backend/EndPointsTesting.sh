# signup
curl -X POST -H "Content-Type: application/json" -d '{"name":"Ashmit Bhardwaj","email":"ashmitbhardwaj05@gmail.com","phonenumber":"9675051814","masterpassword":"carpediem","confirmpassword":"carpediem"}' http://localhost:8000/api/auth/signup

# signin
curl -X POST -H "Content-Type: application/json" -d '{"email":"ashmitbhardwaj05@gmail.com","masterpassword":"carpediem"}' http://localhost:8000/api/auth/signin

# add password
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTA3NDI0NjMsInVzZXJJZCI6IjY0YmQ3NDNlMzk1YTE3ZGZmYjc1NDRlZCJ9.FrmYm8kDXhKfphgNAOSH4j7SkhKU8M1VAIxfJ0Q-wxU" -d '{"name":"Google","username":"ashmitbhardwaj05@gmail.com","url":"www.google.com","password":"something"}' http://localhost:8000/api/passwords/add/64bd743e395a17dffb7544ed

# get all passwords
curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTA3NDI0NjMsInVzZXJJZCI6IjY0YmQ3NDNlMzk1YTE3ZGZmYjc1NDRlZCJ9.FrmYm8kDXhKfphgNAOSH4j7SkhKU8M1VAIxfJ0Q-wxU" http://localhost:8000/api/passwords/all/64bd743e395a17dffb7544ed

# generate password
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTA3NDI0NjMsInVzZXJJZCI6IjY0YmQ3NDNlMzk1YTE3ZGZmYjc1NDRlZCJ9.FrmYm8kDXhKfphgNAOSH4j7SkhKU8M1VAIxfJ0Q-wxU" -d '{"length":8,"numdigits":3,"numsymbols":3,"noupper":false,"allowrepeat":false}' http://localhost:8000/api/passwords/generate

# delete password
curl -X DELETE -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTA3NDI0NjMsInVzZXJJZCI6IjY0YmQ3NDNlMzk1YTE3ZGZmYjc1NDRlZCJ9.FrmYm8kDXhKfphgNAOSH4j7SkhKU8M1VAIxfJ0Q-wxU" http://localhost:8000/api/passwords/delete/64bd743e395a17dffb7544ed\&64bd7783349ee952011c7ef6